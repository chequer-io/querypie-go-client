package dac_access_control

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
	"qpc/entity/dac_connection"
	"qpc/entity/dac_privilege"
	"qpc/entity/user"
	"qpc/utils"
)

const sacHeaderFmt = "%-36s  %-9s  %-9s  %-24s  %-24s  %-10s  %-3s  %-5s\n"
const sacRowFmt = "%-36s  %-9s  %-9s  %-24s  %-24s  %-10s  %-3d  %-5t\n"

func (sac *SummarizedAccessControl) PrintHeader() *SummarizedAccessControl {
	fmt.Printf(sacHeaderFmt,
		"UUID",
		"USER_TYPE",
		"AUTH_TYPE",
		"NAME",
		"MEMBERS",
		"ADMIN_ROLE",
		"CNT",
		"LINKED",
	)
	return sac
}

func (sac *SummarizedAccessControl) Print() *SummarizedAccessControl {
	logrus.Debug(sac)
	fmt.Printf(sacRowFmt,
		sac.Uuid,
		sac.UserType,
		sac.AuthType,
		sac.Name,
		sac.Members.Ellipsis(1),
		sac.AdminRole,
		sac.LinkedCount,
		sac.Linked,
	)
	return sac
}

func (dr *DraftGrantRequest) Print() *DraftGrantRequest {
	fmt.Printf("GRANT REQUEST\n\n")
	fmt.Printf("Users matched: %d\n", len(dr.users))
	(&user.User{}).PrintHeader()
	for _, u := range dr.users {
		u.Print()
	}
	fmt.Println()
	fmt.Printf("\nPrivileges matched: %d\n", len(dr.privileges))
	(&dac_privilege.Privilege{}).PrintHeader()
	for _, p := range dr.privileges {
		p.Print()
	}
	fmt.Println()
	fmt.Printf("\nClusters matched: %d\n", len(dr.clusters))
	(&dac_connection.Cluster{}).PrintHeaderWithConnection()
	for _, c := range dr.clusters {
		c.PrintWithConnection()
	}
	fmt.Println()
	return dr
}

func (r *GrantResponse) Print() {
	utils.PrintHttpRequestLineAndResponseStatus(r.HttpResponse)
	if r.HttpResponse.StatusCode() != 200 {
		fmt.Printf("%s\n",
			pretty.Pretty(r.HttpResponse.Body()),
		)
		return
	}
	format := "%-36s  %-10s  %-5s  %-36s  %-8s  %-16s  %-16s\n"
	fmt.Printf(format,
		"UUID",
		"USER_TYPE",
		"NAME",
		"CLOUD_PROVIDER",
		"CLUSTER_UUID",
		"DB_TYPE",
		"STATUS",
	)
	// Print TODO(JK): Print more on MappedConnection
	fmt.Printf(format,
		r.Uuid,
		r.UserType,
		r.Name,
		r.MappedConnection.CloudProvider.Name,
		r.MappedConnection.ClusterUuid,
		r.MappedConnection.DatabaseType,
		r.MappedConnection.Status,
	)
}
