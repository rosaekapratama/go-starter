package client

import "github.com/go-resty/resty/v2"

type Client struct {
	Resty *resty.Client
}
