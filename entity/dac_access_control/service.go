package dac_access_control

import (
	"qpc/entity/dac_connection"
	"qpc/entity/dac_privilege"
	"qpc/entity/user"
	"regexp"
	"strings"
)

func (dr *DraftGrantRequest) LookUpEntities() *DraftGrantRequest {
	// Look up users.
	// It is rather simple that we can use FindByLoginIdOrEmailOrUuid().
	(&user.User{}).FindByLoginIdOrEmailOrUuid(dr.UserQuery, &dr.users)

	// Look up privileges.
	(&dac_privilege.Privilege{}).FindByNameOrUuid(dr.PrivilegeQuery, &dr.privileges)

	// Look up clusters.
	// We need to check if the query is a hostname:port, a UUID, or a cloud identifier.
	if isHostnameAndPort(dr.ClusterQuery) {
		(&dac_connection.Cluster{}).FindByHostAndPort(dr.ClusterQuery, &dr.clusters)
	} else if isUuid(dr.ClusterQuery) {
		var found = (&dac_connection.Cluster{}).FirstByUuid(dr.ClusterQuery)
		if found != nil {
			dr.clusters = append(dr.clusters, *found)
		}
	} else {
		(&dac_connection.Cluster{}).FindByCloudIdentifier(dr.ClusterQuery, &dr.clusters)
	}

	// Look up connections, if clusters are not found.
	if len(dr.clusters) == 0 {
		(&dac_connection.ConnectionV2{}).FindByNameOrUuid(dr.ClusterQuery, &dr.connections)
		// If connections are found, add their clusters to the list.
		for _, connection := range dr.connections {
			for _, cluster := range connection.Clusters {
				dr.clusters = append(dr.clusters, cluster)
			}
		}
	}
	return dr
}

func isHostnameAndPort(query string) bool {
	// Define the regular expression pattern for host:port
	// Host: FQDN or PQDN (e.g., example.com, sub.example.com)
	// Port: 1-65535
	pattern := `(?i)^(\w+\.)+[a-z]{2,}:[0-9]{1,5}$`
	re := regexp.MustCompile(pattern)

	// Check if the input string matches the pattern
	return re.MatchString(query)
}

func isUuid(query string) bool {
	// Define the regular expression pattern for UUID
	// UUID: 8-4-4-4-12 alphanumeric characters
	pattern := `(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	re := regexp.MustCompile(pattern)

	// Check if the input string matches the pattern
	return re.MatchString(query)
}

func (dr *DraftGrantRequest) Validate(
	successCase func(),
	failedCase func(string),
) bool {
	var reasons []string
	if len(dr.users) == 0 {
		reasons = append(reasons, "User not found")
	} else if len(dr.users) > 1 {
		reasons = append(reasons, "Multiple users found")
	}
	if len(dr.privileges) == 0 {
		reasons = append(reasons, "Privilege not found")
	} else if len(dr.privileges) > 1 {
		reasons = append(reasons, "Multiple privileges found")
	}
	if len(dr.clusters) == 0 {
		reasons = append(reasons, "Cluster not found")
	} else if len(dr.clusters) > 1 {
		reasons = append(reasons, "Multiple clusters found")
	}

	if len(reasons) > 0 {
		failedCase(strings.Join(reasons, ", "))
		return false
	} else {
		successCase()
		return true
	}
}
