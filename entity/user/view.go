package user

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

const userFmt = "%-36s  %-22s  %-22s  %-18s  %-8s  %-8s %-16s  %-16s\n"

func (pul *PagedUserList) Print() {
	first := pul.GetPage().CurrentPage == 0
	last := !pul.GetPage().HasNext()

	if first {
		logrus.Debugf("Page: %v", pul.GetPage())
		sc := &User{}
		sc.PrintHeader()
	}

	for _, u := range pul.GetList() {
		logrus.Debug(u)
		u.Print()
	}

	if last {
		logrus.Infof("TotalElements: %v", pul.GetPage().TotalElements)
	}
}

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
