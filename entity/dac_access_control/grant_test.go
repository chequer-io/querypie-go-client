package dac_access_control

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseFixtureV2DacGrantResponse(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("fixture/v2_grant_response.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a GrantResponse struct
	var grantResponse GrantResponse
	if err := json.Unmarshal(data, &grantResponse); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if grantResponse.Uuid == "" {
		t.Errorf("Expected UUID to be non-empty")
	}
	if grantResponse.Name == "" {
		t.Errorf("Expected name to be non-empty")
	}
	if grantResponse.UserType == "" {
		t.Errorf("Expected user type to be non-empty")
	}
	if grantResponse.MappedConnection.CloudProvider.Name == "" {
		t.Errorf("Expected cloud provider name to be non-empty")
	}
	if grantResponse.MappedConnection.CloudProvider.Uuid == "" {
		t.Errorf("Expected cloud provider UUID to be non-empty")
	}
	if grantResponse.MappedConnection.ClusterUuid == "" {
		t.Errorf("Expected cluster UUID to be non-empty")
	}
	if grantResponse.MappedConnection.DatabaseType == "" {
		t.Errorf("Expected database type to be non-empty")
	}
	if grantResponse.MappedConnection.Name == "" {
		t.Errorf("Expected mapped connection name to be non-empty")
	}
	if grantResponse.MappedConnection.Privilege.Name == "" {
		t.Errorf("Expected privilege name to be non-empty")
	}
	if grantResponse.MappedConnection.Privilege.Uuid == "" {
		t.Errorf("Expected privilege UUID to be non-empty")
	}
	if grantResponse.MappedConnection.Status == "" {
		t.Errorf("Expected status to be non-empty")
	}
}
