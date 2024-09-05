package models

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseFixtureV1Users(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../test/fixture_v1_users.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a PagedUserList struct
	var users PagedUserList
	if err := json.Unmarshal(data, &users); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	if len(users.List) == 0 {
		t.Fatalf("Expected non-empty list of users")
	}

	for _, user := range users.List {
		if user.Uuid == "" {
			t.Errorf("Expected user UUID to be non-empty")
		}
		if user.Email == "" {
			t.Errorf("Expected user email to be non-empty")
		}
	}

	if users.Page.TotalElements == 0 {
		t.Errorf("Expected total elements to be non-zero")
	}
}
