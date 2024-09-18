package dac_access_control

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"qpc/entity/dac_connection"
	"qpc/entity/dac_privilege"
	"qpc/entity/user"
	"qpc/utils"
)

type DraftGrantRequest struct {
	UserQuery      string
	PrivilegeQuery string
	ClusterQuery   string
	Force          bool
	DryRun         bool

	users       []user.User
	privileges  []dac_privilege.Privilege
	connections []dac_connection.ConnectionV2
	clusters    []dac_connection.Cluster
}

func (dr *DraftGrantRequest) ToGrantRequest() *GrantRequest {
	return &GrantRequest{
		UserUuid:      dr.users[0].Uuid,
		ClusterUuid:   dr.clusters[0].Uuid,
		PrivilegeUuid: dr.privileges[0].Uuid,
		Force:         dr.Force,
	}
}

type GrantRequest struct {
	// Required: UUID of target user or a group
	UserUuid string `json:"-"`

	// Required: UUID of DB cluster
	ClusterUuid string `json:"clusterUuid"`

	// Required: Privilege UUID
	PrivilegeUuid string `json:"privilegeUuid"`

	// Optional: Whether to overwrite existing permissions (Default: false)
	Force bool `json:"force"`
}

func (r *GrantRequest) Post(server utils.QueryPieServerConfig) *GrantResponse {
	var grantResponse GrantResponse

	restClient := resty.New()
	uri := fmt.Sprintf("%s/api/external/v2/dac/access-controls/%s/grant",
		server.BaseURL,
		r.UserUuid)
	httpResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(server.AccessToken).
		SetBody(r).
		SetResult(&grantResponse).
		Post(uri)
	logrus.Debugf("Response: %v", httpResponse)
	if err != nil {
		logrus.Fatalf("Failed to grant access to DAC connection: %v", err)
	}
	grantResponse.HttpResponse = httpResponse
	return &grantResponse
}

type GrantResponse struct {
	// UUID of target user or group
	Uuid string `json:"uuid"`

	// Target type, USER or GROUP
	UserType string `json:"userType"`

	// Target name
	Name string `json:"name"`

	// Privilege-assigned DB cluster information
	MappedConnection MappedConnection `json:"mappedConnection"`

	// Internal: HTTP response
	HttpResponse *resty.Response `json:"-"`
}
