package dac_connection

import (
	"encoding/json"
	"fmt"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"os"
	"testing"
)

func TestParseFixtureV2DacConnectionList(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../../test/fixture_v2_dac_connection_list.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a PagedConnectionV2List
	var pagedList PagedConnectionV2List
	if err := json.Unmarshal(data, &pagedList); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if len(pagedList.List) == 0 {
		t.Fatalf("Expected non-empty list of users")
	}

	for _, item := range pagedList.List {
		if item.Uuid == "" {
			t.Errorf("Expected connection UUID to be non-empty")
		}
		if item.Name == "" {
			t.Errorf("Expected connection name to be non-empty")
		}
	}

	if pagedList.Page.TotalElements == 0 {
		t.Errorf("Expected total elements to be non-zero")
	}
}

func TestParseAndMarshalV2DacConnectionList(t *testing.T) {
	// Read the JSON file
	json1, err := os.ReadFile("../../test/fixture_v2_dac_connection_list.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a PagedUserV1List struct
	var pagedList PagedConnectionV2List
	if err := json.Unmarshal(json1, &pagedList); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	json2, err := json.Marshal(pagedList)

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

func TestParseFixtureAndValidateV2DacConnectionDetail(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../../test/fixture_v2_dac_connection_detail.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a PagedConnectionV2List
	var connection ConnectionV2
	if err := json.Unmarshal(data, &connection); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if connection.Uuid == "" {
		t.Errorf("Expected connection UUID to be non-empty")
	}
	if connection.Name == "" {
		t.Errorf("Expected connection name to be non-empty")
	}
	if connection.DatabaseType == "" {
		t.Errorf("Expected database type to be non-empty")
	}
	if connection.CreatedAt == "" {
		t.Errorf("Expected created at to be non-empty")
	}
	if connection.CreatedBy.Uuid == "" {
		t.Errorf("Expected created by UUID to be non-empty")
	}
	if connection.UpdatedAt == "" {
		t.Errorf("Expected updated at to be non-empty")
	}
	if connection.UpdatedBy.Uuid == "" {
		t.Errorf("Expected updated by UUID to be non-empty")
	}
	if len(connection.Zones) == 0 {
		t.Errorf("Expected non-empty list of zones")
	}
	for _, zone := range connection.Zones {
		if zone.Uuid == "" {
			t.Errorf("Expected zone UUID to be non-empty")
		}
		if zone.Name == "" {
			t.Errorf("Expected zone name to be non-empty")
		}
	}
}

func TestParseAndMarshalV2DacConnectionDetail(t *testing.T) {
	// Read the JSON file
	json1, err := os.ReadFile("../../test/fixture_v2_dac_connection_detail.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a Connection struct
	var connection ConnectionV2
	if err := json.Unmarshal(json1, &connection); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	json2, err := json.Marshal(connection)

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
