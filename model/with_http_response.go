package model

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
)

type WithHttpResponseInterface interface {
	GetHttpResponse() *resty.Response
	SetHttpResponse(*resty.Response)
	HandleResponse(
		onSuccess func(),
		onClientFailure func(),
		onServerFailure func(),
	)
	PrintHttpRequestLineAndResponseStatus(silent bool) *WithHttpResponse
	PrintRawBody(silent bool) *WithHttpResponse
}

type WithHttpResponse struct {
	// Internal: HTTP response
	HttpResponse *resty.Response `json:"-" gorm:"-" yaml:"-"`
}

func (r *WithHttpResponse) GetHttpResponse() *resty.Response {
	return r.HttpResponse
}

func (r *WithHttpResponse) SetHttpResponse(response *resty.Response) {
	r.HttpResponse = response
}

func (r *WithHttpResponse) HandleResponse(
	onSuccess func(),
	onClientFailure func(),
	onServerFailure func(),
) *WithHttpResponse {
	if r.HttpResponse == nil {
		// Do nothing.
		return r
	}

	if r.HttpResponse.IsSuccess() {
		onSuccess()
	} else if r.isClientFailure() {
		onClientFailure()
	} else if r.isServerFailure() {
		onServerFailure()
	} else {
		logrus.Fatalf("Unexpected status code: %d", r.HttpResponse.StatusCode())
	}
	return r
}

func (r *WithHttpResponse) PrintHttpRequestLineAndResponseStatus(silent bool) *WithHttpResponse {
	if silent {
		return r
	}
	req := r.HttpResponse.Request.RawRequest
	res := r.HttpResponse.RawResponse
	fmt.Printf("%s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
	fmt.Printf("%s %s\n\n", res.Proto, res.Status)
	return r
}

func (r *WithHttpResponse) PrintRawBody(silent bool) *WithHttpResponse {
	if silent {
		return r
	}
	fmt.Printf("%s\n",
		pretty.Pretty(r.HttpResponse.Body()),
	)
	return r
}

func (r *WithHttpResponse) isClientFailure() bool {
	code := r.HttpResponse.StatusCode()
	return 400 <= code && code <= 499
}

func (r *WithHttpResponse) isServerFailure() bool {
	code := r.HttpResponse.StatusCode()
	return 500 <= code && code <= 599
}
