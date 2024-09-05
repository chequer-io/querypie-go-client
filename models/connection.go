package models

type Connection struct {
	Admin struct {
		Username string `json:"username"`
	} `json:"admin"`
	Common struct {
		Username string `json:"username"`
	} `json:"common"`
	Proxy struct {
		Username string `json:"username"`
	} `json:"proxy"`
	ConnectionOwners []ConnectionOwner `json:"connectionOwners"`
	DatabaseType     string            `json:"databaseType"`
	DatabaseVersion  string            `json:"databaseVersion"`
	Deleted          bool              `json:"deleted"`
	DmlSnapshot      bool              `json:"dmlSnapshot"`
	HideCredential   bool              `json:"hideCredential"`
	Ledger           bool              `json:"ledger"`
	ListenerType     string            `json:"listenerType"`
	MaxDisplayRows   int               `json:"maxDisplayRows"`
	MaxExportRows    int               `json:"maxExportRows"`
	Name             string            `json:"name"`
	NetworkId        string            `json:"networkId"`
	ProxyAuthType    string            `json:"proxyAuthType"`
	RoleName         string            `json:"roleName"`
	SchemaName       string            `json:"schemaName"`
	Tags             []Tag             `json:"tags"`
	UseProxy         bool              `json:"useProxy"`
	Uuid             string            `json:"uuid"`
	WarehouseName    string            `json:"warehouseName"`
	Zones            []Zone            `json:"zones"`
}
