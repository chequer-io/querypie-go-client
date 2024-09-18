package user

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Please do not print extra white spaces in the last column.
const userFmt = "%-36s  %-22s  %-22s  %-18s  %-8s  %-8s %-16s  %s\n"

func (u *User) PrintHeader() *User {
	fmt.Printf(userFmt,
		"UUID",
		"LOGIN_ID",
		"EMAIL",
		"NAME",
		"STATUS",
		"MORE",
		"CREATED",
		"UPDATED",
	)
	return u
}

func (u *User) Print() *User {
	logrus.Debug(u)
	fmt.Printf(userFmt,
		u.Uuid,
		u.LoginId,
		u.Email,
		u.Name,
		u.Status,
		u.StatusMore(),
		u.ShortCreatedAt(),
		u.ShortUpdatedAt(),
	)
	return u
}
