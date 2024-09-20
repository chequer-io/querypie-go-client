package dac_policy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"qpc/entity/dac_connection"
)

type PolicyRequest struct {
	Matched    []SummarizedConnectionForPolicy `json:"-" yaml:"matched"`
	Validation struct {
		Result bool   `json:"-" yaml:"result"`
		Reason string `json:"-" yaml:"reason,omitempty"`
	} `json:"-" yaml:"validation"`

	ConnectionUuid string     `json:"clusterGroupUuid" yaml:"-"`
	PolicyType     PolicyType `json:"policyType"`
	Name           string     `json:"title"`
}

func GeneratePolicyRequest(connection string, policyType PolicyType, name string) *PolicyRequest {
	var found []dac_connection.ConnectionV2
	if len(connection) > 0 {
		(&dac_connection.ConnectionV2{}).FindByNameOrUuid(connection, &found)
	}

	var matched []SummarizedConnectionForPolicy
	for _, it := range found {
		matched = append(matched, SummarizedConnectionForPolicy{
			Uuid:         it.Uuid,
			DatabaseType: it.DatabaseType,
			Name:         it.Name,
		})
	}
	return &PolicyRequest{
		Matched:    matched,
		PolicyType: policyType,
		Name:       name,
	}
}

func (pr *PolicyRequest) Validate() *PolicyRequest {
	pr.Validation.Result = false

	if len(pr.Matched) == 0 {
		pr.Validation.Reason = "Connection not found"
	} else if len(pr.Matched) == 1 {
		pr.ConnectionUuid = pr.Matched[0].Uuid
		pr.Validation.Result = true
	} else {
		pr.Validation.Reason = "Multiple connections found"
	}
	return pr
}

func (pr *PolicyRequest) PrintYaml(silent bool) *PolicyRequest {
	logrus.Debug(pr)
	if silent {
		return pr
	}
	_yaml, err := yaml.Marshal(pr)
	if err == nil {
		fmt.Print(string(_yaml))
	} else {
		logrus.Errorf("Failed to marshal Policy: %v", err)
	}
	return pr
}

func (pr *PolicyRequest) IfValidated(
	valid func(),
	invalid func(),
) *PolicyRequest {
	if pr.Validation.Result {
		valid()
	} else {
		invalid()
	}
	return pr
}
