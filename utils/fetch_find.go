package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"qpc/config"
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

func FindAllAndForEach[T any](
	countFunc func(db *gorm.DB, total *int64) *gorm.DB,
	findAllFunc func(db *gorm.DB, items *[]T) *gorm.DB,
	forEachFunc func(found *T) bool,
) {
	var items []T
	var total, selected int64 = 0, 0
	result := countFunc(config.LocalDatabase, &total)
	if result.Error != nil {
		logrus.Fatalf("Failed to count items: %v", result.Error)
	}
	logrus.Debugf("Found %d items", total)

	result = findAllFunc(config.LocalDatabase, &items)
	if result.Error != nil {
		logrus.Fatalf("Failed to select data from local database: %v", result.Error)
	}
	selected = int64(len(items))
	for i := range items {
		forEachFunc(&items[i])
	}
	if selected != total {
		logrus.Errorf("Selected %d, whereas total count was %d, difference: %d",
			selected, total, total-selected)
	}
}
