package restserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inhies/go-bytesize"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/headers"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/location"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	commonContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/healthcheck"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/log/transport/models"
	"github.com/rosaekapratama/go-starter/log/transport/repositories"
	"github.com/rosaekapratama/go-starter/response"
	"github.com/rosaekapratama/go-starter/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/datatypes"
)

const (
	contentTypeApplicationJson = "application/json"
	realmsPath                 = "realms/"
)

var (
	cfg           *config.Object
	propagator    = otel.GetTextMapPropagator()
	Router        *gin.Engine
	logRepository repositories.ITransportLogRepository
)

func logging(payloadLogSizeLimit int) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Skip if request is health check
		isHealthCheck, _ := regexp.MatchString(healthcheck.URLPathRegex, c.Request.URL.Path)
		if isHealthCheck {
			return
		}

		ctx := c.Request.Context()
		isStdoutLogEnabled := cfg.Transport.Server.Rest.Logging.Stdout
		databaseLog := cfg.Transport.Server.Rest.Logging.Database

		var clonedReq *http.Request
		var clearFunc func()
		var body string
		var err error

		// http request clone process
		if isStdoutLogEnabled || databaseLog != str.Empty {
			clonedReq, clearFunc, err = utils.CloneHttpRequest(c.Request, payloadLogSizeLimit)
			defer clearFunc()
			if err != nil {
				log.Error(ctx, err, "Failed to clone http request for logging")
				c.Next()
				return
			}

			bytes, errIf := io.ReadAll(clonedReq.Body)
			if errIf != nil {
				err = errIf
				log.Error(ctx, err, "Failed to read cloned body request for logging")
				c.Next()
				return
			}

			body = string(bytes)
		}

		// write to stdout
		if isStdoutLogEnabled {
			log.Trace(ctx, "Start stdout http request log writing")
			httpFields := make(map[string]interface{})
			httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeRest
			httpFields[constant.UrlLogKey] = clonedReq.URL
			httpFields[constant.MethodLogKey] = clonedReq.Method
			httpFields[constant.IsServerLogKey] = true
			httpFields[constant.IsRequestLogKey] = true
			httpFields[constant.HeadersLogKey] = clonedReq.Header
			httpFields[constant.BodyLogKey] = body
			log.WithTraceFields(clonedReq.Context()).WithFields(httpFields).GetLogrusLogger().Info()
			log.Trace(ctx, "End stdout http request log writing")
		}

		// write to database
		if databaseLog != str.Empty {
			go func(ctx context.Context) {
				log.Trace(ctx, "Start database http request log writing")
				restLog, err := json.Marshal(&models.TransportRestLog{
					IsServer:   true,
					IsRequest:  true,
					URL:        fmt.Sprintf("%v", clonedReq.URL),
					Method:     clonedReq.Method,
					Headers:    utils.StringP(fmt.Sprintf("%v", clonedReq.Header)),
					Body:       utils.StringP(body),
					StatusCode: nil,
				})
				if err != nil {
					log.Error(ctx, err, "failed to marshal REST server request log to json")
					return
				}

				transportLog := &models.TransportLog{
					ID:        uuid.New(),
					TraceID:   commonContext.TraceIdFromContext(ctx),
					SpanID:    commonContext.SpanIdFromContext(ctx),
					Type:      constant.LogTypeRest,
					Log:       datatypes.JSON(restLog),
					ProcessDT: time.Now().In(location.AsiaJakarta),
					ProcessBy: config.Instance.GetObject().App.Name,
				}
				logRepository.Save(ctx, transportLog)
				log.Trace(ctx, "End database http request log writing")
			}(commonContext.NewContextFromTraceParent(ctx))
		}

		c.Next()
		i := c.Writer.(*WriterInterceptor)

		// write to stdout
		if isStdoutLogEnabled {
			log.Trace(ctx, "Start stdout http response log writing")
			httpFields := make(map[string]interface{})
			httpFields[constant.LogTypeFieldLogKey] = constant.LogTypeRest
			httpFields[constant.UrlLogKey] = clonedReq.URL
			httpFields[constant.MethodLogKey] = clonedReq.Method
			httpFields[constant.IsServerLogKey] = true
			httpFields[constant.IsRequestLogKey] = false
			httpFields[constant.StatusCodeLogKey] = i.status
			httpFields[constant.HeadersLogKey] = i.ResponseWriter.Header()
			httpFields[constant.BodyLogKey] = string(i.body)
			log.WithTraceFields(clonedReq.Context()).WithFields(httpFields).GetLogrusLogger().Info()
			log.Trace(ctx, "End stdout http response log writing")
		}

		if databaseLog != str.Empty {
			if len(i.body) > integer.Zero {
				body = string(i.body)
			}

			go func(ctx context.Context) {
				log.Trace(ctx, "Start database http response log writing")
				restLog, marshalErr := json.Marshal(&models.TransportRestLog{
					IsServer:   true,
					IsRequest:  false,
					URL:        fmt.Sprintf("%v", clonedReq.URL),
					Method:     clonedReq.Method,
					Headers:    utils.StringP(fmt.Sprintf("%v", i.ResponseWriter.Header())),
					Body:       utils.StringP(body),
					StatusCode: utils.StringP(strconv.Itoa(i.status)),
				})
				if marshalErr != nil {
					log.Error(ctx, marshalErr, "failed to marshal REST response log to json")
					return
				}

				transportLog := &models.TransportLog{
					ID:           uuid.New(),
					TraceID:      commonContext.TraceIdFromContext(ctx),
					SpanID:       commonContext.SpanIdFromContext(ctx),
					Type:         constant.LogTypeRest,
					Log:          datatypes.JSON(restLog),
					ErrorMessage: utils.StringP(i.response.Description()),
					ProcessDT:    time.Now().In(location.AsiaJakarta),
					ProcessBy:    config.Instance.GetObject().App.Name,
				}
				logRepository.Save(ctx, transportLog)
				log.Trace(ctx, "End database http response log writing")
			}(commonContext.NewContextFromTraceParent(ctx))
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
		log.Debugf(ctx, "CORS allow credentials : %v", corsConfig.AllowCredentials)
		log.Debugf(ctx, "CORS allow origins     : %v", corsConfig.AllowOrigins)
		log.Debugf(ctx, "CORS allow methods     : %v", corsConfig.AllowMethods)
		log.Debugf(ctx, "CORS allow headers     : %v", corsConfig.AllowHeaders)
		log.Debugf(ctx, "CORS expose headers    : %v", corsConfig.ExposeHeaders)
		log.Debugf(ctx, "CORS max age           : %v second", propsConfig.MaxAge)
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
		c.Request = c.Request.WithContext(commonContext.ContextWithTraceParent(c.Request.Context(), traceparent))
	}
}

func interceptResponse(payloadLogSizeLimit int) func(c *gin.Context) {
	return func(c *gin.Context) {
		r := c.Request
		i := NewWriterInterceptor(r.Context(), c.Writer, payloadLogSizeLimit)
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
			ctx = commonContext.ContextWithUsername(ctx, username)

			// Override request context with new context
			c.Request = c.Request.WithContext(ctx)
		}
	} else if auth != str.Empty && strings.HasPrefix(auth, headers.BearerTokenPrefix) {
		token := strings.TrimPrefix(auth, headers.BearerTokenPrefix)
		ctx = commonContext.ContextWithToken(ctx, token)

		// Override request context with new context
		c.Request = c.Request.WithContext(ctx)
	} else {
		log.Tracef(ctx, "Authorization is empty, skip inject auth context process, path=%s, method=%s", c.Request.URL.Path, c.Request.Method)
	}
}

func Init(ctx context.Context, config config.Config, newLogRepository repositories.ITransportLogRepository) {
	logRepository = newLogRepository

	// Set gin mode
	cfg = config.GetObject()
	gin.SetMode(cfg.App.Mode)
	Router = gin.New()

	// Get rest server payload limit config for logging
	_payloadLogSizeLimit, err := bytesize.Parse(config.GetObject().Transport.Server.Rest.Logging.PayloadLogSizeLimit)
	if err != nil {
		log.Fatal(ctx, err, "Invalid value of REST server payloadLogSizeLimit config")
	}
	payloadLogSizeLimit := int(_payloadLogSizeLimit)

	// List of mandatory middleware
	Router.Use(
		extractTraceParent,
		otelgin.Middleware(cfg.App.Name),
		logging(payloadLogSizeLimit),
		enableCors(),
		interceptResponse(payloadLogSizeLimit),
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
