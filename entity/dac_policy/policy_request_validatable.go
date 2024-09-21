package dac_policy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"qpc/entity/dac_connection"
	"qpc/model"
)

func GeneratePolicyRequest(connection string, policyType PolicyType, name string) *PolicyRequestValidatable {
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
	return &PolicyRequestValidatable{
		UserInput: UserInput{
			Connection: matched,
			PolicyType: policyType,
			Name:       name,
		},
	}
}

type UserInput struct {
	Connection []SummarizedConnectionForPolicy `json:"-" yaml:"connection"`
	PolicyType PolicyType                      `json:"-" yaml:"policyType"`
	Name       string                          `json:"-" yaml:"name"`
}

type PolicyRequestValidatable struct {
	UserInput     UserInput        `json:"-" yaml:"USER_INPUT"`
	Validation    model.Validation `json:"-" yaml:"VALIDATION"`
	PolicyRequest *PolicyRequest   `json:"-" yaml:"POLICY_REQUEST,omitempty"`
}

func (pr *PolicyRequestValidatable) Validate() *PolicyRequestValidatable {
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

func (pr *PolicyRequestValidatable) ValidateForDelete() *PolicyRequestValidatable {
	pr.Validate()
	if !pr.Validation.Result {
		return pr
	}

	policy := (&Policy{}).FirstByClusterGroupUuidAndName(
		pr.UserInput.Connection[0].Uuid,
		pr.UserInput.Name,
	)

	if policy == nil {
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Policy not found")
	} else {
		pr.PolicyRequest.PolicyUuid = policy.Uuid
	}
	return pr
}

func (pr *PolicyRequestValidatable) PrintYaml(silent bool) *PolicyRequestValidatable {
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

func (pr *PolicyRequestValidatable) UnlessValidated(
	invalid func(),
) *PolicyRequestValidatable {
	if !pr.Validation.Result {
		invalid()
	}
	return pr
}
