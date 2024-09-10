package dac_access_control

import (
	"fmt"
	"qpc/models"
)

type SummarizedAccessControl struct {
	Uuid        string   `json:"uuid"`
	UserType    string   `json:"userType"`
	AuthType    string   `json:"authType"`
	Name        string   `json:"name"`
	Members     []string `json:"members"`
	AdminRole   string   `json:"adminRole"`
	LinkedCount int      `json:"linkedCount"`
	Linked      bool     `json:"linked"`
}

func (sac SummarizedAccessControl) Status() string {
	if sac.Linked {
		return "linked"
	}
	return "-"
}

func (sac SummarizedAccessControl) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Type=%s, Name=%s }",
		sac.Uuid, sac.UserType, sac.Name,
	)
}

func (sac SummarizedAccessControl) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, UserType=%s, AuthType=%s, Name=%s, Members=%v, "+
			"AdminRole=%s, LinkedCount=%d, Linked=%t }",
		sac.Uuid, sac.UserType, sac.AuthType, sac.Name, sac.Members,
		sac.AdminRole, sac.LinkedCount, sac.Linked,
	)
}

type SummarizedAccessControlPagedList struct {
	models.PagedList[SummarizedAccessControl]
}

type Privilege struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type MappedConnection struct {
	Name          string               `json:"name"`
	ClusterUuid   string               `json:"clusterUuid"`
	DatabaseType  string               `json:"databaseType"`
	CloudProvider models.CloudProvider `json:"cloudProvider"`
	Privilege     Privilege            `json:"privilege"`
	Status        string               `json:"status"`
	Ledger        bool                 `json:"ledger"`
}

type AccessControl struct {
	Uuid              string             `json:"uuid"`
	UserType          string             `json:"userType"`
	AuthType          string             `json:"authType"`
	Name              string             `json:"name"`
	Members           []string           `json:"members"`
	MappedConnections []MappedConnection `json:"mappedConnections"`
	AdminRole         string             `json:"adminRole"`
	LinkedCount       int                `json:"linkedCount"`
	Linked            bool               `json:"linked"`
}

func (ac AccessControl) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Type=%s, Name=%s }",
		ac.Uuid, ac.UserType, ac.Name,
	)
}

func (ac AccessControl) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, UserType=%s, AuthType=%s, Name=%s, Members=%v, "+
			"AdminRole=%s, LinkedCount=%d, Linked=%t }",
		ac.Uuid, ac.UserType, ac.AuthType, ac.Name, ac.Members,
		ac.AdminRole, ac.LinkedCount, ac.Linked,
	)
}

type AccessControlPagedList struct {
	models.PagedList[AccessControl]
}
