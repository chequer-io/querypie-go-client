package dac_privilege

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"qpc/utils"
)

const privilegeHeaderFmt = "%-36s  %-20s  %-40s  %-18s  %-8s  %-10s\n"
const privilegeRowFmt = "%-36s  %-20s  %-40s  %-18s  %-8s  %-10s\n"

func (p *Privilege) PrintHeader() *Privilege {
	fmt.Printf(privilegeHeaderFmt,
		"UUID",
		"NAME",
		"PRIVILEGE_TYPES",
		"DESCRIPTION",
		"VENDOR",
		"STATUS",
	)
	return p
}

func (p *Privilege) Print() *Privilege {
	logrus.Debug(p)
	fmt.Printf(privilegeRowFmt,
		p.Uuid,
		p.Name,
		p.PrivilegeTypes.Ellipsis(),
		utils.Optional(p.Description),
		p.PrivilegeVendor,
		p.Status,
	)
	return p
}
