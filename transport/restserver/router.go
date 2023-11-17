package restserver

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/location"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	myContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/healthcheck"
	"github.com/rosaekapratama/go-starter/keycloak"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"
	"github.com/rosaekapratama/go-starter/transport/constant"
	"github.com/rosaekapratama/go-starter/transport/logging/models"
	"github.com/rosaekapratama/go-starter/transport/logging/repositories"
	"github.com/rosaekapratama/go-starter/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	contentTypeApplicationJson = "application/json"
)

var (
	cfg           *config.Object
	propagator    = otel.GetTextMapPropagator()
	Router        *gin.Engine
	logRepository repositories.IRestLogRepository
)

func cloneRequest(req *http.Request) (*http.Request, error) {
	// Clone the request
	clone := new(http.Request)
	*clone = *req

	// Copy the Headers
	clone.Header = make(http.Header, len(req.Header))
	for k, v := range req.Header {
		clone.Header[k] = append([]string(nil), v...)
	}

	// Copy the Body
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	clone.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Restore the Body in the original request
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return clone, nil
}

func logging(c *gin.Context) {
	// Skip if request is health check
	isHealthCheck, _ := regexp.MatchString(healthcheck.URLPathRegex, c.Request.URL.Path)
	if isHealthCheck {
		return
	}

	var ctx context.Context
	var body *string
	r, err := cloneRequest(c.Request)
	if err != nil {
		log.Error(ctx, err, "Failed to clone http request for logging")
	} else {
		ctx = r.Context()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(ctx, err, "Failed to read a copy of body request for logging")
		}
		body = utils.StringP(string(bodyBytes))
		if cfg.Transport.Server.Rest.Logging.Stdout {
			httpFields := make(map[string]interface{})
			httpFields[constant.LogTypeFieldKey] = constant.LogTypeHttp
			httpFields[constant.UrlFieldKey] = r.URL
			httpFields[constant.MethodFieldKey] = r.Method
			httpFields[constant.IsServerFieldKey] = true
			httpFields[constant.IsRequestFieldKey] = true
			httpFields[constant.HeadersFieldKey] = r.Header
			if r.MultipartForm != nil {
				httpFields[constant.FormDataFieldKey] = r.MultipartForm
			}
			if body != nil {
				httpFields[constant.BodyFieldKey] = body
			}
			log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Info()
		}

		if cfg.Transport.Server.Rest.Logging.Database != str.Empty {
			err := logRepository.Save(ctx, &models.TransportRestLog{
				ID:           uuid.New(),
				TraceID:      myContext.TraceIdFromContext(ctx),
				SpanID:       myContext.SpanIdFromContext(ctx),
				IsServer:     true,
				IsRequest:    true,
				URL:          fmt.Sprintf("%v", r.URL),
				Method:       r.Method,
				Headers:      utils.StringP(fmt.Sprintf("%v", r.Header)),
				Body:         body,
				StatusCode:   nil,
				ErrorMessage: nil,
				ProcessDT:    time.Now().In(location.AsiaJakarta),
				ProcessBy:    config.Instance.GetObject().App.Name,
			})
			if err != nil {
				log.Error(ctx, err, "Failed to save REST server request log")
			}
		}
	}

	c.Next()
	i := c.Writer.(*WriterInterceptor)

	if cfg.Transport.Server.Rest.Logging.Stdout {
		httpFields := make(map[string]interface{})
		httpFields[constant.LogTypeFieldKey] = constant.LogTypeHttp
		httpFields[constant.UrlFieldKey] = r.URL
		httpFields[constant.MethodFieldKey] = r.Method
		httpFields[constant.IsServerFieldKey] = true
		httpFields[constant.IsRequestFieldKey] = false
		httpFields[constant.StatusCodeFieldKey] = i.status
		httpFields[constant.HeadersFieldKey] = i.ResponseWriter.Header()
		if len(i.body) > integer.Zero {
			httpFields[constant.BodyFieldKey] = string(i.body)
		}
		log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Info()
	}

	if cfg.Transport.Server.Rest.Logging.Database != str.Empty {
		if len(i.body) > integer.Zero {
			body = utils.StringP(string(i.body))
		}
		err := logRepository.Save(ctx, &models.TransportRestLog{
			ID:           uuid.New(),
			TraceID:      myContext.TraceIdFromContext(ctx),
			SpanID:       myContext.SpanIdFromContext(ctx),
			IsServer:     true,
			IsRequest:    false,
			URL:          fmt.Sprintf("%v", r.URL),
			Method:       r.Method,
			Headers:      utils.StringP(fmt.Sprintf("%v", i.ResponseWriter.Header())),
			Body:         body,
			StatusCode:   utils.StringP(strconv.Itoa(i.status)),
			ErrorMessage: utils.StringP(i.response.Description()),
			ProcessDT:    time.Now().In(location.AsiaJakarta),
			ProcessBy:    config.Instance.GetObject().App.Name,
		})
		if err != nil {
			log.Error(ctx, err, "Failed to save REST server response log")
		}
	}
}

func enableCors() gin.HandlerFunc {
	ctx := context.Background()
	// Skip if cors is disabled
	propsConfig := cfg.Cors
	if propsConfig != nil && propsConfig.Enabled {
		corsConfig := cors.DefaultConfig()
		if propsConfig.AllowOrigins != nil && len(propsConfig.AllowOrigins) > integer.Zero {
			corsConfig.AllowOrigins = propsConfig.AllowOrigins
		}
		if propsConfig.AllowMethods != nil && len(propsConfig.AllowMethods) > integer.Zero {
			corsConfig.AllowMethods = propsConfig.AllowMethods
		}
		if propsConfig.AllowHeaders != nil && len(propsConfig.AllowHeaders) > integer.Zero {
			corsConfig.AllowHeaders = propsConfig.AllowHeaders
		}
		if propsConfig.ExposeHeaders != nil && len(propsConfig.ExposeHeaders) > integer.Zero {
			corsConfig.ExposeHeaders = propsConfig.ExposeHeaders
		}
		corsConfig.AllowCredentials = propsConfig.AllowCredentials
		if propsConfig.MaxAge > integer.Zero {
			corsConfig.MaxAge = time.Duration(propsConfig.MaxAge) * time.Second
		}
		log.Info(ctx, "CORS is enabled")
		log.Debugf(ctx, "CORS allow credentials : %s", corsConfig.AllowCredentials)
		log.Debugf(ctx, "CORS allow origins     : %v", corsConfig.AllowOrigins)
		log.Debugf(ctx, "CORS allow methods     : %v", corsConfig.AllowMethods)
		log.Debugf(ctx, "CORS allow headers     : %v", corsConfig.AllowHeaders)
		log.Debugf(ctx, "CORS expose headers    : %v", corsConfig.ExposeHeaders)
		log.Debugf(ctx, "CORS max age           : %i second", propsConfig.MaxAge)
		return cors.New(corsConfig)
	} else {
		log.Info(ctx, "CORS is disabled")
		return func(c *gin.Context) {
			// It is intended to be empty function so it does nothing
		}
	}
}

func isHealthCheckPath(c *gin.Context) bool {
	r := c.Request
	w := c.Writer
	isHealthCheck, err := regexp.MatchString(healthcheck.URLPathRegex, r.URL.Path)
	if err != nil {
		log.Error(r.Context(), err)
		SetResponse(w, response.GeneralError)
		c.Abort()
		return false
	}
	if isHealthCheck && r.Method == http.MethodGet {
		SetResponse(w, response.Success)
		return true
	}
	return false
}

func extractTraceParent(c *gin.Context) {
	traceparent := strings.TrimSpace(c.GetHeader(headers.Traceparent))
	if traceparent != str.Empty {
		// Override request context with new context if traceparent is exists
		c.Request = c.Request.WithContext(myContext.ContextWithTraceParent(c.Request.Context(), traceparent))
	}
}

func interceptResponse(c *gin.Context) {
	r := c.Request
	i := NewWriterInterceptor(r.Context(), c.Writer)
	c.Writer = i
	c.Next()
	contentType := strings.ToLower(strings.Split(i.ResponseWriter.Header().Get(headers.ContentType), sym.SemiColon)[0])
	if contentType == str.Empty {
		i.ResponseWriter.Header().Set(headers.ContentType, contentTypeApplicationJson)
		contentType = contentTypeApplicationJson
	}
	if !i.IsWritten && contentType == contentTypeApplicationJson {
		_, err := i.Write(nil)
		if err != nil {
			log.Error(r.Context(), err)
		}
	}
}

func injectTraceParent(c *gin.Context) {
	// Skip if request is health check
	if isHealthCheckPath(c) {
		c.Next()
		return
	}

	// Start middleware logic
	headerCarrier := propagation.HeaderCarrier{}
	propagator.Inject(c.Request.Context(), headerCarrier)
	for _, key := range headerCarrier.Keys() {
		c.Writer.Header().Add(key, headerCarrier.Get(key))
	}
}

func injectAuthContext(c *gin.Context) {
	// Skip if request is health check
	if isHealthCheckPath(c) {
		c.Next()
		return
	}

	// Start middleware logic
	ctx := c.Request.Context()
	auth := c.GetHeader(headers.Authorization)
	if auth != str.Empty && strings.HasPrefix(auth, headers.BasicAuthPrefix) {
		auth = strings.TrimPrefix(auth, headers.BasicAuthPrefix)

		// Decode basic auth
		bs, err := base64.StdEncoding.DecodeString(auth)
		if err != nil {
			log.Error(ctx, err, "Failed to decode basic auth")
		} else {
			// Get the username from decoded value and set to context
			username := strings.Split(string(bs), sym.Colon)[integer.Zero]
			ctx = myContext.ContextWithUsername(ctx, username)

			// Override request context with new context
			c.Request = c.Request.WithContext(ctx)
		}
	} else if auth != str.Empty && strings.HasPrefix(auth, headers.BearerTokenPrefix) {
		token := strings.TrimPrefix(auth, headers.BearerTokenPrefix)
		ctx = myContext.ContextWithToken(ctx, token)

		// Override request context with new context
		c.Request = c.Request.WithContext(ctx)
	} else {
		log.Tracef(ctx, "Authorization is empty, skip inject auth context process, path=%s, method=%s", c.Request.URL.Path, c.Request.Method)
	}
}

func InjectKeycloakContext(additionalClaims ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip if request is health check
		if isHealthCheckPath(c) {
			c.Next()
			return
		}

		// start middleware logic
		ctx := c.Request.Context()
		w := c.Writer

		// get keycloak token from context which already set by injectAuthContext function
		if tokenStr, ok := myContext.TokenFromContext(ctx); ok {

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
			var claims map[string]interface{}
			if err := json.Unmarshal(payloadBytes, &claims); err != nil {
				log.Error(ctx, err, "Error parsing token JSON payload")
				SetResponse(w, response.GeneralError)
				c.Abort()
				return
			}

			// set common token claim to context
			// set sub claim to context if exists
			if v, ok := claims[keycloak.ClaimSub]; ok {
				sc := v.(string)
				if sc == str.Empty {
					log.Warn(ctx, "Unable to set context, sub claim is empty")
				} else {
					ctx = myContext.ContextWithUserId(ctx, sc)
				}
			} else {
				log.Warn(ctx, "Unable to set context, missing sub claim")
			}

			// set preferred_username claim to context if exists
			var puc string
			if v, ok := claims[keycloak.ClaimPreferredUsername]; ok {
				puc = v.(string)
				if puc == str.Empty {
					log.Warn(ctx, "Unable to set context, preferred_username claim is empty")
				} else {
					ctx = myContext.ContextWithUsername(ctx, puc)
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
			ctx = myContext.ContextWithFullName(ctx, nc)

			// set email claim to context if exists
			if v, ok := claims[keycloak.ClaimEmail]; ok {
				ec := v.(string)
				if ec == str.Empty {
					log.Trace(ctx, "Unable to set context, email claim is empty")
				} else {
					ctx = myContext.ContextWithEmail(ctx, ec)
				}
			} else {
				log.Trace(ctx, "Unable to set context, missing email claim")
			}

			// set developer provided claim to context if exists
			for _, additionalClaim := range additionalClaims {
				if v, ok := claims[additionalClaim]; ok {
					ctx = context.WithValue(ctx, additionalClaim, v)
				} else {
					log.Tracef(ctx, "Unable to set context, missing %s claim", additionalClaim)
				}
			}

			// override request context with new context
			c.Request = c.Request.WithContext(ctx)
		} else {
			log.Trace(ctx, "Unable to set context, missing keycloak token, path=%s, method=%s", c.Request.URL.Path, c.Request.Method)
		}
	}
}

func Init(_ context.Context, config config.Config, newLogRepository repositories.IRestLogRepository) {
	logRepository = newLogRepository

	// Set gin mode
	cfg = config.GetObject()
	gin.SetMode(cfg.App.Mode)
	Router = gin.New()

	// List of mandatory middleware
	Router.Use(
		extractTraceParent,
		otelgin.Middleware(cfg.App.Name),
		logging,
		enableCors(),
		interceptResponse,
		gin.Recovery(),
		injectTraceParent,
		injectAuthContext,
	)

	// Handle no route error
	Router.NoRoute(func(c *gin.Context) {
		log.Warnf(c.Request.Context(), response.APINotRegistered.Description(), c.Request.URL.Path, c.Request.Method)
		SetResponse(c.Writer, response.APINotRegistered, c.Request.RequestURI, c.Request.Method)
	})

	// Set health check endpoint
	Router.GET("/v1/health", gin.WrapH(healthcheck.HandlerV1()))
}

func Run() {
	ctx := context.Background()

	// Skip if disabled
	if cfg.Transport.Server.Rest.Disabled {
		log.Warn(ctx, "REST server is disabled")
		return
	}

	port := cfg.Transport.Server.Rest.Port.Http
	log.Infof(ctx, "Starting REST server on port %d", port)
	err := Router.Run(fmt.Sprintf("%s:%d", "0.0.0.0", port))
	if err != nil {
		log.Fatalf(ctx, err, "Failed to run REST server, port=%d", port)
	}
}
