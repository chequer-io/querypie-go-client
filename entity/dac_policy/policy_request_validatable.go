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

func (pr *PolicyRequestValidatable) validateConnection() {
	if len(pr.UserInput.Connection) == 0 {
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Connection not found")
	} else if len(pr.UserInput.Connection) == 1 {
		// do nothing
	} else {
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Multiple connections found")
	}
}

func (pr *PolicyRequestValidatable) validatePolicyType() {
	switch pr.UserInput.PolicyType {
	case DataLevel, DataAccess, DataMasking, Notification, Ledger:
		// do nothing
	default:
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Invalid policy type")
	}
}

func (pr *PolicyRequestValidatable) validateName() {
	if len(pr.UserInput.Name) > 0 {
		// do nothing
	} else {
		pr.Validation.Result = false
		pr.Validation.Reason = append(pr.Validation.Reason, "Name is empty")
	}
}

func (pr *PolicyRequestValidatable) createPolicyRequestIfValidated() {
	if pr.Validation.Result {
		pr.PolicyRequest = &PolicyRequest{
			ClusterGroupUuid: pr.UserInput.Connection[0].Uuid,
			PolicyType:       pr.UserInput.PolicyType,
			Title:            pr.UserInput.Name,
		}
	}
}

func (pr *PolicyRequestValidatable) tryToFindPolicyLocallyIfValidated() {
	if pr.Validation.Result {
		policy := (&Policy{}).FirstByClusterGroupUuidAndPolicyType(
			pr.UserInput.Connection[0].Uuid,
			pr.UserInput.PolicyType,
		)
		if policy != nil {
			pr.PolicyRequest.PolicyUuid = policy.Uuid
		}
	}
}

func (pr *PolicyRequestValidatable) Validate() *PolicyRequestValidatable {
	pr.Validation.Result = true

	pr.validateConnection()
	pr.validatePolicyType()
	pr.validateName()
	pr.createPolicyRequestIfValidated()
	pr.tryToFindPolicyLocallyIfValidated()
	return pr
}

func (pr *PolicyRequestValidatable) ValidateForDelete() *PolicyRequestValidatable {
	pr.Validation.Result = true

	pr.validateConnection()
	pr.validatePolicyType()
	// No need to check name for delete operation.
	pr.createPolicyRequestIfValidated()
	pr.tryToFindPolicyLocallyIfValidated()

	if pr.Validation.Result {
		if pr.PolicyRequest.PolicyUuid == "" {
			pr.Validation.Result = false
			pr.Validation.Reason = append(pr.Validation.Reason, "Policy not found")
		}
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
