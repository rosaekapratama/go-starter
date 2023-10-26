package restserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rosaekapratama/go-starter/avro"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	myContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/database"
	"github.com/rosaekapratama/go-starter/elasticsearch"
	"github.com/rosaekapratama/go-starter/google"
	"github.com/rosaekapratama/go-starter/google/cloud/oauth"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub"
	"github.com/rosaekapratama/go-starter/google/cloud/pubsub/sub"
	"github.com/rosaekapratama/go-starter/google/cloud/scheduler"
	"github.com/rosaekapratama/go-starter/google/cloud/storage"
	"github.com/rosaekapratama/go-starter/google/firebase"
	"github.com/rosaekapratama/go-starter/healthcheck"
	"github.com/rosaekapratama/go-starter/keycloak"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/loginit"
	"github.com/rosaekapratama/go-starter/mocks"
	myOtel "github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/redis"
	"github.com/rosaekapratama/go-starter/response"
	"github.com/rosaekapratama/go-starter/transport/constant"
	"github.com/rosaekapratama/go-starter/transport/restclient"
	"github.com/rosaekapratama/go-starter/zeebe"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	contentTypeApplicationJson = "application/json"
)

var (
	cfg        *config.Object
	propagator = otel.GetTextMapPropagator()
	Router     *gin.Engine
)

func logging(c *gin.Context) {
	r := c.Request
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
	if r.Body != nil {
		httpFields[constant.BodyFieldKey] = r.Body
	}
	log.WithTraceFields(r.Context()).WithFields(httpFields).GetLogrusLogger().Info()

	c.Next()

	i := c.Writer.(*WriterInterceptor)
	httpFields = make(map[string]interface{})
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
		bytes, err := base64.StdEncoding.DecodeString(auth)
		if err != nil {
			log.Error(ctx, err, "Failed to decode basic auth")
		} else {
			// Get the username from decoded value and set to context
			username := strings.Split(string(bytes), sym.Colon)[integer.Zero]
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
					c := v.(string)
					if c == str.Empty {
						log.Tracef(ctx, "Unable to set context, %s claim is empty", additionalClaim)
					} else {
						ctx = context.WithValue(ctx, additionalClaim, c)
					}
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

func initRun(ctx context.Context) {
	// Init config package
	config.Init()
	config := config.Instance
	cfg = config.GetObject()

	// Init google credential
	projectId := cfg.App.Mode
	credentials, jsonKey := google.CreateCredentials(ctx, config)

	// Extract project ID from credentials
	if credentials != nil && credentials.ProjectID != str.Empty {
		projectId = credentials.ProjectID
	}

	// Set project ID for loginit
	loginit.SetProjectId(projectId)
	log.Infof(ctx, "projectId=%s", projectId)

	// Init google package
	if credentials != nil {
		firebaseApp := firebase.New(ctx, credentials)
		oauthClient := oauth.New(ctx)
		pubsubClient := pubsub.New(ctx, credentials)
		sub.Init(pubsubClient)
		scheduler.Init(ctx, credentials)
		schedulerService := scheduler.Service
		storage.Init(ctx, credentials)
		storageClient := storage.Client
		google.Init(
			ctx,
			credentials,
			jsonKey,
			firebaseApp,
			oauthClient,
			pubsubClient,
			schedulerService,
			storageClient)
	}

	// Init log package
	log.Init(ctx, config, projectId)

	// Init otel package
	myOtel.Init(ctx, config)

	// Init avro package
	avro.Init(ctx)

	// Init database package
	database.Init(ctx, config)

	// Init elasticsearch package
	elasticsearch.Init(ctx, config)

	// Init redis package
	redis.Init(ctx, config)

	// Init http client package
	restclient.Init(ctx, config)

	// Init zeebe package
	zeebe.Init(ctx, config)

	// Set gin mode
	gin.SetMode(cfg.App.Mode)
	Router = gin.New()

	// List of mandatory middleware
	Router.Use(
		logging,
		enableCors(),
		extractTraceParent,
		interceptResponse,
		gin.Recovery(),
		otelgin.Middleware(cfg.App.Name),
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

func initTest(_ context.Context) {
	// To handle init() function which calls config
	mockConfig := mocks.GetMockConfig()
	mockConfig.On("GetString", mock.Anything).Return("string", nil)
	mockConfig.On("GetInt", mock.Anything).Return(integer.Zero, nil)
	mockConfig.On("GetBool", mock.Anything).Return(false, nil)
	mockConfig.On("GetSlice", mock.Anything).Return(make([]interface{}, integer.Zero), nil)
	mockConfig.On("GetStringAndThrowFatalIfEmpty", mock.Anything).Return("string", nil)
	config.Instance = mockConfig
}

func init() {
	ctx := context.Background()

	args := os.Args
	if strings.HasSuffix(args[0], ".test") {
		initTest(ctx)
	} else {
		initRun(ctx)
	}
}

func Run() {
	ctx := context.Background()

	// run the http server
	err := Router.Run(fmt.Sprintf("%s:%d", "0.0.0.0", cfg.Transport.Server.Rest.Port.Http))
	if err != nil {
		log.Fatal(ctx, err, "Failed to start server")
	}
}
