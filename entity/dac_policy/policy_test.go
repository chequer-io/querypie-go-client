package dac_policy

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParsePolicyCreateRequest1(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/policy_create_request_1.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var request Policy
	err = json.Unmarshal(fileData, &request)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the fields
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", request.ClusterGroupUuid, "ClusterGroupUuid mismatch")
	assert.Equal(t, DataAccess, request.PolicyType, "PolicyType mismatch")
	assert.Equal(t, "Policy Title", request.Name, "Name mismatch")
}

func TestParsePolicyCreateResponse1(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/policy_create_response_1.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var response Policy
	err = json.Unmarshal(fileData, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the fields
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", response.Uuid, "Uuid mismatch")
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", response.ClusterGroupUuid, "ClusterGroupUuid mismatch")
	assert.Equal(t, "Policy Title", response.Name, "Name mismatch")
	assert.Equal(t, 10, response.NumberOfRules, "NumberOfRules mismatch")
	assert.Equal(t, true, response.Enabled, "Enabled mismatch")
	assert.Equal(t, "2019-08-24T14:15:22Z", response.CreatedAt, "CreatedAt mismatch")
	assert.Equal(t, "Bruce Han", response.CreatedBy.Name, "CreatedUser mismatch")
	assert.Equal(t, "2019-08-24T14:15:22Z", response.UpdatedAt, "UpdatedAt mismatch")
	assert.Equal(t, "Bruce Han", response.UpdatedBy.Name, "UpdatedUser mismatch")
}

func TestParsePolicyUpdateRequest1(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/policy_update_request_1.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var request Policy
	err = json.Unmarshal(fileData, &request)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the fields
	assert.Equal(t, "", request.ClusterGroupUuid, "ClusterGroupUuid mismatch")
	assert.Equal(t, "Policy Title", request.Name, "Name mismatch")
	assert.Equal(t, UnknownPolicyType, request.PolicyType, "PolicyType mismatch")
}

func TestParsePolicyUpdateResponse1(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/policy_update_response_1.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var response Policy
	err = json.Unmarshal(fileData, &response)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the fields
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", response.Uuid, "Uuid mismatch")
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", response.ClusterGroupUuid, "ClusterGroupUuid mismatch")
	assert.Equal(t, "Policy Title", response.Name, "Name mismatch")
	assert.Equal(t, 10, response.NumberOfRules, "NumberOfRules mismatch")
	assert.Equal(t, true, response.Enabled, "Enabled mismatch")
	assert.Equal(t, "2019-08-24T14:15:22Z", response.CreatedAt, "CreatedAt mismatch")
	assert.Equal(t, "Bruce Han", response.CreatedBy.Name, "CreatedUser mismatch")
	assert.Equal(t, "2019-08-24T14:15:22Z", response.UpdatedAt, "UpdatedAt mismatch")
	assert.Equal(t, "Bruce Han", response.UpdatedBy.Name, "UpdatedUser mismatch")
}

func TestParsePolicyList1(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/policy_list_1.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var list PolicyPagedList
	err = json.Unmarshal(fileData, &list)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the fields
	assert.Equal(t, 1, len(list.List), "List length mismatch")
	assert.Equal(t, 0, list.Page.CurrentPage, "Page number mismatch")
	assert.Equal(t, 40, list.Page.PageSize, "Page size mismatch")
	assert.Equal(t, 1, list.Page.TotalElements, "Total elements mismatch")
	assert.Equal(t, 1, list.Page.TotalPages, "Total pages mismatch")

	// Validate the first item
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", list.List[0].Uuid, "Uuid mismatch")
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", list.List[0].ClusterGroupUuid, "ClusterGroupUuid mismatch")
	assert.Equal(t, "Policy Title", list.List[0].Name, "Name mismatch")
	assert.Equal(t, 10, list.List[0].NumberOfRules, "NumberOfRules mismatch")
	assert.Equal(t, true, list.List[0].Enabled, "Enabled mismatch")
	assert.Equal(t, "2019-08-24T14:15:22Z", list.List[0].CreatedAt, "CreatedAt mismatch")
	assert.Equal(t, "Bruce Han", list.List[0].CreatedBy.Name, "CreatedUser mismatch")
	assert.Equal(t, "2019-08-24T14:15:22Z", list.List[0].UpdatedAt, "UpdatedAt mismatch")
	assert.Equal(t, "Bruce Han", list.List[0].UpdatedBy.Name, "UpdatedUser mismatch")
}

func TestParsePolicyList2(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/policy_list_2.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var list PolicyPagedList
	err = json.Unmarshal(fileData, &list)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the fields
	assert.Equal(t, 7, len(list.List), "List length mismatch")
	assert.Equal(t, 0, list.Page.CurrentPage, "Page number mismatch")
	assert.Equal(t, 40, list.Page.PageSize, "Page size mismatch")
	assert.Equal(t, 7, list.Page.TotalElements, "Total elements mismatch")
	assert.Equal(t, 1, list.Page.TotalPages, "Total pages mismatch")

	// Validate the first item
	assert.Equal(t, "c15f0d96-829c-4271-b8b1-4dd9dc3ffd82", list.List[0].Uuid, "Uuid mismatch")
	assert.Equal(t, "59db4f83-d3c8-4836-9de7-d691b105c9ba", list.List[0].ClusterGroupUuid, "ClusterGroupUuid mismatch")
	assert.Equal(t, "MySQL Sensitive Data Policy", list.List[0].Name, "Name mismatch")
	assert.Equal(t, 2, list.List[0].NumberOfRules, "NumberOfRules mismatch")
	assert.Equal(t, true, list.List[0].Enabled, "Enabled mismatch")
	assert.Equal(t, "2024-09-05T08:02:22.923Z", list.List[0].CreatedAt, "CreatedAt mismatch")
	assert.Equal(t, "Daniel Son", list.List[0].CreatedBy.Name, "CreatedUser mismatch")
	assert.Equal(t, "2024-09-05T08:03:36.708Z", list.List[0].UpdatedAt, "UpdatedAt mismatch")
	assert.Equal(t, "Daniel Son", list.List[0].UpdatedBy.Name, "UpdatedUser mismatch")
}
