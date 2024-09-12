package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"qpc/model"
)

func FetchPrintAndSave[T any, P model.PagedList[T]](
	uri string,
	result P,
	printFunc func(object P, first bool, last bool),
	saveFunc func(object P) bool,
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

		printFunc(result, page == 0, !result.GetPage().HasNext())
		shouldBreak := saveFunc(result)
		if shouldBreak {
			break
		}

		page++
	}
}
