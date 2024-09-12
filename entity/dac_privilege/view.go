package dac_privilege

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"qpc/utils"
)

func (pl *PrivilegePagedList) Print() {
	first := pl.GetPage().CurrentPage == 0
	last := !pl.GetPage().HasNext()

	headerFmt := "%-36s  %-20s  %-28s  %-16s  %-12s  %-10s\n"
	rowFmt := "%-36s  %-20s  %-28s  %-16s  %-12s  %-10s\n"
	if first {
		logrus.Debugf("Page: %v", pl.GetPage())
		fmt.Printf(headerFmt,
			"UUID",
			"NAME",
			"PRIVILEGE_TYPES",
			"DESCRIPTION",
			"PRIVILEGE_VENDOR",
			"STATUS",
		)
	}

	for _, p := range pl.GetList() {
		p.PrivilegeTypesStr = utils.JsonFromStringArray(p.PrivilegeTypes)
		logrus.Debug(p)
		fmt.Printf(rowFmt,
			p.Uuid,
			p.Name,
			p.PrivilegeTypesStr,
			p.Description,
			p.PrivilegeVendor,
			p.Status,
		)
	}

	if last {
		logrus.Infof("TotalElements: %v", pl.GetPage().TotalElements)
	}
}
