package model

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
)

type WithHttpResponse struct {
	// Internal: HTTP response
	HttpResponse *resty.Response `json:"-" gorm:"-" yaml:"-"`
}

func (r WithHttpResponse) PrintHttpReqRes(silent bool, printFunc func()) WithHttpResponse {
	if silent {
		return r
	} else if r.HttpResponse == nil {
		printFunc()
		return r
	}

	printHttpRequestLineAndResponseStatus(r.HttpResponse)
	var printRawBody = false
	if r.HttpResponse.StatusCode() != 200 {
		printRawBody = true
	} else if r.HttpResponse.Request.Method == "DELETE" {
		printRawBody = true
	}

	if printRawBody {
		fmt.Printf("%s\n",
			pretty.Pretty(r.HttpResponse.Body()),
		)
	} else {
		printFunc()
	}
	return r
}

func printHttpRequestLineAndResponseStatus(r *resty.Response) {
	req := r.Request.RawRequest
	res := r.RawResponse
	fmt.Printf("%s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
	fmt.Printf("%s %s\n\n", res.Proto, res.Status)
}
