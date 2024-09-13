package dac_privilege

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseFixtureV2DacPrivilegeList(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../../test/fixture_v2_dac_privilege_list.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a PrivilegePagedList struct
	var privilegeList PrivilegePagedList
	if err := json.Unmarshal(data, &privilegeList); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if len(privilegeList.List) == 0 {
		t.Fatalf("Expected non-empty list of privileges")
	}

	for _, privilege := range privilegeList.List {
		if privilege.Uuid == "" {
			t.Errorf("Expected privilege UUID to be non-empty")
		}
		if privilege.Name == "" {
			t.Errorf("Expected privilege name to be non-empty")
		}
	}

	if privilegeList.Page.TotalElements == 0 {
		t.Errorf("Expected total elements to be non-zero")
	}
}

func TestParseFixtureAndValidateV2DacPrivilegeDetail(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../../test/fixture_v2_dac_privilege_detail.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a Privilege struct
	var privilege Privilege
	if err := json.Unmarshal(data, &privilege); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if privilege.Uuid == "" {
		t.Errorf("Expected UUID to be non-empty")
	}
	if privilege.Name == "" {
		t.Errorf("Expected privilege name to be non-empty")
	}
	if len(privilege.PrivilegeTypes.Values) == 0 {
		t.Errorf("Expected non-empty list of privilege types")
	}
	if privilege.PrivilegeVendor == "" {
		t.Errorf("Expected privilege vendor to be non-empty")
	}
	if privilege.Status == "" {
		t.Errorf("Expected status to be non-empty")
	}
	if privilege.CreatedAt == "" {
		t.Errorf("Expected created at to be non-empty")
	}
	if privilege.UpdatedAt == "" {
		t.Errorf("Expected updated at to be non-empty")
	}
}
