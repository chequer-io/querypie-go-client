package dac_connection

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"qpc/model"
	"qpc/utils"
)

type SummarizedConnectionV2 struct {
	Uuid              string `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	DatabaseType      string `json:"databaseType" yaml:"databaseType"`
	CloudProviderType string `json:"cloudProviderType" yaml:"cloudProviderType"`
	CloudProviderUuid string `json:"cloudProviderUuid" yaml:"cloudProviderUuid"`

	Name string `json:"name" yaml:"name"`

	AdditionalInfo SummarizedAdditionalInfo `json:"additionalInfo" gorm:"-" yaml:"additionalInfo"`

	Zones   []model.SummarizedZone `json:"zones" gorm:"-" yaml:"zones"`
	Ledger  bool                   `json:"ledger" yaml:"ledger"`
	Deleted bool                   `json:"deleted" yaml:"deleted"`

	CreatedAt     string         `json:"createdAt" yaml:"createdAt"`
	CreatedBy     model.Modifier `json:"createdBy" gorm:"foreignKey:CreatedByUuid" yaml:"createdBy"`
	CreatedByUuid string         `json:"-" yaml:"-"`
	UpdatedAt     string         `json:"updatedAt" gorm:"autoUpdateTime:false" yaml:"updatedAt"`
	UpdatedBy     model.Modifier `json:"updatedBy" gorm:"foreignKey:UpdatedByUuid" yaml:"updatedBy"`
	UpdatedByUuid string         `json:"-" yaml:"-"`
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
	AuditEnabled       bool   `json:"auditEnabled" yaml:"auditEnabled"`
	DmlSnapshotEnabled bool   `json:"dmlSnapshotEnabled" yaml:"dmlSnapshotEnabled"`
	ProxyEnabled       bool   `json:"proxyEnabled" yaml:"proxyEnabled"`
	ProxyAuthType      string `json:"proxyAuthType" yaml:"proxyAuthType"`
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
	Page model.Page               `json:"page"`
}

func (cl *PagedConnectionV2List) GetPage() model.Page {
	return cl.Page
}

func (cl *PagedConnectionV2List) GetList() []SummarizedConnectionV2 {
	return cl.List
}

type ConnectionV2 struct {
	Uuid              string  `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	DatabaseType      string  `json:"databaseType" yaml:"databaseType"`
	CloudProviderType *string `json:"cloudProviderType" yaml:"cloudProviderType"`
	CloudProviderUuid *string `json:"cloudProviderUuid" yaml:"cloudProviderUuid"`

	Name string `json:"name"`

	Clusters          []Cluster         `json:"clusters" gorm:"foreignKey:ConnectionUuid" yaml:"clusters"`
	ConnectionAccount ConnectionAccount `json:"connectionAccount" gorm:"embedded" yaml:"connectionAccount"`
	HideCredential    bool              `json:"hideCredential" yaml:"hideCredential"`

	// More tabs
	AdditionalInfo        AdditionalInfo        `json:"additionalInfo" gorm:"embedded" yaml:"additionalInfo"`
	JustificationSettings JustificationSettings `json:"justificationSettings" gorm:"embedded" yaml:"justificationSettings"`

	SslSetting SslSetting `json:"sslSetting" gorm:"embedded" yaml:"sslSetting"`
	SshSetting SshSetting `json:"sshSetting" gorm:"embedded" yaml:"sshSetting"`

	ConnectionOwners []ConnectionOwner `json:"connectionOwners" gorm:"many2many" yaml:"connectionOwners"`

	AdvancedPrivilegeSetting []AdvancedPrivilegeSetting `json:"advancedPrivilegeSetting" gorm:"json" yaml:"advancedPrivilegeSetting"`

	Zones        []model.Zone `json:"zones" gorm:"json" yaml:"zones"`
	Ledger       bool         `json:"ledger" yaml:"ledger"`
	VendorDetail VendorDetail `json:"vendorDetail" gorm:"json" yaml:"vendorDetail"`

	CreatedAt     string         `json:"createdAt" yaml:"createdAt"`
	CreatedBy     model.Modifier `json:"createdBy" gorm:"foreignKey:CreatedByUuid" yaml:"createdBy"`
	CreatedByUuid string         `json:"-" yaml:"-"`
	UpdatedAt     string         `json:"updatedAt" gorm:"autoUpdateTime:false" yaml:"updatedAt"`
	UpdatedBy     model.Modifier `json:"updatedBy" gorm:"foreignKey:UpdatedByUuid" yaml:"updatedBy"`
	UpdatedByUuid string         `json:"-" yaml:"-"`

	// Internal: HTTP response
	HttpResponse *resty.Response `json:"-" gorm:"-" yaml:"-"`
}

func (c *ConnectionV2) SetHttpResponse(response *resty.Response) {
	c.HttpResponse = response
}

func (c *ConnectionV2) GetHttpResponse() *resty.Response {
	return c.HttpResponse
}

// Ensure ConnectionV2 implements RestResponse
var _ model.RestResponse = (*ConnectionV2)(nil)

func (c *ConnectionV2) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, DatabaseType=%s, "+
			"CloudProviderType=%s, CloudProviderUuid=%s, "+
			"Name=%s,, AdditionalInfo=%v, "+
			"Zones=%v, Ledger=%t, "+
			"CreatedAt=%s, CreatedBy=%v, UpdatedAt=%s, UpdatedBy=%v }",
		c.Uuid, c.DatabaseType,
		utils.OptionalPtr(c.CloudProviderType), utils.OptionalPtr(c.CloudProviderUuid),
		c.Name, c.AdditionalInfo,
		c.Zones, c.Ledger,
		c.CreatedAt, c.CreatedBy, c.UpdatedAt, c.UpdatedBy,
	)
}

func (c *ConnectionV2) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Name=%s }",
		c.Uuid, c.Name,
	)
}

type AdditionalInfo struct {
	AccessEndTime       *string          `json:"accessEndTime" yaml:"accessEndTime"`
	AccessStartTime     *string          `json:"accessStartTime" yaml:"accessStartTime"`
	AuditEnabled        bool             `json:"auditEnabled" yaml:"auditEnabled"`
	DatabaseVersion     *string          `json:"databaseVersion" yaml:"databaseVersion"`
	Description         string           `json:"description" yaml:"description"`
	DmlSnapshotEnabled  bool             `json:"dmlSnapshotEnabled" yaml:"dmlSnapshotEnabled"`
	LoginRules          LoginRules       `json:"loginRules,omitempty" gorm:"embedded" yaml:"loginRules,omitempty"`
	MaxDisplayRows      int              `json:"maxDisplayRows" yaml:"maxDisplayRows"`
	MaxExportRows       int              `json:"maxExportRows" yaml:"maxExportRows"`
	NetworkId           *string          `json:"networkId" yaml:"networkId"`
	ProxyAuthType       *string          `json:"proxyAuthType" yaml:"proxyAuthType"`
	ProxyEnabled        bool             `json:"proxyEnabled" yaml:"proxyEnabled"`
	WeekdayAccessDenied model.StringList `json:"weekdayAccessDenied,omitempty" yaml:"weekdayAccessDenied,omitempty"`
}

type LoginRules struct {
	Interval         string `json:"interval,omitempty" yaml:"interval,omitempty"`
	MaxLoginFailures int    `json:"maxLoginFailures" yaml:"maxLoginFailures"`
}

type AdvancedPrivilegeSetting struct {
	PrivilegeUuid string `json:"privilegeUuid,omitempty" yaml:"privilegeUuid,omitempty"`
	PrivilegeName string `json:"privilegeName,omitempty" yaml:"privilegeName,omitempty"`
	DbAccountName string `json:"dbAccountName,omitempty" yaml:"dbAccountName,omitempty"`

	ConnectionUuid string `json:"-"`
}

type Cluster struct {
	Uuid            string        `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	CloudIdentifier *string       `json:"cloudIdentifier" yaml:"cloudIdentifier"`
	Host            string        `json:"host" yaml:"host"`
	Port            string        `json:"port" yaml:"port"`
	ReplicationType string        `json:"replicationType" yaml:"replicationType"`
	Deleted         bool          `json:"deleted" yaml:"deleted"`
	ConnectionUuid  string        `json:"-" yaml:"-"`
	Connection      *ConnectionV2 `json:"connection,omitempty" gorm:"foreignKey:ConnectionUuid" yaml:"connection,omitempty"`
}

func (c *Cluster) Status() string {
	if c.Deleted {
		return "deleted"
	}
	return "-"
}

func (c *Cluster) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Host=%s, Port=%s, ReplicationType=%s, "+
			"CloudIdentifier=%s, Deleted=%t, Connection=%s }",
		c.Uuid, c.Host, c.Port, c.ReplicationType,
		utils.OptionalPtr(c.CloudIdentifier), c.Deleted,
		c.Connection.ShortID(),
	)
}

type KerberosProtocol struct {
	Principal   string `json:"principal,omitempty" yaml:"principal,omitempty"`
	Realm       string `json:"realm,omitempty" yaml:"realm,omitempty"`
	ServiceName string `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`
}

type ConnectionAccount struct {
	DbAccountName      *string            `json:"dbAccountName" yaml:"dbAccountName"`
	KerberosProtocol   *KerberosProtocol  `json:"kerberosProtocol" gorm:"json" yaml:"kerberosProtocol"`
	SecretStore        *model.SecretStore `json:"secretStore" gorm:"json" yaml:"secretStore"`
	SecretStoreEnabled bool               `json:"secretStoreEnabled" yaml:"secretStoreEnabled"`
	Type               string             `json:"type" yaml:"type"`
}

type OwnedBy struct {
	Uuid     string `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	UserType string `json:"userType" yaml:"userType"`
	AuthType string `json:"authType" yaml:"authType"`
	LoginId  string `json:"loginId" yaml:"loginId"`
	Name     string `json:"name" yaml:"name"`
	Email    string `json:"email" yaml:"email"`
	Active   bool   `json:"active" yaml:"active"`
}

type ConnectionOwner struct {
	ObjectUuid string        `json:"objectUuid" gorm:"primaryKey" yaml:"objectUuid"`
	Connection *ConnectionV2 `json:"connection,omitempty" gorm:"foreignKey:ObjectUuid" yaml:"connection,omitempty"`

	RoleUuid string     `json:"-" gorm:"primaryKey" yaml:"-"`
	Role     model.Role `json:"role" gorm:"foreignKey:RoleUuid" yaml:"role"`

	Uuid    string  `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	OwnedBy OwnedBy `json:"ownedBy" gorm:"foreignKey:Uuid" yaml:"ownedBy"`
}

type JustificationSettings struct {
	RequireExecuteReason      bool `json:"requireExecuteReason" yaml:"requireExecuteReason"`
	RequireExportDataReason   bool `json:"requireExportDataReason" yaml:"requireExportDataReason"`
	RequireExportSchemaReason bool `json:"requireExportSchemaReason" yaml:"requireExportSchemaReason"`
	RequireImportDataReason   bool `json:"requireImportDataReason" yaml:"requireImportDataReason"`
	RequireImportSchemaReason bool `json:"requireImportSchemaReason" yaml:"requireImportSchemaReason"`
}

type SshSetting struct {
	UseSsh        bool    `json:"useSsh" yaml:"useSsh"`
	SshConfigUuid *string `json:"sshConfigUuid" yaml:"sshConfigUuid"`
	SshConfigName *string `json:"sshConfigName" yaml:"sshConfigName"`
}

type SslSetting struct {
	UseSsl        bool    `json:"useSsl" yaml:"useSsl"`
	SslConfigName *string `json:"sslConfigName" yaml:"sslConfigName"`
	SslConfigUuid *string `json:"sslConfigUuid" yaml:"sslConfigUuid"`
}

type VendorDetail struct {
	DatabaseName  string  `json:"databaseName" yaml:"databaseName"`
	Charset       *string `json:"charset" yaml:"charset"`
	Collation     *string `json:"collation" yaml:"collation"`
	CloudRegion   *string `json:"cloudRegion,omitempty" yaml:"cloudRegion,omitempty"`
	WorkGroupName *string `json:"workGroupName,omitempty" yaml:"workGroupName,omitempty"`
}
