package utils

import (
	"errors"
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

		if IsClientError(resp) || IsServerError(resp) {
			PrintHttpRequestLineAndResponseStatus(resp)
			logrus.Fatalf("Failed to fetch resources: %v", resp.String())
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
	result.SetHttpResponse(response)
	return result, err
}

func IsClientError(r *resty.Response) bool {
	return 399 < r.StatusCode() && r.StatusCode() < 500
}

func IsServerError(r *resty.Response) bool {
	return 499 < r.StatusCode() && r.StatusCode() < 600
}

func PrintHttpRequestLineAndResponseStatus(r *resty.Response) {
	req := r.Request.RawRequest
	res := r.RawResponse
	fmt.Printf("%s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
	fmt.Printf("%s %s\n\n", res.Proto, res.Status)
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

func FindMultiple[T any](
	items *[]T,
	queryFunc func(db *gorm.DB) *gorm.DB,
) {
	result := queryFunc(config.LocalDatabase)

	if result.Error == nil {
		if result.RowsAffected == 0 {
			logrus.Debugf("No rows found")
		} else {
			logrus.Debugf("Found %d rows: %v", result.RowsAffected, &(*items)[0])
		}
	} else if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logrus.Debugf("No rows found: %v", result.Error)
		} else {
			logrus.Fatalf("Failed to find items by error: %s", result.Error)
		}
	}
}

func First[T any](
	item *T,
	queryFunc func(db *gorm.DB) *gorm.DB,
) *T {
	result := queryFunc(config.LocalDatabase)

	if result.Error == nil {
		if result.RowsAffected == 0 {
			logrus.Debugf("No rows found")
		} else {
			logrus.Debugf("Found %d rows: %v", result.RowsAffected, item)
			return item
		}
	} else if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logrus.Debugf("No rows found: %s", result.Error)
		} else {
			logrus.Fatalf("Failed to find items by error: %s", result.Error)
		}
	}
	return nil
}
