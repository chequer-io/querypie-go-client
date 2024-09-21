package dac_policy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"qpc/entity/dac_connection"
)

type PolicyRequest struct {
	ClusterGroupUuid string     `json:"clusterGroupUuid" yaml:"clusterGroupUuid"`
	PolicyType       PolicyType `json:"policyType" yaml:"policyType"`
	Title            string     `json:"title" yaml:"title"`
}

func (pr PolicyRequest) Post() {
	logrus.Warnf("Not Yet Implemented: %v", pr)
}

type UserInput struct {
	Connection []SummarizedConnectionForPolicy `json:"-" yaml:"connection"`
	PolicyType PolicyType                      `json:"-" yaml:"policyType"`
	Name       string                          `json:"-" yaml:"name"`
}

type ValidatablePolicyRequest struct {
	UserInput UserInput `json:"-" yaml:"USER_INPUT"`

	Validation struct {
		Result bool     `json:"-" yaml:"result"`
		Reason []string `json:"-" yaml:"reason,omitempty"`
	} `json:"-" yaml:"VALIDATION"`

	PolicyRequest *PolicyRequest `json:"-" yaml:"POLICY_REQUEST,omitempty"`
}

func GeneratePolicyRequest(connection string, policyType PolicyType, name string) *ValidatablePolicyRequest {
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
	return &ValidatablePolicyRequest{
		UserInput: UserInput{
			Connection: matched,
			PolicyType: policyType,
			Name:       name,
		},
	}
}

func (pr *ValidatablePolicyRequest) Validate() *ValidatablePolicyRequest {
	pr.Validation.Result = true

	if len(pr.UserInput.Connection) == 0 {
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Connection not found")
	} else if len(pr.UserInput.Connection) == 1 {
		// do nothing
	} else {
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Multiple connections found")
	}

	switch pr.UserInput.PolicyType {
	case DataLevel, DataAccess, DataMasking, Notification, Ledger:
		// do nothing
	default:
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Invalid policy type")
	}

	if len(pr.UserInput.Name) > 0 {
		// do nothing
	} else {
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Name is empty")
	}

	if pr.Validation.Result {
		pr.PolicyRequest = &PolicyRequest{
			ClusterGroupUuid: pr.UserInput.Connection[0].Uuid,
			PolicyType:       pr.UserInput.PolicyType,
			Title:            pr.UserInput.Name,
		}
	}
	return pr
}

func (pr *ValidatablePolicyRequest) PrintYaml(silent bool) *ValidatablePolicyRequest {
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

func (pr *ValidatablePolicyRequest) IfValidated(
	valid func(request *PolicyRequest),
	invalid func(),
) *ValidatablePolicyRequest {
	if pr.Validation.Result {
		valid(pr.PolicyRequest)
	} else {
		invalid()
	}
	return pr
}
