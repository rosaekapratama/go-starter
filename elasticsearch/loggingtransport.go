package elasticsearch

import (
	"bytes"
	"context"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/transport/restclient"
	"io"
	"net/http"
)

func NewLoggingTransport(ctx context.Context, isStdoutLogEnabled bool, databaseLog string, payloadLogSizeLimit int) http.RoundTripper {
	client, err := restclient.Manager.NewClient(ctx, restclient.WithLogging(isStdoutLogEnabled, databaseLog))
	if err != nil {
		log.Fatalf(ctx, err, "error on create rest client for elastic search")
	}

	return &loggingTransport{
		restClient:          client,
		isStdoutLogEnabled:  isStdoutLogEnabled,
		payloadLogSizeLimit: payloadLogSizeLimit,
	}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	ctx := req.Context()
	restyReq := t.restClient.NewRequest(ctx)

	// Copy header
	for key, values := range req.Header {
		for _, value := range values {
			restyReq.SetHeader(key, value)
		}
	}

	// Set the body
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body for further use
		restyReq.SetBody(bodyBytes)
	}

	// Perform the request
	restyRes, err := restyReq.Execute(req.Method, req.URL.String())
	if err != nil {
		log.Error(ctx, err)
		return
	}

	// Create http.Response from resty.Response
	res = &http.Response{
		Status:        restyRes.Status(),
		StatusCode:    restyRes.StatusCode(),
		Header:        restyRes.Header(),
		Body:          io.NopCloser(bytes.NewBuffer(restyRes.Body())),
		ContentLength: int64(len(restyRes.Body())),
		Request:       req,
	}

	return
}
