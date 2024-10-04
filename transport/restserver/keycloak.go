package restserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	commonContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/keycloak"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
	"github.com/rosaekapratama/go-starter/slices"
	commonStrings "github.com/rosaekapratama/go-starter/strings"
	"strings"
)

var additionalClaims []string

type KeycloakContextOption interface {
	Apply(ctx context.Context, c *gin.Context, claims map[string]interface{})
}

type KeycloakContextAdditionalClaimOption struct {
}

type KeycloakContextWhitelistedHostOption struct {
	whilistedHosts []string
}

func (o *KeycloakContextAdditionalClaimOption) Apply(_ context.Context, _ *gin.Context, _ map[string]interface{}) {
}

func (o *KeycloakContextWhitelistedHostOption) Apply(ctx context.Context, c *gin.Context, claims map[string]interface{}) {
	// Only process if host is within whitelist
	host := c.GetHeader(headers.XForwardedHost)
	if host != str.Empty && slices.ContainStringCaseInsensitive(o.whilistedHosts, host) {
		log.Tracef(ctx, "%s is whitelisted, host=%s", headers.XForwardedHost, host)
		// Try to get token claim value from header and put it to claims map
		for k, v := range c.Request.Header {
			claim := commonStrings.DashToSnake(k)
			claims[claim] = slices.ToString(v, sym.Comma)
			log.Tracef(ctx, "Set claim from header, claim=%s, headerKey=%s, headerValue=%s", claim, k, v)
		}

		// Try to get token claim value from query param and put it to claims map
		queryParams := c.Request.URL.Query()
		for k, v := range queryParams {
			claim := commonStrings.CamelToSnake(k)
			claims[claim] = slices.ToString(v, sym.Comma)
			log.Tracef(ctx, "Set claim from query param, claim=%s, queryKey=%s, queryValue=%s", claim, k, v)
		}
	}
}

func WithAdditionalClaim(claims ...string) KeycloakContextOption {
	additionalClaims = append(additionalClaims, claims...)
	return &KeycloakContextAdditionalClaimOption{}
}

func WithWhitelistedHost(hosts ...string) KeycloakContextOption {
	return &KeycloakContextWhitelistedHostOption{whilistedHosts: hosts}
}

func InjectKeycloakContext(options ...KeycloakContextOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip if request is health check
		if isHealthCheckPath(c) {
			c.Next()
			return
		}

		// start middleware logic
		ctx := c.Request.Context()
		w := c.Writer

		claims := make(map[string]interface{})
		for _, option := range options {
			option.Apply(ctx, c, claims)
		}
		// get keycloak token from context which already set by injectAuthContext function
		if tokenStr, ok := commonContext.TokenFromContext(ctx); ok {

			// Split JWT into header, payload, and signature
			parts := strings.Split(tokenStr, ".")
			payloadBase64 := parts[1]

			// Decode payload from base64
			payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadBase64)
			if err != nil {
				log.Error(ctx, err, "Error decoding token payload")
				SetResponse(w, response.GeneralError)
				c.Abort()
				return
			}

			// Parse JSON payload to access claims
			if err := json.Unmarshal(payloadBytes, &claims); err != nil {
				log.Error(ctx, err, "Error parsing token JSON payload")
				SetResponse(w, response.GeneralError)
				c.Abort()
				return
			}
		} else {
			log.Tracef(ctx, "Unable to set context from keycloak token, path=%s, method=%s", c.Request.URL.Path, c.Request.Method)
		}

		// set common token claim to context
		// set sub claim to context if exists
		if v, ok := claims[keycloak.ClaimSub]; ok {
			sc := v.(string)
			if sc == str.Empty {
				log.Warn(ctx, "Unable to set context, sub claim is empty")
			} else {
				ctx = commonContext.ContextWithUserId(ctx, sc)
			}
		} else {
			log.Warn(ctx, "Unable to set context, missing sub claim")
		}

		// set realm from issuer claim to context if exists
		if v, ok := claims[keycloak.ClaimIss]; ok {
			iss := v.(string)
			if iss == str.Empty {
				log.Warn(ctx, "Unable to set context, iss claim is empty")
			} else {
				startIdx := strings.Index(iss, realmsPath) + len(realmsPath)
				ctx = commonContext.ContextWithRealm(ctx, iss[startIdx:])
			}
		} else {
			log.Warn(ctx, "Unable to set context, missing iss claim")
		}

		// set preferred_username claim to context if exists
		var puc string
		if v, ok := claims[keycloak.ClaimPreferredUsername]; ok {
			puc = v.(string)
			if puc == str.Empty {
				log.Warn(ctx, "Unable to set context, preferred_username claim is empty")
			} else {
				ctx = commonContext.ContextWithUsername(ctx, puc)
			}
		} else {
			log.Warn(ctx, "Unable to set context, missing preferred_username claim")
		}

		// set name claim to context if exists, or set with username if empty
		var nc string
		if v, ok := claims[keycloak.ClaimName]; ok {
			nc = v.(string)
			if nc == str.Empty {
				log.Trace(ctx, "Name claim is empty, set context using preferred_username claim value instead")
				nc = puc
			}
		} else {
			log.Trace(ctx, "Missing name claim, set context using preferred_username claim value instead")
			nc = puc
		}
		ctx = commonContext.ContextWithFullName(ctx, nc)

		// set email claim to context if exists
		if v, ok := claims[keycloak.ClaimEmail]; ok {
			ec := v.(string)
			if ec == str.Empty {
				log.Trace(ctx, "Unable to set context, email claim is empty")
			} else {
				ctx = commonContext.ContextWithEmail(ctx, ec)
			}
		} else {
			log.Trace(ctx, "Unable to set context, missing email claim")
		}

		// set realm access claim to context if exists
		if v, ok := claims[keycloak.ClaimRealmAccess]; ok {
			rac := v.(map[string]interface{})
			if len(rac) == 0 {
				log.Trace(ctx, "Unable to set context, realm access claim is empty")
			} else if v, exists := rac[keycloak.ClaimRealmAccessRoles]; exists {
				rc := v.([]interface{})
				if len(rc) == 0 {
					log.Trace(ctx, "Unable to set context, roles of realm access claim is empty")
				} else {
					roles := make([]string, 0)
					for _, role := range rc {
						roles = append(roles, role.(string))
					}
					ctx = commonContext.ContextWithRoles(ctx, roles)
				}
			}
		} else {
			log.Trace(ctx, "Unable to set context, missing email claim")
		}

		// set custom provided claim to context if exists
		for _, additionalClaim := range additionalClaims {
			k := commonStrings.SnakeToCamel(additionalClaim)
			if v, ok := claims[additionalClaim]; ok {
				ctx = context.WithValue(ctx, k, v)
				commonContext.AddManagedKey(k)
			} else {
				log.Tracef(ctx, "Unable to set context, missing %s claim", additionalClaim)
			}
		}

		// override request context with new context
		c.Request = c.Request.WithContext(ctx)
	}
}
