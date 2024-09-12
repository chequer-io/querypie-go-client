package dac_access_control

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"qpc/model"
	"qpc/utils"
)

type SummarizedAccessControl struct {
	Uuid        string   `json:"uuid"`
	UserType    string   `json:"userType"`
	AuthType    string   `json:"authType"`
	Name        string   `json:"name"`
	Members     []string `json:"members" gorm:"-"`
	MembersStr  string   `json:"-"`
	AdminRole   string   `json:"adminRole"`
	LinkedCount int      `json:"linkedCount"`
	Linked      bool     `json:"linked"`
}

func (sac *SummarizedAccessControl) Status() string {
	if sac.Linked {
		return "linked"
	}
	return "-"
}

func (sac *SummarizedAccessControl) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Type=%s, Name=%s }",
		sac.Uuid, sac.UserType, sac.Name,
	)
}

func (sac *SummarizedAccessControl) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, UserType=%s, AuthType=%s, Name=%s, Members=%v, MembersStr=%v "+
			"AdminRole=%s, LinkedCount=%d, Linked=%t }",
		sac.Uuid, sac.UserType, sac.AuthType, sac.Name, sac.Members, sac.MembersStr,
		sac.AdminRole, sac.LinkedCount, sac.Linked,
	)
}

func (sac *SummarizedAccessControl) PopulateMemberStr() {
	sac.MembersStr = utils.JsonFromStringArray(sac.Members)
	logrus.Debugf("Populated MembersStr: %v from Members: %v", sac.MembersStr, sac.Members)
	return
}

func (sac *SummarizedAccessControl) PopulateMembers() {
	sac.Members = utils.StringArrayFromJson(sac.MembersStr)
	logrus.Debugf("Populated Members: %v from MembersStr: %v", sac.Members, sac.MembersStr)
	return
}

func (sac *SummarizedAccessControl) MembersString() string {
	if sac.Members == nil {
		return "-"
	} else if len(sac.Members) == 0 {
		return "[]"
	} else if len(sac.Members) == 1 {
		return fmt.Sprintf("[%v]", sac.Members)
	} else {
		return fmt.Sprintf("[%v, +%d]", sac.Members[0], len(sac.Members)-1)
	}
}

type SummarizedAccessControlPagedList struct {
	List []SummarizedAccessControl `json:"list"`
	Page model.Page                `json:"page"`
}

func (acl *SummarizedAccessControlPagedList) GetPage() model.Page {
	return acl.Page
}

func (acl *SummarizedAccessControlPagedList) GetList() []SummarizedAccessControl {
	return acl.List
}

type Privilege struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type MappedConnection struct {
	Name          string              `json:"name"`
	ClusterUuid   string              `json:"clusterUuid"`
	DatabaseType  string              `json:"databaseType"`
	CloudProvider model.CloudProvider `json:"cloudProvider"`
	Privilege     Privilege           `json:"privilege"`
	Status        string              `json:"status"`
	Ledger        bool                `json:"ledger"`
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

func (ac *AccessControl) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Type=%s, Name=%s }",
		ac.Uuid, ac.UserType, ac.Name,
	)
}

func (ac *AccessControl) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, UserType=%s, AuthType=%s, Name=%s, Members=%v, "+
			"AdminRole=%s, LinkedCount=%d, Linked=%t }",
		ac.Uuid, ac.UserType, ac.AuthType, ac.Name, ac.Members,
		ac.AdminRole, ac.LinkedCount, ac.Linked,
	)
}

type AccessControlPagedList struct {
	List []AccessControl `json:"list"`
	Page model.Page      `json:"page"`
}

func (acl *AccessControlPagedList) GetPage() model.Page {
	return acl.Page
}

func (acl *AccessControlPagedList) GetList() []AccessControl {
	return acl.List
}
