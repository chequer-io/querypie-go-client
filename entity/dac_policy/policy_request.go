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
	PolicyUuid       string     `json:"-" yaml:"policyUuid"`
}

func (pr PolicyRequest) UpdateOrCreateRemotely(server utils.QueryPieServerConfig) *Policy {
	if len(pr.PolicyUuid) > 0 {
		return pr.UpdateRemotely(server)
	} else {
		return pr.CreateRemotely(server)
	}
}

func (pr PolicyRequest) CreateRemotely(server utils.QueryPieServerConfig) *Policy {
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
		logrus.Fatalf("Failed to create a policy: %v", err)
	}
	response.HttpResponse = httpResponse
	return &response
}

func (pr PolicyRequest) UpdateRemotely(server utils.QueryPieServerConfig) *Policy {
	var response Policy

	restClient := resty.New()
	uri := fmt.Sprintf("%s/api/external/policies/%s", server.BaseURL, pr.PolicyUuid)
	httpResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(server.AccessToken).
		SetBody(pr).
		SetResult(&response).
		Put(uri)
	logrus.Debugf("Response: %v", httpResponse)
	if err != nil {
		logrus.Fatalf("Failed to update a policy: %v", err)
	}
	response.HttpResponse = httpResponse
	return &response
}

func (pr PolicyRequest) DeleteRemotely(server utils.QueryPieServerConfig) *Policy {
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
	response.Uuid = pr.PolicyUuid
	response.HttpResponse = httpResponse
	return &response
}
