package dac_access_control

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseFixtureV2DacAccessControlList(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../../test/fixture_v2_dac_access_control_list.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a SummarizedAccessControlPagedList struct
	var acl SummarizedAccessControlPagedList
	if err := json.Unmarshal(data, &acl); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if len(acl.List) == 0 {
		t.Fatalf("Expected non-empty list of access controls")
	}

	for _, accessControl := range acl.List {
		if accessControl.Uuid == "" {
			t.Errorf("Expected access control UUID to be non-empty")
		}
		if accessControl.Name == "" {
			t.Errorf("Expected access control name to be non-empty")
		}
	}

	if acl.Page.TotalElements == 0 {
		t.Errorf("Expected total elements to be non-zero")
	}
}

func TestParseFixtureAndValidateV2DacAccessControlDetail(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../../test/fixture_v2_dac_access_control_detail.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a AccessControl
	var ac AccessControl
	if err := json.Unmarshal(data, &ac); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if ac.Uuid == "" {
		t.Errorf("Expected UUID to be non-empty")
	}
	if ac.UserType == "" {
		t.Errorf("Expected user type to be non-empty")
	}
	if ac.AuthType == "" {
		t.Errorf("Expected auth type to be non-empty")
	}
	if ac.Name == "" {
		t.Errorf("Expected connection name to be non-empty")
	}
	if len(ac.Members) == 0 {
		t.Errorf("Expected non-empty list of members")
	}
	if len(ac.MappedConnections) == 0 {
		t.Errorf("Expected non-empty list of mapped connections")
	}
	if ac.AdminRole == "" {
		t.Errorf("Expected admin role to be non-empty")
	}
	if ac.LinkedCount == 0 {
		t.Errorf("Expected linked count be non-zero")
	}
}
