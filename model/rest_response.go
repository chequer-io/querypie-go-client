package model

import "github.com/go-resty/resty/v2"

type RestResponse interface {
	SetHttpResponse(response *resty.Response)
	GetHttpResponse() *resty.Response
}
