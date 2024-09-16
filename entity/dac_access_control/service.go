package dac_access_control

import (
	"github.com/sirupsen/logrus"
	"qpc/entity/dac_connection"
	"qpc/entity/dac_privilege"
	"qpc/entity/user"
	"regexp"
)

func (dr *DraftGrantRequest) LookUpEntities() *DraftGrantRequest {
	(&user.User{}).FindByLoginIdOrEmailOrUuid(dr.UserQuery, &dr.users)
	(&dac_privilege.Privilege{}).FindByNameOrUuid(dr.PrivilegeQuery, &dr.privileges)
	if isHostnameAndPort(dr.ClusterQuery) {
		(&dac_connection.Cluster{}).FindByHostAndPort(dr.ClusterQuery, &dr.clusters)
	} else {
		var cluster = (&dac_connection.Cluster{}).FindByUuid(dr.ClusterQuery)
		if cluster != nil {
			dr.clusters = append(dr.clusters, *cluster)
		}
	}
	return dr
}

func isHostnameAndPort(clusterQuery string) bool {
	// Define the regular expression pattern for host:port
	// Host: FQDN or PQDN (e.g., example.com, sub.example.com)
	// Port: 1-65535
	pattern := `^(\w+\.)+[a-zA-Z]{2,}:[0-9]{1,5}$`
	re := regexp.MustCompile(pattern)

	// Check if the input string matches the pattern
	return re.MatchString(clusterQuery)
}

func (dr *DraftGrantRequest) Validate() bool {
	logrus.Warnf("TODO(JK): Not Yet Implemented")
	return false
}
