package models

type SummarizedConnectionV2 struct {
	AdditionalInfo    SummarizedAdditionalInfo `json:"additionalInfo"`
	CloudProviderType string                   `json:"cloudProviderType"`
	CloudProviderUuid string                   `json:"cloudProviderUuid"`
	CreatedAt         string                   `json:"createdAt"`
	CreatedBy         Modifier                 `json:"createdBy"`
	DatabaseType      string                   `json:"databaseType"`
	Deleted           bool                     `json:"deleted"`
	Ledger            bool                     `json:"ledger"`
	Name              string                   `json:"name"`
	UpdatedAt         string                   `json:"updatedAt"`
	UpdatedBy         Modifier                 `json:"updatedBy"`
	Uuid              string                   `json:"uuid"`
	Zones             []SummarizedZone         `json:"zones"`
}

type SummarizedAdditionalInfo struct {
	AuditEnabled       bool   `json:"auditEnabled"`
	DmlSnapshotEnabled bool   `json:"dmlSnapshotEnabled"`
	ProxyEnabled       bool   `json:"proxyEnabled"`
	ProxyAuthType      string `json:"proxyAuthType"`
}

type PagedConnectionV2List struct {
	PagedList[SummarizedConnectionV2]
}

type ConnectionV2 struct {
	AdditionalInfo           AdditionalInfo             `json:"additionalInfo"`
	AdvancedPrivilegeSetting []AdvancedPrivilegeSetting `json:"advancedPrivilegeSetting"`
	CloudProviderType        string                     `json:"cloudProviderType"`
	CloudProviderUuid        string                     `json:"cloudProviderUuid"`
	Clusters                 []Cluster                  `json:"clusters"`
	ConnectionAccount        ConnectionAccount          `json:"connectionAccount"`
	ConnectionOwners         []ConnectionOwner          `json:"connectionOwners"`
	CreatedAt                string                     `json:"createdAt"`
	CreatedBy                Modifier                   `json:"createdBy"`
	DatabaseType             string                     `json:"databaseType"`
	HideCredential           bool                       `json:"hideCredential"`
	JustificationSettings    JustificationSettings      `json:"justificationSettings"`
	Ledger                   bool                       `json:"ledger"`
	Name                     string                     `json:"name"`
	SshSetting               SshSetting                 `json:"sshSetting"`
	SslSetting               SslSetting                 `json:"sslSetting"`
	UpdatedAt                string                     `json:"updatedAt"`
	UpdatedBy                Modifier                   `json:"updatedBy"`
	Uuid                     string                     `json:"uuid"`
	VendorDetail             VendorDetail               `json:"vendorDetail"`
	Zones                    []Zone                     `json:"zones"`
}

type AdditionalInfo struct {
	AccessEndTime       string     `json:"accessEndTime"`
	AccessStartTime     string     `json:"accessStartTime"`
	AuditEnabled        bool       `json:"auditEnabled"`
	DatabaseVersion     string     `json:"databaseVersion"`
	Description         string     `json:"description"`
	DmlSnapshotEnabled  bool       `json:"dmlSnapshotEnabled"`
	LoginRules          LoginRules `json:"loginRules"`
	MaxDisplayRows      int        `json:"maxDisplayRows"`
	MaxExportRows       int        `json:"maxExportRows"`
	NetworkId           string     `json:"networkId"`
	ProxyAuthType       string     `json:"proxyAuthType"`
	ProxyEnabled        bool       `json:"proxyEnabled"`
	WeekdayAccessDenied string     `json:"weekdayAccessDenied"`
}

type LoginRules struct {
	Interval         string `json:"interval"`
	MaxLoginFailures int    `json:"maxLoginFailures"`
}

type AdvancedPrivilegeSetting struct {
	DbAccountName string `json:"dbAccountName"`
	PrivilegeName string `json:"privilegeName"`
	PrivilegeUuid string `json:"privilegeUuid"`
}

type Cluster struct {
	CloudIdentifier string `json:"cloudIdentifier"`
	Deleted         bool   `json:"deleted"`
	Host            string `json:"host"`
	Port            string `json:"port"`
	ReplicationType string `json:"replicationType"`
	Uuid            string `json:"uuid"`
}

type KerberosProtocol struct {
	Principal   string `json:"principal"`
	Realm       string `json:"realm"`
	ServiceName string `json:"serviceName"`
}

type SecretStore struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Uuid    string `json:"uuid"`
}

type ConnectionAccount struct {
	DbAccountName      string           `json:"dbAccountName"`
	KerberosProtocol   KerberosProtocol `json:"kerberosProtocol"`
	SecretStore        SecretStore      `json:"secretStore"`
	SecretStoreEnabled bool             `json:"secretStoreEnabled"`
	Type               string           `json:"type"`
}

type OwnedBy struct {
	Active   bool   `json:"active"`
	AuthType string `json:"authType"`
	Email    string `json:"email"`
	LoginId  string `json:"loginId"`
	Name     string `json:"name"`
	UserType string `json:"userType"`
	Uuid     string `json:"uuid"`
}

type ConnectionOwner struct {
	ObjectUuid string  `json:"objectUuid"`
	OwnedBy    OwnedBy `json:"ownedBy"`
	Role       Role    `json:"role"`
	Uuid       string  `json:"uuid"`
}

type JustificationSettings struct {
	RequireExecuteReason      bool `json:"requireExecuteReason"`
	RequireExportDataReason   bool `json:"requireExportDataReason"`
	RequireExportSchemaReason bool `json:"requireExportSchemaReason"`
	RequireImportDataReason   bool `json:"requireImportDataReason"`
	RequireImportSchemaReason bool `json:"requireImportSchemaReason"`
}

type SshSetting struct {
	SshConfigName string `json:"sshConfigName"`
	SshConfigUuid string `json:"sshConfigUuid"`
	UseSsh        bool   `json:"useSsh"`
}

type SslSetting struct {
	SslConfigName string `json:"sslConfigName"`
	SslConfigUuid string `json:"sslConfigUuid"`
	UseSsl        bool   `json:"useSsl"`
}

type VendorDetail struct {
	CloudRegion   string `json:"cloudRegion"`
	DatabaseName  string `json:"databaseName"`
	WorkGroupName string `json:"workGroupName"`
}
