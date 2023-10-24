package oauth

import "context"

type Client interface {
	VerifyToken(ctx context.Context, token string) (response *InfotokenResponse, err error)
}

type ClientImpl struct {
	tokeninfoEndpoint string
}

type InfotokenResponse struct {
	Aud           string `json:"aud"`
	Azp           string `json:"azp"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Exp           string `json:"exp"`
	Iat           string `json:"iat"`
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	Typ           string `json:"typ"`
}
