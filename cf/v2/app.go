package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/happytobi/cf-puppeteer/ui"
	"net/url"
	"time"
)

var (
	ErrAppNotFound = errors.New("application not found")
)

//AppRoutesResponse application response struct
type AppRoutesResponse struct {
	NextURL   string `json:"next_url,omitempty"`
	PrevURL   string `json:"prev_url,omitempty"`
	Resources []struct {
		Entity struct {
			AppsURL             string      `json:"apps_url"`
			DomainGUID          string      `json:"domain_guid"`
			DomainURL           string      `json:"domain_url"`
			Host                string      `json:"host"`
			Path                string      `json:"path"`
			Port                interface{} `json:"port"`
			RouteMappingsURL    string      `json:"route_mappings_url"`
			ServiceInstanceGUID interface{} `json:"service_instance_guid"`
			SpaceGUID           string      `json:"space_guid"`
			SpaceURL            string      `json:"space_url"`
		} `json:"entity"`
		Metadata struct {
			CreatedAt time.Time `json:"created_at"`
			GUID      string    `json:"guid"`
			UpdatedAt time.Time `json:"updated_at"`
			URL       string    `json:"url"`
		} `json:"metadata"`
	} `json:"resources"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type AppResourcesEntity struct {
	Metadata Metadata `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

type MetaDataEntity struct {
	AppResourcesEntity []AppResourcesEntity `json:"resources"`
}
type Metadata struct {
	GUID string `json:"guid"`
}
type Entity struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

//LoadAppRoutes
func (resource *ResourcesData) LoadAppRoutes(appGUID string) (*AppRoutesResponse, error) {
	ui.Say("LoadAppRoutes called")
	path := fmt.Sprintf(`/v2/apps/%s/routes`, appGUID)

	jsonResult, err := resource.cli.GetJSON(path)
	if err != nil {
		return nil, err
	}

	var response AppRoutesResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

//GetAppMetadata
func (resource *ResourcesData) GetAppMetadata(appName string) (*AppResourcesEntity, error) {
	space, err := resource.connection.GetCurrentSpace()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf(`v2/apps?q=name:%s&q=space_guid:%s`, url.QueryEscape(appName), space.Guid)
	jsonResult, err := resource.cli.GetJSON(path)

	if err != nil {
		return nil, err
	}

	var metaDataResponseEntity MetaDataEntity
	err = json.Unmarshal([]byte(jsonResult), &metaDataResponseEntity)

	if err != nil {
		ui.Failed("no response / parsable response from %s", path)
		return nil, err
	}

	if len(metaDataResponseEntity.AppResourcesEntity) == 0 {
		return nil, ErrAppNotFound
	}

	return &metaDataResponseEntity.AppResourcesEntity[0], nil
}
