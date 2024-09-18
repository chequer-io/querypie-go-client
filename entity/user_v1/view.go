package user_v1

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Please do not print extra white spaces in the last column.
const userFmt = "%-36s  %-24s  %-24s  %-20s  %-8s  %-16s  %s\n"

func (u *UserV1) PrintHeader() *UserV1 {
	fmt.Printf(userFmt,
		"UUID",
		"LOGIN_ID",
		"EMAIL",
		"NAME",
		"STATUS",
		"CREATED",
		"UPDATED",
	)
	return u
}

func (u *UserV1) Print() *UserV1 {
	logrus.Debug(u)
	fmt.Printf(userFmt,
		u.Uuid,
		u.LoginId,
		u.Email,
		u.Name,
		u.Status(),
		u.ShortCreatedAt(),
		u.ShortUpdatedAt(),
	)
	return u
}
