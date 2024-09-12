package model

type SecurityV2 struct {
	AccountLockoutPolicy         AccountLockoutPolicy         `json:"accountLockoutPolicy"`
	PasswordSetting              PasswordSetting              `json:"passwordSetting"`
	Timeout                      Timeout                      `json:"timeout"`
	WebIpAccessControl           WebIpAccessControl           `json:"webIpAccessControl"`
	DbConnectionSecurity         DbConnectionSecurity         `json:"dbConnectionSecurity"`
	ServerConnectionSecurity     ServerConnectionSecurity     `json:"serverConnectionSecurity"`
	KubernetesConnectionSecurity KubernetesConnectionSecurity `json:"kubernetesConnectionSecurity"`
	VaultEnabled                 bool                         `json:"vaultEnabled"`
	SecretStores                 []SecretStore                `json:"secretStores"`
	Others                       Others                       `json:"others"`
}

type AccountLockoutPolicy struct {
	ExpirationPeriod      int     `json:"expirationPeriod"`
	AccountLockoutEnabled bool    `json:"accountLockoutEnabled"`
	AccountLockoutSetting *string `json:"accountLockoutSetting"`
}

type PasswordSetting struct {
	MaxPasswordAge       int `json:"maxPasswordAge"`
	PasswordHistoryCount int `json:"passwordHistoryCount"`
}

type Timeout struct {
	WebInactivityTimeout int `json:"webInactivityTimeout"`
	AgentSessionTimeout  int `json:"agentSessionTimeout"`
}

type WebIpAccessControl struct {
	IpBands       []string `json:"ipBands"`
	Individualize bool     `json:"individualize"`
}

type DbConnectionSecurity struct {
	PrivilegeDeactivationPeriod int  `json:"privilegeDeactivationPeriod"`
	MaxAccessDuration           int  `json:"maxAccessDuration"`
	DisplayRows                 bool `json:"displayRows"`
	AdvancedPrivilegeSetting    bool `json:"advancedPrivilegeSetting"`
}

type ServerConnectionSecurity struct {
	OsAccountLockoutEnabled   bool    `json:"osAccountLockoutEnabled"`
	OsAccountLockoutSetting   *string `json:"osAccountLockoutSetting"`
	SessionTerminationEnabled bool    `json:"sessionTerminationEnabled"`
	ServerSessionTimeout      int     `json:"serverSessionTimeout"`
	SessionTerminationSetting *string `json:"sessionTerminationSetting"`
}

type KubernetesConnectionSecurity struct {
	IdleSessionTimeoutMinutes int `json:"idleSessionTimeoutMinutes"`
}

type Others struct {
	FileExportEncryption string `json:"fileExportEncryption"`
}
