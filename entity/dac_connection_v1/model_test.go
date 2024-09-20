package dac_connection_v1

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseConnectionsList(t *testing.T) {
	// Read the JSON file
	filePath := "fixture/connections_list_1.json"
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read JSON file")

	// Unmarshal the JSON data into the struct
	var connectionList ConnectionPagedList
	err = json.Unmarshal(fileData, &connectionList)
	assert.NoError(t, err, "Failed to unmarshal JSON data")

	// Validate the parsed data
	assert.NotEmpty(t, connectionList.List, "Expected non-empty list of connections")
	assert.NotNil(t, connectionList.Page, "Expected non-nil page information")

	for _, connection := range connectionList.List {
		assert.NotEmpty(t, connection.Uuid, "Expected connection UUID to be non-empty")
		assert.NotEmpty(t, connection.Name, "Expected connection name to be non-empty")
		assert.NotEmpty(t, connection.DatabaseType, "Expected database type to be non-empty")
		assert.NotEmpty(t, connection.ConnectionAccount.Type, "Expected connection account type to be non-empty")
		assert.NotEmpty(t, connection.Clusters, "Expected non-empty list of clusters")
		assert.NotEmpty(t, connection.ConnectionOwners, "Expected non-empty list of connection owners")
		assert.NotEmpty(t, connection.Tags, "Expected non-empty list of tags")
		assert.NotEmpty(t, connection.Zones, "Expected non-empty list of zones")
	}
}
