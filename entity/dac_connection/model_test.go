package dac_connection

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"os"
	"qpc/model"
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

func TestParseAndMarshalV2DacConnectionDetail2(t *testing.T) {
	// Read the JSON file
	json1, err := os.ReadFile("../../test/fixture_v2_dac_connection_detail-2.json")
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

func TestParseConnectionV2FromFixture2(t *testing.T) {
	// Read fixture file
	fixtureData, err := os.ReadFile("../../test/fixture_v2_dac_connection_detail-2.json")
	assert.NoError(t, err, "Failed to read fixture file")

	// Parse JSON
	var connection ConnectionV2
	err = json.Unmarshal(fixtureData, &connection)
	assert.NoError(t, err, "Failed to parse JSON")

	// Validate field values
	assert.Equal(t, "59db4f83-d3c8-4836-9de7-d691b105c9ba", connection.Uuid, "UUID mismatch")
	assert.Equal(t, "mysql.querypie.io", connection.Name, "Name mismatch")
	assert.Equal(t, "MYSQL", connection.DatabaseType, "DatabaseType mismatch")
	assert.Nil(t, connection.CloudProviderUuid, "CloudProviderUuid should be nil")
	assert.Empty(t, connection.CloudProviderType, "CloudProviderType should be empty")

	// Validate time fields
	assert.Equal(t, "2024-06-18T06:24:01.609Z", connection.CreatedAt, "CreatedAt mismatch")
	assert.Equal(t, "2024-09-05T08:22:16.936Z", connection.UpdatedAt, "UpdatedAt mismatch")

	// Validate CreatedBy
	assert.Equal(t, "ella@chequer.io", connection.CreatedBy.Uuid, "CreatedBy UUID mismatch")
	assert.Equal(t, "97bb817d-1fa0-4d8b-9df3-492605655d7f", connection.CreatedBy.Name, "CreatedBy Name mismatch")
	assert.Equal(t, "Ella Lee", connection.CreatedBy.LoginId, "CreatedBy LoginId mismatch")
	assert.Equal(t, "user.masked@chequer.io", connection.CreatedBy.Email, "CreatedBy Email mismatch")

	// Validate UpdatedBy
	assert.Equal(t, "user.masked@chequer.io", connection.UpdatedBy.Uuid, "UpdatedBy UUID mismatch")
	assert.Equal(t, "4afff027-35d0-4fca-8e67-85f757d9c44c", connection.UpdatedBy.Name, "UpdatedBy Name mismatch")
	assert.Equal(t, "Daniel Son", connection.UpdatedBy.LoginId, "UpdatedBy LoginId mismatch")
	assert.Equal(t, "user.masked@chequer.io", connection.UpdatedBy.Email, "UpdatedBy Email mismatch")

	// Validate Zones
	assert.Len(t, connection.Zones, 1, "Zones length mismatch")
	assert.Equal(t, "b21abaaa-bbf2-11ed-9e22-0242ac110002", connection.Zones[0].Uuid, "Zone UUID mismatch")
	assert.Equal(t, "All Allowed (Any)", connection.Zones[0].Name, "Zone Name mismatch")
	assert.Equal(t, model.StringList{Values: []string{"0.0.0.0/0"}}, connection.Zones[0].IpBands, "Zone IpBands mismatch")
	assert.Equal(t, "2024-04-30T11:34:30.000Z", connection.Zones[0].CreatedAt, "Zone CreatedAt mismatch")
	assert.Equal(t, "2024-05-28T10:03:33.806Z", connection.Zones[0].UpdatedAt, "Zone UpdatedAt mismatch")

	// Validate Clusters
	assert.Len(t, connection.Clusters, 1, "Clusters length mismatch")
	assert.Equal(t, "4156920f-b536-4fc6-b1f7-bcf93aec4bd9", connection.Clusters[0].Uuid, "Cluster UUID mismatch")
	assert.Empty(t, connection.Clusters[0].CloudIdentifier, "Cluster CloudIdentifier should be empty")
	assert.Equal(t, "SINGLE", connection.Clusters[0].ReplicationType, "Cluster ReplicationType mismatch")
	assert.Equal(t, "mysql.querypie.io", connection.Clusters[0].Host, "Cluster Host mismatch")
	assert.Equal(t, "3306", connection.Clusters[0].Port, "Cluster Port mismatch")
	assert.False(t, connection.Clusters[0].Deleted, "Cluster Deleted mismatch")

	// Validate ConnectionAccount
	assert.Equal(t, "UIDPWD", connection.ConnectionAccount.Type, "ConnectionAccount Type mismatch")
	assert.False(t, connection.ConnectionAccount.SecretStoreEnabled, "ConnectionAccount SecretStoreEnabled mismatch")
	assert.Empty(t, connection.ConnectionAccount.SecretStore, "ConnectionAccount SecretStore should be empty")
	assert.Empty(t, connection.ConnectionAccount.KerberosProtocol, "ConnectionAccount KerberosProtocol should be empty")
	assert.Equal(t, "querypie", *connection.ConnectionAccount.DbAccountName, "ConnectionAccount DbAccountName mismatch")

	assert.Empty(t, connection.ConnectionOwners, "ConnectionOwners should be empty")
	assert.False(t, connection.HideCredential, "HideCredential mismatch")
	assert.True(t, connection.Ledger, "Ledger mismatch")

	// Validate AdditionalInfo
	assert.Equal(t, 1000, connection.AdditionalInfo.MaxDisplayRows, "AdditionalInfo MaxDisplayRows mismatch")
	assert.Equal(t, 1000, connection.AdditionalInfo.MaxExportRows, "AdditionalInfo MaxExportRows mismatch")
	assert.Empty(t, connection.AdditionalInfo.AccessStartTime, "AdditionalInfo AccessStartTime should be empty")
	assert.Empty(t, connection.AdditionalInfo.AccessEndTime, "AdditionalInfo AccessEndTime should be empty")
	assert.Equal(t, model.StringList{Values: []string{}}, connection.AdditionalInfo.WeekdayAccessDenied, "AdditionalInfo WeekdayAccessDenied mismatch")
	assert.Equal(t, 5, connection.AdditionalInfo.LoginRules.MaxLoginFailures, "AdditionalInfo LoginRules MaxLoginFailures mismatch")
	assert.Equal(t, "10MIN", connection.AdditionalInfo.LoginRules.Interval, "AdditionalInfo LoginRules Interval mismatch")
	assert.Empty(t, connection.AdditionalInfo.DatabaseVersion, "AdditionalInfo DatabaseVersion should be empty")
	assert.True(t, connection.AdditionalInfo.AuditEnabled, "AdditionalInfo AuditEnabled mismatch")
	assert.True(t, connection.AdditionalInfo.DmlSnapshotEnabled, "AdditionalInfo DmlSnapshotEnabled mismatch")
	assert.False(t, connection.AdditionalInfo.ProxyEnabled, "AdditionalInfo ProxyEnabled mismatch")
	assert.Empty(t, connection.AdditionalInfo.ProxyAuthType, "AdditionalInfo ProxyAuthType should be empty")
	assert.Empty(t, connection.AdditionalInfo.NetworkId, "AdditionalInfo NetworkId should be empty")
	assert.Equal(t, "", connection.AdditionalInfo.Description, "AdditionalInfo Description mismatch")

	// Validate JustificationSettings
	assert.True(t, connection.JustificationSettings.RequireExportSchemaReason, "JustificationSettings RequireExportSchemaReason mismatch")
	assert.True(t, connection.JustificationSettings.RequireExportDataReason, "JustificationSettings RequireExportDataReason mismatch")
	assert.False(t, connection.JustificationSettings.RequireExecuteReason, "JustificationSettings RequireExecuteReason mismatch")
	assert.False(t, connection.JustificationSettings.RequireImportSchemaReason, "JustificationSettings RequireImportSchemaReason mismatch")
	assert.False(t, connection.JustificationSettings.RequireImportDataReason, "JustificationSettings RequireImportDataReason mismatch")

	// Validate SSL and SSH settings
	assert.False(t, connection.SslSetting.UseSsl, "SslSetting UseSsl mismatch")
	assert.Empty(t, connection.SslSetting.SslConfigName, "SslSetting SslConfigName should be empty")
	assert.Empty(t, connection.SslSetting.SslConfigUuid, "SslSetting SslConfigUuid should be empty")
	assert.False(t, connection.SshSetting.UseSsh, "SshSetting UseSsh mismatch")
	assert.Empty(t, connection.SshSetting.SshConfigName, "SshSetting SshConfigName should be empty")
	assert.Empty(t, connection.SshSetting.SshConfigUuid, "SshSetting SshConfigUuid should be empty")

	assert.Empty(t, connection.AdvancedPrivilegeSetting, "AdvancedPrivilegeSetting should be empty")

	// Validate VendorDetail
	assert.Equal(t, "", connection.VendorDetail.DatabaseName, "VendorDetail DatabaseName mismatch")
	assert.Nil(t, connection.VendorDetail.Charset, "VendorDetail Charset should be nil")
	assert.Nil(t, connection.VendorDetail.Collation, "VendorDetail Collation should be nil")
}
