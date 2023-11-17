package repositories

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/location"
	myContext "github.com/rosaekapratama/go-starter/context"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/transport/logging/models"
	"github.com/rosaekapratama/go-starter/utils"
	"gorm.io/gorm"
	"time"
)

const (
	spanSave         = "repository.logging.database.Save"
	spanSaveRequest  = "repository.logging.database.SaveRequest"
	spanSaveResponse = "repository.logging.database.SaveResponse"
	spanSaveError    = "repository.logging.database.SaveError"
)

type IRestLogRepository interface {
	Save(ctx context.Context, log *models.TransportRestLog) error
	SaveRequest(req *resty.Request, isServer bool) error
	SaveResponse(res *resty.Response, isServer bool) error
	SaveError(req *resty.Request, res *resty.Response, isServer bool, err error) error
}

type RestLogRepository struct {
	DB *gorm.DB
}

func (repo *RestLogRepository) Save(ctx context.Context, log *models.TransportRestLog) error {
	ctx, span := otel.Trace(ctx, spanSave)
	defer span.End()

	result := repo.DB.WithContext(ctx).Create(log)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RestLogRepository) SaveRequest(req *resty.Request, isServer bool) error {
	ctx := req.Context()
	ctx, span := otel.Trace(ctx, spanSaveRequest)
	defer span.End()

	headers := fmt.Sprintf("%v", req.Header)
	var body *string
	if req.Body != nil {
		body = utils.StringP(fmt.Sprintf("%v", req.Body))
	}
	log := &models.TransportRestLog{
		ID:        uuid.New(),
		TraceID:   myContext.TraceIdFromContext(ctx),
		SpanID:    myContext.SpanIdFromContext(ctx),
		IsServer:  isServer,
		IsRequest: true,
		URL:       req.URL,
		Method:    req.Method,
		Headers:   &headers,
		Body:      body,
		ProcessDT: time.Now().In(location.AsiaJakarta),
		ProcessBy: config.Instance.GetObject().App.Name,
	}
	result := repo.DB.WithContext(ctx).Create(log)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RestLogRepository) SaveResponse(res *resty.Response, isServer bool) error {
	ctx := res.Request.Context()
	ctx, span := otel.Trace(ctx, spanSaveResponse)
	defer span.End()

	headers := fmt.Sprintf("%v", res.Header())
	var err interface{}
	if res.IsError() {
		err = res.Error()
	}
	var body *string
	if res.Body() != nil && len(res.Body()) > integer.Zero {
		body = utils.StringP(string(res.Body()))
	}
	log := &models.TransportRestLog{
		ID:           uuid.New(),
		TraceID:      myContext.TraceIdFromContext(ctx),
		SpanID:       myContext.SpanIdFromContext(ctx),
		IsServer:     isServer,
		IsRequest:    false,
		URL:          res.Request.URL,
		Method:       res.Request.Method,
		Headers:      &headers,
		Body:         body,
		StatusCode:   utils.StringP(res.Status()),
		ErrorMessage: utils.StringP(fmt.Sprintf("%v", err)),
		ProcessDT:    time.Now().In(location.AsiaJakarta),
		ProcessBy:    config.Instance.GetObject().App.Name,
	}
	result := repo.DB.WithContext(ctx).Create(log)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RestLogRepository) SaveError(req *resty.Request, res *resty.Response, isServer bool, err error) error {
	ctx := req.Context()
	ctx, span := otel.Trace(ctx, spanSaveError)
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
	log := &models.TransportRestLog{
		ID:           uuid.New(),
		TraceID:      myContext.TraceIdFromContext(ctx),
		SpanID:       myContext.SpanIdFromContext(ctx),
		IsServer:     isServer,
		IsRequest:    false,
		URL:          req.URL,
		Method:       req.Method,
		Headers:      &headers,
		Body:         body,
		ErrorMessage: utils.StringP(err.Error()),
		ProcessDT:    time.Now().In(location.AsiaJakarta),
		ProcessBy:    config.Instance.GetObject().App.Name,
	}
	result := repo.DB.WithContext(ctx).Create(log)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewRestLogRepository(db *gorm.DB) IRestLogRepository {
	return &RestLogRepository{DB: db}
}
