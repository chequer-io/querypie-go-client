package dac_connection

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
	"qpc/utils"
)

// Please do not print extra white spaces in the last column.
const scFmt = "%-36s  %-10s  %-5s  %-36s  %-8s  %-16s  %s\n"

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

func (sc *SummarizedConnectionV2) PrintYamlHeader(silent bool) *SummarizedConnectionV2 {
	if silent {
		return sc
	}
	fmt.Printf("# begin of summarized connection v2\n")
	return sc
}

func (sc *SummarizedConnectionV2) PrintYaml(silent bool) *SummarizedConnectionV2 {
	logrus.Debug(sc)
	if silent {
		return sc
	}
	slice := []SummarizedConnectionV2{*sc}
	_yaml, err := yaml.Marshal(slice)
	if err == nil {
		fmt.Print(string(_yaml))
	} else {
		logrus.Errorf("Failed to marshal SummarizedConnectionV2: %v", err)
	}
	return sc
}

func (sc *SummarizedConnectionV2) PrintYamlFooter(silent bool, count int) *SummarizedConnectionV2 {
	if silent {
		return sc
	}
	fmt.Printf("# end of summarized connection v2, total elements: %d\n", count)
	logrus.Infof("TotalElements: %v", count)
	return sc
}

// Please do not print extra white spaces in the last column.
const connHeaderFmt = "%s%-8s  %-6s  %-13s  %-12s  %s\n"
const connRowFmt = "%s%-8d  %-6d  %-13s  %-12s  %s\n"

func (c *ConnectionV2) PrintHeader(prefix string) *ConnectionV2 {
	fmt.Printf(connHeaderFmt,
		prefix,
		"CLUSTERS",
		"OWNERS",
		"ACNT_TYPE", // UIDPWD, SASL_KERBEROS
		"DB_ACNT_NAME",
		"DESCRIPTION",
	)
	return c
}

func (c *ConnectionV2) Print(prefix string) *ConnectionV2 {
	fmt.Printf(connRowFmt,
		prefix,
		len(c.Clusters),
		len(c.ConnectionOwners),
		utils.Optional(c.ConnectionAccount.Type),
		utils.OptionalPtr(c.ConnectionAccount.DbAccountName),
		utils.Optional(c.AdditionalInfo.Description),
	)
	return c
}

func (c *ConnectionV2) PrintJson() *ConnectionV2 {
	if c == nil {
		return c
	}
	_json, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Failed to marshal ConnectionV2: %v", err)
		return c
	}
	fmt.Printf("%s\n", pretty.Pretty(_json))
	return c
}

func (c *ConnectionV2) PrintYamlHeader(silent bool) *ConnectionV2 {
	if silent {
		return c
	}
	fmt.Printf("# begin of connection v2\n")
	return c
}

func (c *ConnectionV2) PrintYaml(silent bool) *ConnectionV2 {
	logrus.Debug(c)
	if silent {
		return c
	}
	slice := []ConnectionV2{*c}
	_yaml, err := yaml.Marshal(slice)
	if err == nil {
		fmt.Print(string(_yaml))
	} else {
		logrus.Errorf("Failed to marshal ConnectionV2: %v", err)
	}
	return c
}

func (c *ConnectionV2) PrintYamlFooter(silent bool, count int) *ConnectionV2 {
	if silent {
		return c
	}
	fmt.Printf("# end of connection v2, total elements: %d\n", count)
	logrus.Infof("TotalElements: %v", count)
	return c
}

// Please do not print extra white spaces in the last column.
const clusterFmt = "%s%-36s  %-24s  %-5s  %-8s  %-16s  %s\n"

func (c *Cluster) PrintHeader(prefix string) *Cluster {
	fmt.Printf(clusterFmt,
		prefix,
		"UUID",
		"HOST",
		"PORT",
		"TYPE",
		"CLOUD_IDENTIFIER",
		"STATUS",
	)
	return c
}

func (c *Cluster) Print(prefix string) *Cluster {
	fmt.Printf(clusterFmt,
		prefix,
		c.Uuid,
		utils.Optional(c.Host),
		utils.Optional(c.Port),
		c.ReplicationType,
		utils.OptionalPtr(c.CloudIdentifier),
		c.Status(),
	)
	return c
}

// Please do not print extra white spaces in the last column.
const clusterFmtWithConnection = "%-36s  %-24s  %-5s  %-8s  %-16s  %-8s  %-36s  %s\n"

func (c *Cluster) PrintHeaderWithConnection() *Cluster {
	fmt.Printf(clusterFmtWithConnection,
		"UUID",
		"HOST",
		"PORT",
		"TYPE",
		"CLOUD_IDENTIFIER",
		"STATUS",
		"CONNECTION_UUID",
		"CONNECTION_NAME",
	)
	return c
}

func (c *Cluster) PrintWithConnection() *Cluster {
	fmt.Printf(clusterFmtWithConnection,
		c.Uuid,
		utils.Optional(c.Host),
		utils.Optional(c.Port),
		c.ReplicationType,
		utils.OptionalPtr(c.CloudIdentifier),
		c.Status(),
		utils.Optional(c.ConnectionUuid),
		utils.Optional(c.Connection.Name),
	)
	return c
}
