package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/transport/restclient"
)

const (
	errFailedToGetTokenInfo = "failed to get token info from google"

	spanVerifyToken   = "common.google.cloud.oauth.VerifyToken"
	tokeninfoEndpoint = "https://oauth2.googleapis.com/tokeninfo?id_token=%s"
)

func (c *clientImpl) VerifyToken(ctx context.Context, token string) (response *InfotokenResponse, err error) {
	ctx, span := otel.Trace(ctx, spanVerifyToken)
	defer span.End()

	var res *resty.Response
	req := restclient.Manager.GetDefaultClient().NewRequest(ctx)
	res, err = req.Get(fmt.Sprintf(c.tokeninfoEndpoint, token))
	if err != nil {
		log.Error(ctx, err, errFailedToGetTokenInfo)
		return
	}

	if res.IsSuccess() {
		response = &InfotokenResponse{}
		err = json.Unmarshal(res.Body(), response)
		if err != nil {
			log.Error(ctx, err)
			return nil, err
		}
		return
	} else {
		log.Tracef(ctx, "Failed to get token info from google, httpCode=%d, httpBody=%s", res.StatusCode(), string(res.Body()))
		return
	}
}

func NewClient(ctx context.Context) Client {
	log.Info(ctx, "Google OAuth client service is initiated")
	return &clientImpl{tokeninfoEndpoint: tokeninfoEndpoint}
}
