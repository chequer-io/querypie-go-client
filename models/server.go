package models

type Server struct {
	Host    string `json:"host"`
	Name    string `json:"name"`
	OsType  string `json:"osType"`
	SshPort int    `json:"sshPort"`
	Tags    struct {
		CustomTags   []Tag `json:"customTags"`
		ProviderTags []Tag `json:"providerTags"`
	} `json:"tags"`
	Uuid string `json:"uuid"`
}
