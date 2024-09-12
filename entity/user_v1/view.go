package user_v1

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func (pul *PagedUserV1List) Print() {
	first := pul.GetPage().CurrentPage == 0
	last := !pul.GetPage().HasNext()

	format := "%-36s  %-24s  %-24s  %-20s  %-8s  %-16s  %-16s\n"
	if first {
		logrus.Debugf("Page: %v", pul.GetPage())
		fmt.Printf(format,
			"UUID",
			"LOGIN_ID",
			"EMAIL",
			"NAME",
			"STATUS",
			"CREATED",
			"UPDATED",
		)
	}

	for _, u := range pul.GetList() {
		logrus.Debug(u)
		fmt.Printf(format,
			u.Uuid,
			u.LoginId,
			u.Email,
			u.Name,
			u.Status(),
			u.ShortCreatedAt(),
			u.ShortUpdatedAt(),
		)
	}

	if last {
		logrus.Infof("TotalElements: %v", pul.GetPage().TotalElements)
	}
}
