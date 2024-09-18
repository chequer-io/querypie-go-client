package dac_connection

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"qpc/model"
	"qpc/utils"
	"time"
)

type SummarizedConnectionV2 struct {
	Uuid              string `json:"uuid" gorm:"primaryKey"`
	DatabaseType      string `json:"databaseType"`
	CloudProviderType string `json:"cloudProviderType"`
	CloudProviderUuid string `json:"cloudProviderUuid"`

	Name string `json:"name"`

	AdditionalInfo SummarizedAdditionalInfo `json:"additionalInfo" gorm:"-"`

	Zones   []model.SummarizedZone `json:"zones" gorm:"-"`
	Ledger  bool                   `json:"ledger"`
	Deleted bool                   `json:"deleted"`

	CreatedAt     string         `json:"createdAt"`
	CreatedBy     model.Modifier `json:"createdBy" gorm:"foreignKey:CreatedByUuid"`
	CreatedByUuid string         `json:"-"`
	UpdatedAt     string         `json:"updatedAt" gorm:"autoUpdateTime:false"`
	UpdatedBy     model.Modifier `json:"updatedBy" gorm:"foreignKey:UpdatedByUuid"`
	UpdatedByUuid string         `json:"-"`
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
	Page model.Page               `json:"page"`
}

func (cl *PagedConnectionV2List) GetPage() model.Page {
	return cl.Page
}

func (cl *PagedConnectionV2List) GetList() []SummarizedConnectionV2 {
	return cl.List
}

type ConnectionV2 struct {
	Uuid              string  `json:"uuid" gorm:"primaryKey"`
	DatabaseType      string  `json:"databaseType"`
	CloudProviderType *string `json:"cloudProviderType"`
	CloudProviderUuid *string `json:"cloudProviderUuid"`

	Name string `json:"name"`

	Clusters          []Cluster         `json:"clusters" gorm:"foreignKey:ConnectionUuid"`
	ConnectionAccount ConnectionAccount `json:"connectionAccount" gorm:"embedded"`
	HideCredential    bool              `json:"hideCredential"`

	// More tabs
	AdditionalInfo        AdditionalInfo        `json:"additionalInfo" gorm:"embedded"`
	JustificationSettings JustificationSettings `json:"justificationSettings" gorm:"embedded"`

	SslSetting SslSetting `json:"sslSetting" gorm:"embedded"`
	SshSetting SshSetting `json:"sshSetting" gorm:"embedded"`

	ConnectionOwners []ConnectionOwner `json:"connectionOwners" gorm:"foreignKey:ObjectUuid"`

	AdvancedPrivilegeSetting []AdvancedPrivilegeSetting `json:"advancedPrivilegeSetting" gorm:"json"`

	Zones        []model.Zone `json:"zones" gorm:"json"`
	Ledger       bool         `json:"ledger"`
	VendorDetail VendorDetail `json:"vendorDetail" gorm:"json"`

	CreatedAt time.Time      `json:"createdAt"`
	CreatedBy model.Modifier `json:"createdBy" gorm:"json"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime:false"`
	UpdatedBy model.Modifier `json:"updatedBy" gorm:"json"`

	// Internal: HTTP response
	HttpResponse *resty.Response `json:"-" gorm:"-"`
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
	AccessEndTime       *string          `json:"accessEndTime"`
	AccessStartTime     *string          `json:"accessStartTime"`
	AuditEnabled        bool             `json:"auditEnabled"`
	DatabaseVersion     *string          `json:"databaseVersion"`
	Description         string           `json:"description"`
	DmlSnapshotEnabled  bool             `json:"dmlSnapshotEnabled"`
	LoginRules          LoginRules       `json:"loginRules" gorm:"embedded"`
	MaxDisplayRows      int              `json:"maxDisplayRows"`
	MaxExportRows       int              `json:"maxExportRows"`
	NetworkId           *string          `json:"networkId"`
	ProxyAuthType       *string          `json:"proxyAuthType"`
	ProxyEnabled        bool             `json:"proxyEnabled"`
	WeekdayAccessDenied model.StringList `json:"weekdayAccessDenied"`
}

type LoginRules struct {
	Interval         string `json:"interval"`
	MaxLoginFailures int    `json:"maxLoginFailures"`
}

type AdvancedPrivilegeSetting struct {
	PrivilegeUuid string `json:"privilegeUuid"`
	PrivilegeName string `json:"privilegeName"`
	DbAccountName string `json:"dbAccountName"`

	ConnectionUuid string `json:"-"`
}

type Cluster struct {
	Uuid            string       `json:"uuid" gorm:"primaryKey"`
	CloudIdentifier *string      `json:"cloudIdentifier"`
	Host            string       `json:"host"`
	Port            string       `json:"port"`
	ReplicationType string       `json:"replicationType"`
	Deleted         bool         `json:"deleted"`
	ConnectionUuid  string       `json:"-"`
	Connection      ConnectionV2 `json:"-" gorm:"foreignKey:ConnectionUuid"`
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
	Principal   string `json:"principal"`
	Realm       string `json:"realm"`
	ServiceName string `json:"serviceName"`
}

type ConnectionAccount struct {
	DbAccountName      *string            `json:"dbAccountName"`
	KerberosProtocol   *KerberosProtocol  `json:"kerberosProtocol" gorm:"json"`
	SecretStore        *model.SecretStore `json:"secretStore" gorm:"json"`
	SecretStoreEnabled bool               `json:"secretStoreEnabled"`
	Type               string             `json:"type"`
}

type OwnedBy struct {
	Uuid     string `json:"uuid" gorm:"primaryKey"`
	UserType string `json:"userType"`
	AuthType string `json:"authType"`
	LoginId  string `json:"loginId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

type ConnectionOwner struct {
	Uuid string `json:"uuid" gorm:"primaryKey"`

	RoleUuid string     `json:"-"`
	Role     model.Role `json:"role" gorm:"foreignKey:RoleUuid"`

	OwnerUuid string  `json:"-"`
	OwnedBy   OwnedBy `json:"ownedBy" gorm:"foreignKey:OwnerUuid"`

	ObjectUuid string `json:"objectUuid"`
}

type JustificationSettings struct {
	RequireExecuteReason      bool `json:"requireExecuteReason"`
	RequireExportDataReason   bool `json:"requireExportDataReason"`
	RequireExportSchemaReason bool `json:"requireExportSchemaReason"`
	RequireImportDataReason   bool `json:"requireImportDataReason"`
	RequireImportSchemaReason bool `json:"requireImportSchemaReason"`
}

type SshSetting struct {
	SshConfigName *string `json:"sshConfigName"`
	SshConfigUuid *string `json:"sshConfigUuid"`
	UseSsh        bool    `json:"useSsh"`
}

type SslSetting struct {
	SslConfigName *string `json:"sslConfigName"`
	SslConfigUuid *string `json:"sslConfigUuid"`
	UseSsl        bool    `json:"useSsl"`
}

type VendorDetail struct {
	DatabaseName  string  `json:"databaseName"`
	Charset       *string `json:"charset,omitempty"`
	Collation     *string `json:"collation,omitempty"`
	CloudRegion   *string `json:"cloudRegion,omitempty"`
	WorkGroupName *string `json:"workGroupName,omitempty"`
}
