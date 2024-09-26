package dac_policy

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseSensitiveDataRule_1(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/sensitive-data-rule_create_response_1.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var rule SensitiveDataRule
	err = json.Unmarshal(fileData, &rule)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the parsed data
	assert.Equal(t, "2019-08-24T14:15:22Z", rule.CreatedAt, "Expected createdAt to match")
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", rule.PolicyUuid, "Expected policyUuid to match")
	assert.Equal(t, "COLUMN", rule.ObjectType, "Expected objectType to match")
	assert.Equal(t, 1, rule.Level, "Expected level to match")
	assert.Equal(t, []string{"table", "column"}, rule.ObjectPath, "Expected objectPath to match")
	assert.Equal(t, "35cfa847-165c-48b5-bb4e-af271e490f19", rule.Uuid, "Expected uuid to match")
	assert.Equal(t, "2019-08-24T14:15:22Z", rule.UpdatedAt, "Expected updatedAt to match")

}

func TestParseAndMarshalSensitiveDataRuleList_1(t *testing.T) {
	// Read the JSON file
	json1, err := os.ReadFile("fixture/sensitive-data-rule_list_1.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a SensitiveDataRuleList
	var list []SensitiveDataRule
	if err := json.Unmarshal(json1, &list); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Marshal the struct back to JSON
	json2, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	// Compare the original and marshaled JSON data
	var originalData []map[string]interface{}
	if err := json.Unmarshal(json1, &originalData); err != nil {
		t.Fatalf("Failed to unmarshal original JSON data: %v", err)
	}

	var marshaledData []map[string]interface{}
	if err := json.Unmarshal(json2, &marshaledData); err != nil {
		t.Fatalf("Failed to unmarshal marshaled JSON data: %v", err)
	}

	assert.Equal(t, originalData, marshaledData, "JSON data should be identical")
}
