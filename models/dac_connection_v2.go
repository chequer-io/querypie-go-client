package models

import (
	"fmt"
	"qpc/utils"
)

type SummarizedConnectionV2 struct {
	Uuid              string `json:"uuid" gorm:"primaryKey"`
	DatabaseType      string `json:"databaseType"`
	CloudProviderType string `json:"cloudProviderType"`
	CloudProviderUuid string `json:"cloudProviderUuid"`

	Name string `json:"name"`

	AdditionalInfo SummarizedAdditionalInfo `json:"additionalInfo" gorm:"-"`

	Zones   []SummarizedZone `json:"zones" gorm:"-"`
	Ledger  bool             `json:"ledger"`
	Deleted bool             `json:"deleted"`

	CreatedAt     string   `json:"createdAt"`
	CreatedBy     Modifier `json:"createdBy" gorm:"foreignKey:CreatedByUuid"`
	CreatedByUuid string   `json:"-"`
	UpdatedAt     string   `json:"updatedAt"`
	UpdatedBy     Modifier `json:"updatedBy" gorm:"foreignKey:UpdatedByUuid"`
	UpdatedByUuid string   `json:"-"`
}

func (c SummarizedConnectionV2) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, DatabaseType=%s, "+
			"CloudProviderType=%s, CloudProviderUuid=%s, "+
			"Name=%s,, AdditionalInfo=%v, "+
			"Zones=%v, Ledger=%t, Deleted=%t, "+
			"CreatedAt=%s, CreatedBy=%v, UpdatedAt=%s, UpdatedBy=%v }",
		c.Uuid, c.DatabaseType,
		c.CloudProviderType, c.CloudProviderUuid,
		c.Name, c.AdditionalInfo,
		c.Zones, c.Ledger, c.Deleted,
		c.CreatedAt, c.CreatedBy, c.UpdatedAt, c.UpdatedBy,
	)
}

func (c SummarizedConnectionV2) Status() string {
	if c.Deleted {
		return "deleted"
	}
	if c.Ledger {
		return "ledger"
	}
	return "-"
}

func (c SummarizedConnectionV2) ShortCreatedAt() string {
	return utils.ShortDatetimeWithTZ(c.CreatedAt)
}

func (c SummarizedConnectionV2) ShortUpdatedAt() string {
	return utils.ShortDatetimeWithTZ(c.UpdatedAt)
}

func (c SummarizedConnectionV2) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Name=%s }",
		c.Uuid, c.Name,
	)
}

type SummarizedAdditionalInfo struct {
	AuditEnabled       bool   `json:"auditEnabled"`
	DmlSnapshotEnabled bool   `json:"dmlSnapshotEnabled"`
	ProxyEnabled       bool   `json:"proxyEnabled"`
	ProxyAuthType      string `json:"proxyAuthType"`
}

func (i SummarizedAdditionalInfo) String() string {
	return fmt.Sprintf(
		"{ AuditEnabled=%t, DmlSnapshotEnabled=%t "+
			"ProxyEnabled=%t, ProxyAuthType=%s }",
		i.AuditEnabled, i.DmlSnapshotEnabled,
		i.ProxyEnabled, i.ProxyAuthType,
	)
}

type PagedConnectionV2List struct {
	PagedList[SummarizedConnectionV2]
}

type ConnectionV2 struct {
	Uuid              string `json:"uuid"`
	DatabaseType      string `json:"databaseType"`
	CloudProviderType string `json:"cloudProviderType"`
	CloudProviderUuid string `json:"cloudProviderUuid"`

	Name string `json:"name"`

	Clusters          []Cluster         `json:"clusters"`
	ConnectionAccount ConnectionAccount `json:"connectionAccount"`
	HideCredential    bool              `json:"hideCredential"`

	// More tabs
	AdditionalInfo        AdditionalInfo        `json:"additionalInfo"`
	JustificationSettings JustificationSettings `json:"justificationSettings"`

	SslSetting SslSetting `json:"sslSetting"`
	SshSetting SshSetting `json:"sshSetting"`

	ConnectionOwners []ConnectionOwner `json:"connectionOwners"`

	AdvancedPrivilegeSetting []AdvancedPrivilegeSetting `json:"advancedPrivilegeSetting"`

	Zones        []Zone       `json:"zones"`
	Ledger       bool         `json:"ledger"`
	VendorDetail VendorDetail `json:"vendorDetail"`

	CreatedAt string   `json:"createdAt"`
	CreatedBy Modifier `json:"createdBy"`
	UpdatedAt string   `json:"updatedAt"`
	UpdatedBy Modifier `json:"updatedBy"`
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
