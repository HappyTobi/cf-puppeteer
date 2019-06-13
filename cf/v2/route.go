package v2

import (
	"encoding/json"
	"fmt"
	"time"
)

type Route struct {
	DomainGUID string `json:"domain_guid"`
	SpaceGUID  string `json:"space_guid"`
	Port       int    `json:"port,omitempty"`
	Host       string `json:"host,omitempty"`
}

//RouteResponse
type RouteResponse struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Host                string      `json:"host"`
		Path                string      `json:"path"`
		DomainGUID          string      `json:"domain_guid"`
		SpaceGUID           string      `json:"space_guid"`
		ServiceInstanceGUID interface{} `json:"service_instance_guid"`
		Port                int         `json:"port"`
		DomainURL           string      `json:"domain_url"`
		SpaceURL            string      `json:"space_url"`
		AppsURL             string      `json:"apps_url"`
		RouteMappingsURL    string      `json:"route_mappings_url"`
	} `json:"entity"`
}

//LoadSharedDomains
func (resource *ResourcesData) CreateRoute(spaceGUID string, domainGUID string, host string) (*RouteResponse, error) {
	path := fmt.Sprintf(`/v2/routes`)

	var v2Route Route
	v2Route.DomainGUID = domainGUID
	v2Route.SpaceGUID = spaceGUID
	v2Route.Host = host

	appJSON, err := json.Marshal(v2Route)
	if err != nil {
		return nil, err
	}

	jsonResult, err := resource.cli.PostJSON(path, string(appJSON))

	if err != nil {
		return nil, err
	}

	var response RouteResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
