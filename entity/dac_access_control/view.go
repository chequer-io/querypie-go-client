package dac_access_control

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func (acl *SummarizedAccessControlPagedList) Print() {
	first := acl.GetPage().CurrentPage == 0
	last := !acl.GetPage().HasNext()

	headerFmt := "%-36s  %-9s  %-9s  %-24s  %-24s  %-10s  %-3s  %-5s\n"
	rowFmt := "%-36s  %-9s  %-9s  %-24s  %-24s  %-10s  %-3d  %-5t\n"
	if first {
		logrus.Debugf("Page: %v", acl.GetPage())
		fmt.Printf(headerFmt,
			"UUID",
			"USER_TYPE",
			"AUTH_TYPE",
			"NAME",
			"MEMBERS",
			"ADMIN_ROLE",
			"CNT",
			"LINKED",
		)
	}

	for _, ac := range acl.GetList() {
		logrus.Debug(ac)
		fmt.Printf(rowFmt,
			ac.Uuid,
			ac.UserType,
			ac.AuthType,
			ac.Name,
			ac.MembersString(),
			ac.AdminRole,
			ac.LinkedCount,
			ac.Linked,
		)
	}

	if last {
		logrus.Infof("TotalElements: %v", acl.GetPage().TotalElements)
	}
}
