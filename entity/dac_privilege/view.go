package dac_privilege

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"qpc/utils"
)

func (pl *PrivilegePagedList) Print() {
	first := pl.GetPage().CurrentPage == 0
	last := !pl.GetPage().HasNext()

	headerFmt := "%-36s  %-20s  %-40s  %-18s  %-8s  %-10s\n"
	rowFmt := "%-36s  %-20s  %-40s  %-18s  %-8s  %-10s\n"
	if first {
		logrus.Debugf("Page: %v", pl.GetPage())
		fmt.Printf(headerFmt,
			"UUID",
			"NAME",
			"PRIVILEGE_TYPES",
			"DESCRIPTION",
			"VENDOR",
			"STATUS",
		)
	}

	for _, p := range pl.GetList() {
		logrus.Debug(p)
		fmt.Printf(rowFmt,
			p.Uuid,
			p.Name,
			p.PrivilegeTypes.Ellipsis(),
			utils.Optional(p.Description),
			p.PrivilegeVendor,
			p.Status,
		)
	}

	if last {
		logrus.Infof("TotalElements: %v", pl.GetPage().TotalElements)
	}
}
