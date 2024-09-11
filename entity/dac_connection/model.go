package dac_connection

import (
	"fmt"
	"qpc/models"
	"qpc/utils"
)

type SummarizedConnectionV2 struct {
	Uuid              string `json:"uuid" gorm:"primaryKey"`
	DatabaseType      string `json:"databaseType"`
	CloudProviderType string `json:"cloudProviderType"`
	CloudProviderUuid string `json:"cloudProviderUuid"`

	Name string `json:"name"`

	AdditionalInfo SummarizedAdditionalInfo `json:"additionalInfo" gorm:"-"`

	Zones   []models.SummarizedZone `json:"zones" gorm:"-"`
	Ledger  bool                    `json:"ledger"`
	Deleted bool                    `json:"deleted"`

	CreatedAt     string          `json:"createdAt"`
	CreatedBy     models.Modifier `json:"createdBy" gorm:"foreignKey:CreatedByUuid"`
	CreatedByUuid string          `json:"-"`
	UpdatedAt     string          `json:"updatedAt"`
	UpdatedBy     models.Modifier `json:"updatedBy" gorm:"foreignKey:UpdatedByUuid"`
	UpdatedByUuid string          `json:"-"`
}

func (sc *SummarizedConnectionV2) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, DatabaseType=%s, "+
			"CloudProviderType=%s, CloudProviderUuid=%s, "+
			"Name=%s,, AdditionalInfo=%v, "+
			"Zones=%v, Ledger=%t, Deleted=%t, "+
			"CreatedAt=%s, CreatedBy=%v, UpdatedAt=%s, UpdatedBy=%v }",
		sc.Uuid, sc.DatabaseType,
		sc.CloudProviderType, sc.CloudProviderUuid,
		sc.Name, sc.AdditionalInfo,
		sc.Zones, sc.Ledger, sc.Deleted,
		sc.CreatedAt, sc.CreatedBy, sc.UpdatedAt, sc.UpdatedBy,
	)
}

func (sc *SummarizedConnectionV2) Status() string {
	if sc.Deleted {
		return "deleted"
	}
	if sc.Ledger {
		return "ledger"
	}
	return "-"
}

func (sc *SummarizedConnectionV2) ShortCreatedAt() string {
	return utils.ShortDatetimeWithTZ(sc.CreatedAt)
}

func (sc *SummarizedConnectionV2) ShortUpdatedAt() string {
	return utils.ShortDatetimeWithTZ(sc.UpdatedAt)
}

func (sc *SummarizedConnectionV2) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Name=%s }",
		sc.Uuid, sc.Name,
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
	List []SummarizedConnectionV2 `json:"list"`
	Page models.Page              `json:"page"`
}

func (cl *PagedConnectionV2List) GetPage() models.Page {
	return cl.Page
}

func (cl *PagedConnectionV2List) GetList() []SummarizedConnectionV2 {
	return cl.List
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

	Zones        []models.Zone `json:"zones"`
	Ledger       bool          `json:"ledger"`
	VendorDetail VendorDetail  `json:"vendorDetail"`

	CreatedAt string          `json:"createdAt"`
	CreatedBy models.Modifier `json:"createdBy"`
	UpdatedAt string          `json:"updatedAt"`
	UpdatedBy models.Modifier `json:"updatedBy"`
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

type ConnectionAccount struct {
	DbAccountName      string             `json:"dbAccountName"`
	KerberosProtocol   KerberosProtocol   `json:"kerberosProtocol"`
	SecretStore        models.SecretStore `json:"secretStore"`
	SecretStoreEnabled bool               `json:"secretStoreEnabled"`
	Type               string             `json:"type"`
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
	ObjectUuid string      `json:"objectUuid"`
	OwnedBy    OwnedBy     `json:"ownedBy"`
	Role       models.Role `json:"role"`
	Uuid       string      `json:"uuid"`
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
