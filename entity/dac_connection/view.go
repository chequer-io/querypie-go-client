package dac_connection

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
	"qpc/utils"
)

const scFmt = "%-36s  %-10s  %-5s  %-36s  %-8s  %-16s  %-16s\n"

func (cl *PagedConnectionV2List) Print() *PagedConnectionV2List {
	first := cl.GetPage().CurrentPage == 0
	last := !cl.GetPage().HasNext()

	if first {
		sc := &SummarizedConnectionV2{}
		sc.PrintHeader()
	}
	for _, conn := range cl.List {
		conn.Print()
	}
	if last {
		logrus.Infof("TotalElements: %v", cl.Page.TotalElements)
	}
	return cl
}

func (sc *SummarizedConnectionV2) PrintHeader() *SummarizedConnectionV2 {
	fmt.Printf(scFmt,
		"UUID",
		"DB_TYPE",
		"CLOUD",
		"NAME",
		"STATUS",
		"CREATED",
		"UPDATED",
	)
	return sc
}

func (sc *SummarizedConnectionV2) Print() *SummarizedConnectionV2 {
	logrus.Debug(sc)
	fmt.Printf(scFmt,
		sc.Uuid,
		sc.DatabaseType,
		utils.Optional(sc.CloudProviderType),
		sc.Name,
		sc.Status(),
		sc.ShortCreatedAt(),
		sc.ShortUpdatedAt(),
	)
	return sc
}

const connHeaderFmt = "  %-8s  %-6s  %-13s  %-12s  %-16s\n"
const connRowFmt = "  %-8d  %-6d  %-13s  %-12s  %-16s\n"

func (c *ConnectionV2) PrintHeader() *ConnectionV2 {
	fmt.Printf(connHeaderFmt,
		"CLUSTERS",
		"OWNERS",
		"ACNT_TYPE", // UIDPWD, SASL_KERBEROS
		"DB_ACNT_NAME",
		"DESCRIPTION",
	)
	return c
}

func (c *ConnectionV2) Print() *ConnectionV2 {
	fmt.Printf(connRowFmt,
		len(c.Clusters),
		len(c.ConnectionOwners),
		utils.Optional(c.ConnectionAccount.Type),
		utils.OptionalPtr(c.ConnectionAccount.DbAccountName),
		utils.Optional(c.AdditionalInfo.Description),
	)
	return c
}

func (c *ConnectionV2) printHttpRequestLineAndResponseStatus() {
	req := c.HttpResponse.Request.RawRequest
	res := c.HttpResponse.RawResponse
	fmt.Printf("%s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
	fmt.Printf("%s %s\n\n", res.Proto, res.Status)
}

func (c *ConnectionV2) PrintJson() *ConnectionV2 {
	if c == nil {
		return c
	} else if c.HttpResponse != nil {
		c.printHttpRequestLineAndResponseStatus()
		if c.HttpResponse.IsError() {
			fmt.Printf("%s\n", pretty.Pretty(c.HttpResponse.Body()))
			return c
		}
	}
	_json, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Failed to marshal ConnectionV2: %v", err)
		return c
	}
	fmt.Printf("%s\n", pretty.Pretty(_json))
	return c
}
