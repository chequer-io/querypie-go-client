package dac_policy

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"qpc/utils"
)

type PolicyRequest struct {
	ClusterGroupUuid string     `json:"clusterGroupUuid" yaml:"clusterGroupUuid"`
	PolicyType       PolicyType `json:"policyType" yaml:"policyType"`
	Title            string     `json:"title" yaml:"title"`
	PolicyUuid       string     `json:"-" yaml:"-"`
}

func (pr PolicyRequest) CreateByPost(server utils.QueryPieServerConfig) *Policy {
	var response Policy

	restClient := resty.New()
	uri := fmt.Sprintf("%s/api/external/policies", server.BaseURL)
	httpResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(server.AccessToken).
		SetBody(pr).
		SetResult(&response).
		Post(uri)
	logrus.Debugf("Response: %v", httpResponse)
	if err != nil {
		logrus.Fatalf("Failed to grant access to DAC connection: %v", err)
	}
	response.HttpResponse = httpResponse
	return &response
}

func (pr PolicyRequest) DeleteByDelete(server utils.QueryPieServerConfig) *Policy {
	var response Policy

	restClient := resty.New()
	uri := fmt.Sprintf("%s/api/external/policies/%s", server.BaseURL, pr.PolicyUuid)
	httpResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(server.AccessToken).
		SetResult(&response).
		Delete(uri)
	logrus.Debugf("Response: %v", httpResponse)
	if err != nil {
		logrus.Fatalf("Failed to delete policy: %v", err)
	}
	response.HttpResponse = httpResponse
	return &response
}
