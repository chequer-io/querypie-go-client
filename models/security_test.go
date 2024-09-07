package models

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseFixtureV2Security(t *testing.T) {
	// Read the JSON file
	data, err := os.ReadFile("../test/fixture_v2_security.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a SecurityV2 struct
	var security SecurityV2
	if err := json.Unmarshal(data, &security); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Validate the parsed data
	validateSecurityV2(t, security)
}

func validateSecurityV2(t *testing.T, security SecurityV2) {
	if security.AccountLockoutPolicy.AccountLockoutEnabled == true {
		t.Errorf("AccountLockoutPolicy.AccountLockoutEnabled is true")
	}
	if security.PasswordSetting.MaxPasswordAge == 0 {
		t.Errorf("PasswordSetting.MaxPasswordAge is zero")
	}
	if security.Timeout.WebInactivityTimeout == 0 {
		t.Errorf("Timeout.WebInactivityTimeout is zero")
	}
	if len(security.WebIpAccessControl.IpBands) == 0 {
		t.Errorf("WebIpAccessControl.IpBands is empty")
	}
	if security.DbConnectionSecurity.AdvancedPrivilegeSetting == false {
		t.Errorf("DbConnectionSecurity.AdvancedPrivilegeSetting is false")
	}
	if security.ServerConnectionSecurity.OsAccountLockoutEnabled == true {
		t.Errorf("ServerConnectionSecurity.OsAccountLockoutEnabled is true")
	}
	if security.KubernetesConnectionSecurity.IdleSessionTimeoutMinutes == 0 {
		t.Errorf("KubernetesConnectionSecurity.IdleSessionTimeoutMinutes is zero")
	}
	if security.VaultEnabled == false {
		t.Errorf("VaultEnabled is false")
	}
	if len(security.SecretStores) == 0 {
		t.Errorf("SecretStores is empty")
	}
	if security.Others.FileExportEncryption == "" {
		t.Errorf("Others.FileExportEncryption is empty")
	}
}
