package server

import (
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/page"
	"github.com/rosaekapratama/go-starter/response"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

func GetPageFromRequest(r *http.Request) (*page.PageRequest, error) {
	ctx := r.Context()
	var err error

	pageNumInt := page.DefaultPageNum
	pageSizeInt := page.DefaultPageSize

	queryMap := r.URL.Query()

	log.Trace(ctx, "Trying to extract page number from query param...")
	pageNum := queryMap[page.PageNumQueryKey]
	if len(pageNum) > integer.Zero && pageNum[integer.Zero] != str.Empty {
		log.Tracef(ctx, "Found page number in query param, pageNum=%s", pageNum[integer.Zero])
		pageNumInt, err = strconv.Atoi(pageNum[integer.Zero])
		if err != nil {
			log.Error(ctx, err, "Convert page number string to int failed")
			return nil, err
		}
	} else {
		log.Trace(ctx, "Page number is not found in query param, trying to extract from headers...")
		pageNum := r.Header.Get(page.PageNumHeaderKey)
		if pageNum != str.Empty {
			log.Tracef(ctx, "Found page number in headers, pageNum=%s", pageNum)
			pageNumInt, err = strconv.Atoi(pageNum)
			if err != nil {
				log.Error(ctx, err, "Convert page number string to int failed")
				return nil, err
			}
		}
	}

	log.Trace(ctx, "Trying to extract page size from query param...")
	pageSize := queryMap[page.PageSizeQueryKey]
	if len(pageSize) > integer.Zero && pageSize[integer.Zero] != str.Empty {
		log.Tracef(ctx, "Found page size in query param, pageSize=%s", pageSize[integer.Zero])
		pageSizeInt, err = strconv.Atoi(pageSize[integer.Zero])
		if err != nil {
			log.Error(ctx, err, "Convert page size string to int failed")
			return nil, err
		}
	} else {
		log.Trace(ctx, "Page size is not found in query param, trying to extract from headers...")
		pageSize := r.Header.Get(page.PageSizeHeaderKey)
		if pageSize != str.Empty {
			log.Tracef(ctx, "Found page size in headers, pageSize=%s", pageSize)
			pageSizeInt, err = strconv.Atoi(pageSize)
			if err != nil {
				log.Error(ctx, err, "Convert page size string to int failed")
				return nil, err
			}
		}
	}

	pageRequest := page.NewPageRequest(pageNumInt, pageSizeInt)
	if !pageRequest.IsValid() {
		log.Trace(ctx, response.InvalidPageRequest, "Invalid page request")
		return nil, response.InvalidPageRequest
	}

	return pageRequest, nil
}

// SetResponse Set response code and description to writer interceptor
func SetResponse(w http.ResponseWriter, err error, a ...any) {
	i := castInterceptor(w)
	i.a = a
	span := trace.SpanFromContext(i.ctx)
	switch t := err.(type) {
	case response.IResponse:
		if span != nil {
			span.SetStatus(t.OtelCode(), t.Description())
		}
		i.response = t
	default:
		res := response.GeneralError
		if span != nil {
			span.SetStatus(res.OtelCode(), err.Error())
		}
		i.response = res
	}
}

// SetRawResponse same as SetResponse but without BaseResponse json format
func SetRawResponse(w http.ResponseWriter, err error) {
	i := castInterceptor(w)
	i.IsRaw = true
	span := trace.SpanFromContext(i.ctx)
	switch t := err.(type) {
	case response.IResponse:
		if span != nil {
			span.SetStatus(t.OtelCode(), t.Description())
		}
		i.response = t
	default:
		res := response.GeneralError
		if span != nil {
			span.SetStatus(res.OtelCode(), err.Error())
		}
		i.response = res
	}
}

func SetPagination(w http.ResponseWriter, p *page.PageResponse) {
	i := castInterceptor(w)
	i.page = p
}
