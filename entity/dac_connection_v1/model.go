package dac_connection_v1

import (
	"qpc/model"
)

type Cluster struct {
	Deleted bool   `json:"deleted"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Type    string `json:"type"`
	Uuid    string `json:"uuid"`
}

type KerberosProtocol struct {
	Principal   string `json:"principal"`
	Realm       string `json:"realm"`
	ServiceName string `json:"serviceName"`
}

type UsernamePassword struct {
	Username string `json:"username"`
}

type ConnectionAccount struct {
	KerberosProtocols struct {
		Common KerberosProtocol `json:"common"`
	} `json:"kerberosProtocols"`
	Type                       string                      `json:"type"`
	UseMultipleDatabaseAccount bool                        `json:"useMultipleDatabaseAccount"`
	UsernamePasswords          map[string]UsernamePassword `json:"usernamePasswords"`
}

type Role struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type User struct {
	Uuid         string `json:"uuid"`
	Id           int    `json:"id"`
	LoginId      string `json:"loginId"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Active       bool   `json:"active"`
	AuthProvider string `json:"authProvider"`
}

type ConnectionOwner struct {
	ObjectUuid string `json:"objectUuid"`
	ObjectType string `json:"objectType"`
	Uuid       string `json:"uuid"`
	User       User   `json:"user"`
	Role       Role   `json:"role"`
}

type Tag struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	Synchronized bool   `json:"synchronized"`
}

type Zone struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	IpBands   string `json:"ipBands"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Connection struct {
	Audited           bool              `json:"audited"`
	AuthenticationDb  string            `json:"authenticationDb"`
	AwsRegionValue    string            `json:"awsRegionValue"`
	CloudProviderUuid string            `json:"cloudProviderUuid"`
	Clusters          []Cluster         `json:"clusters"`
	ConnectionAccount ConnectionAccount `json:"connectionAccount"`
	ConnectionOwners  []ConnectionOwner `json:"connectionOwners"`
	DatabaseType      string            `json:"databaseType"`
	DatabaseVersion   string            `json:"databaseVersion"`
	Deleted           bool              `json:"deleted"`
	DmlSnapshot       bool              `json:"dmlSnapshot"`
	HideCredential    bool              `json:"hideCredential"`
	Ledger            bool              `json:"ledger"`
	ListenerType      string            `json:"listenerType"`
	MaxDisplayRows    int               `json:"maxDisplayRows"`
	MaxExportRows     int               `json:"maxExportRows"`
	Name              string            `json:"name"`
	NetworkId         string            `json:"networkId"`
	ProxyAuthType     string            `json:"proxyAuthType"`
	RoleName          string            `json:"roleName"`
	SchemaName        string            `json:"schemaName"`
	Tags              []Tag             `json:"tags"`
	UseProxy          bool              `json:"useProxy"`
	Uuid              string            `json:"uuid"`
	WarehouseName     string            `json:"warehouseName"`
	Zones             []Zone            `json:"zones"`
}

type ConnectionPagedList struct {
	List []Connection `json:"list"`
	Page model.Page   `json:"page"`
}
