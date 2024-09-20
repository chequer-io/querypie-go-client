package dac_access_control

import (
	"encoding/json"
	"fmt"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"os"
	"testing"
)

func TestParseFixtureV2DacAccessControlList(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("fixture/v2_access_control_list.json")
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
	data, err := os.ReadFile("fixture/v2_access_control_detail.json")
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

func TestParseFixtureV2DacAccessControlList2(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("fixture/v2_access_control_list-2.json")
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

func TestParseAndMarshalV2DacAccessControlList3(t *testing.T) {
	if true {
		// Actual response of /api/external/v2/dac/access-controls
		// seems like to be AccessControl instead of SummarizedAccessControl.
		// TODO(JK): Investigate why the response is different from the fixture.
		t.Skip("Skipping TestParseAndMarshalV2DacAccessControlList3")
	}

	// Read the JSON file
	json1, err := os.ReadFile("fixture/v2_access_control_list-3.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a SummarizedAccessControlPagedList struct
	var acl SummarizedAccessControlPagedList
	if err := json.Unmarshal(json1, &acl); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	json2, err := json.Marshal(acl)

	d := gojsondiff.New()
	diff, err := d.Compare(json1, json2)
	if err != nil {
		t.Fatalf("Failed to compare JSON data: %v", err)
	}

	if diff.Modified() {
		fmt.Println("JSONs are different:")
		f := formatter.NewDeltaFormatter()
		diffString, _ := f.Format(diff)
		fmt.Println(diffString)
		t.Errorf("Expected JSON data to be identical")
	}

}
