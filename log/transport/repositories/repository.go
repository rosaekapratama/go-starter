package repositories

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/location"
	myContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/log/constant"
	"github.com/rosaekapratama/go-starter/log/transport/models"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const (
	spanSave             = "repository.logging.database.Save"
	spanSaveRestRequest  = "repository.logging.database.SaveRestRequest"
	spanSaveRestResponse = "repository.logging.database.SaveRestResponse"
	spanSaveRestError    = "repository.logging.database.SaveRestError"
)

type ITransportLogRepository interface {
	Save(ctx context.Context, log *models.TransportLog)
	SavePubSubRequest(message pubsub.Message, isSubscriber bool)
	SavePubSubResponse(message pubsub.Message, isSubscriber bool)
	SavePubSubError(message pubsub.Message, isSubscriber bool, err error)
	SaveRestRequest(req *resty.Request, isServer bool)
	SaveRestResponse(res *resty.Response, isServer bool)
	SaveRestError(req *resty.Request, res *resty.Response, isServer bool, err error)
	SaveSoapRequest(req *http.Request, soapAction string, isServer bool)
	SaveSoapResponse(res *http.Response, soapAction string, isServer bool)
	SaveSoapError(req *http.Request, soapAction string, isServer bool, err error)
}

type TransportLogRepository struct {
	DB *gorm.DB
}

func (repo *TransportLogRepository) SavePubSubRequest(message pubsub.Message, isSubscriber bool) {
	//TODO implement me
}

func (repo *TransportLogRepository) SavePubSubResponse(message pubsub.Message, isSubscriber bool) {
	//TODO implement me
}

func (repo *TransportLogRepository) SavePubSubError(message pubsub.Message, isSubscriber bool, err error) {
	//TODO implement me
}

func (repo *TransportLogRepository) Save(ctx context.Context, transportLog *models.TransportLog) {
	go func(ctx context.Context, transportLog *models.TransportLog) {
		ctx, span := otel.Trace(ctx, spanSave)
		defer span.End()

		result := repo.DB.WithContext(ctx).Create(transportLog)
		if result.Error != nil {
			log.Error(ctx, result.Error, "failed to create transport log")
		}
	}(myContext.NewContextFromTraceParent(ctx), transportLog)
}

func (repo *TransportLogRepository) SaveRestRequest(req *resty.Request, isServer bool) {
	go func(ctx context.Context, req *resty.Request, isServer bool) {
		ctx, span := otel.Trace(ctx, spanSaveRestRequest)
		defer span.End()

		headers := fmt.Sprintf("%v", req.Header)
		var body *string
		if req.Body != nil {
			body = utils.StringP(fmt.Sprintf("%v", req.Body))
		}

		restLog, err := json.Marshal(&models.TransportRestLog{
			IsServer:  isServer,
			IsRequest: true,
			URL:       req.URL,
			Method:    req.Method,
			Headers:   &headers,
			Body:      body,
		})
		if err != nil {
			log.Error(ctx, err, "failed to marshal REST request log to json")
			return
		}

		transportLog := &models.TransportLog{
			ID:        uuid.New(),
			TraceID:   myContext.TraceIdFromContext(ctx),
			SpanID:    myContext.SpanIdFromContext(ctx),
			Type:      constant.LogTypeRest,
			Log:       datatypes.JSON(restLog),
			ProcessDT: time.Now().In(location.AsiaJakarta),
			ProcessBy: config.Instance.GetObject().App.Name,
		}
		result := repo.DB.WithContext(ctx).Create(transportLog)
		if result.Error != nil {
			log.Error(ctx, result.Error, "failed to create REST request log")
		}
	}(myContext.NewContextFromTraceParent(req.Context()), req, isServer)
}

func (repo *TransportLogRepository) SaveRestResponse(res *resty.Response, isServer bool) {
	go func(ctx context.Context, req *resty.Response, isServer bool) {
		ctx, span := otel.Trace(ctx, spanSaveRestResponse)
		defer span.End()

		headers := fmt.Sprintf("%v", res.Header())
		var restErr interface{}
		if res.IsError() {
			restErr = res.Error()
		}
		var body *string
		if res.Body() != nil && len(res.Body()) > integer.Zero {
			body = utils.StringP(string(res.Body()))
		}

		restLog, marshalErr := json.Marshal(&models.TransportRestLog{
			IsServer:   isServer,
			IsRequest:  false,
			URL:        res.Request.URL,
			Method:     res.Request.Method,
			Headers:    &headers,
			Body:       body,
			StatusCode: utils.StringP(res.Status()),
		})
		if marshalErr != nil {
			log.Error(ctx, marshalErr, "failed to marshal REST response log to json")
			return
		}

		transportLog := &models.TransportLog{
			ID:           uuid.New(),
			TraceID:      myContext.TraceIdFromContext(ctx),
			SpanID:       myContext.SpanIdFromContext(ctx),
			Type:         constant.LogTypeRest,
			Log:          datatypes.JSON(restLog),
			ErrorMessage: utils.StringP(fmt.Sprintf("%v", restErr)),
			ProcessDT:    time.Now().In(location.AsiaJakarta),
			ProcessBy:    config.Instance.GetObject().App.Name,
		}
		result := repo.DB.WithContext(ctx).Create(transportLog)
		if result.Error != nil {
			log.Error(ctx, result.Error, "failed to create REST response log")
		}
	}(myContext.NewContextFromTraceParent(res.Request.Context()), res, isServer)
}

func (repo *TransportLogRepository) SaveRestError(req *resty.Request, res *resty.Response, isServer bool, err error) {
	go func(ctx context.Context, req *resty.Request, res *resty.Response, isServer bool, err error) {
		ctx, span := otel.Trace(ctx, spanSaveRestError)
		defer span.End()

		headers := fmt.Sprintf("%v", req.Header)
		if res != nil {
			headers = fmt.Sprintf("%v", res.Header())
		}
		var body *string
		if req.Body != nil {
			body = utils.StringP(fmt.Sprintf("%v", req.Body))
		} else if res != nil && res.Body() != nil && len(res.Body()) > integer.Zero {
			body = utils.StringP(string(res.Body()))
		}

		restLog, marshalErr := json.Marshal(&models.TransportRestLog{
			IsServer:  isServer,
			IsRequest: false,
			URL:       req.URL,
			Method:    req.Method,
			Headers:   &headers,
			Body:      body,
		})
		if marshalErr != nil {
			log.Error(ctx, err, "failed to marshal REST error log to json")
			return
		}

		transportLog := &models.TransportLog{
			ID:           uuid.New(),
			TraceID:      myContext.TraceIdFromContext(ctx),
			SpanID:       myContext.SpanIdFromContext(ctx),
			Type:         constant.LogTypeRest,
			Log:          datatypes.JSON(restLog),
			ErrorMessage: utils.StringP(err.Error()),
			ProcessDT:    time.Now().In(location.AsiaJakarta),
			ProcessBy:    config.Instance.GetObject().App.Name,
		}
		result := repo.DB.WithContext(ctx).Create(transportLog)
		if result.Error != nil {
			log.Error(ctx, result.Error, "failed to create REST error log")
		}
	}(myContext.NewContextFromTraceParent(req.Context()), req, res, isServer, err)
}

func (repo *TransportLogRepository) SaveSoapRequest(req *http.Request, soapAction string, isServer bool) {
	//TODO implement me
}

func (repo *TransportLogRepository) SaveSoapResponse(res *http.Response, soapAction string, isServer bool) {
	//TODO implement me
}

func (repo *TransportLogRepository) SaveSoapError(req *http.Request, soapAction string, isServer bool, err error) {
	//TODO implement me
}

func NewTransportLogRepository(db *gorm.DB) ITransportLogRepository {
	return &TransportLogRepository{DB: db}
}
