package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"qpc/model"
)

func FetchPagedListAndForEach[T any, P model.PagedList[T]](
	uri string,
	result P,
	forEachFunc func(page P) bool,
) {
	page := 0
	size := 40 // Set the desired page size
	restClient := resty.New()

	logrus.Debugf("Type of result: %T", result)
	for {
		resp, err := restClient.R().
			SetQueryParams(
				map[string]string{
					"pageSize":   fmt.Sprintf("%d", size),
					"pageNumber": fmt.Sprintf("%d", page),
				},
			).
			SetHeader("Accept", "application/json").
			SetAuthToken(DefaultQuerypieServer.AccessToken).
			SetResult(&result).
			Get(DefaultQuerypieServer.BaseURL + uri)
		logrus.Debugf("Response: %v", resp)
		if err != nil {
			logrus.Fatalf("Failed to fetch resources: %v", err)
		}

		ok := forEachFunc(result)
		if !(ok && result.GetPage().HasNext()) {
			break
		}
		page++
	}
}

func Fetch[T model.RestResponse](
	uri string,
	result T,
) (T, error) {
	restClient := resty.New()
	response, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(DefaultQuerypieServer.AccessToken).
		SetResult(&result).
		Get(DefaultQuerypieServer.BaseURL + uri)
	logrus.Debugf("Response: %v", response)
	if err != nil {
		return result, err
	}
	result.SetHttpResponse(response)
	return result, nil
}

func FetchAndPrint[T model.RestResponse](
	uri string,
	result T,
	printFunc func(object T),
) (T, error) {
	restClient := resty.New()
	response, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(DefaultQuerypieServer.AccessToken).
		SetResult(&result).
		Get(DefaultQuerypieServer.BaseURL + uri)
	logrus.Debugf("Response: %v", response)
	if err != nil {
		logrus.Fatalf("Failed to fetch a resource: %v", err)
		return result, err
	}
	result.SetHttpResponse(response)
	printFunc(result)
	return result, nil
}
