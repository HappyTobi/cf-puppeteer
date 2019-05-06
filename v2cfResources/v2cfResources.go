package v2cfResources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
)

//V2CfResourcesInterface interface for cfResource usage
type V2CfResourcesInterface interface {
	LoadAppRoutes(appGUID string) (*V2AppRoutesResponse, error)
	LoadSharedDomains(domainGUID string) (*V2SharedDomainResponse, error)
}

//ResourcesData struct to hold important instances to run push
type ResourcesData struct {
	Connection   plugin.CliConnection
	TraceLogging bool
}

//NewResources create a new instance of CFResources
func NewResources(conn plugin.CliConnection, traceLogging bool) *ResourcesData {
	return &ResourcesData{
		Connection:   conn,
		TraceLogging: traceLogging,
	}
}

//V2Apps application struct
type V2Apps struct {
	Name      string `json:"name"`
	Lifecycle struct {
		LifecycleType string `json:"type"`
		LifecycleData struct {
			Buildpacks []string `json:"buildpacks,omitempty"`
			Stack      string   `json:"stack,omitempty"`
		} `json:"data"`
	} `json:"lifecycle"`
	Relationships struct {
		Space struct {
			Data struct {
				GUID string `json:"guid"`
			} `json:"data"`
		} `json:"space"`
	} `json:"relationships"`
	EnvironmentVariables struct {
		Vars map[string]string `json:"var,omitempty"`
	} `json:"environment_variables,omitempty"`
}

//V2AppRoutesResponse application response struct
type V2AppRoutesResponse struct {
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

//LoadAppRoutes
func (resource *ResourcesData) LoadAppRoutes(appGUID string) (*V2AppRoutesResponse, error) {
	path := fmt.Sprintf(`/v2/apps/%s/routes`, appGUID)
	if resource.TraceLogging {
		fmt.Printf("get routes: %s:\n", path)
	}
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")
	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("response from http call to path: %s was:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V2AppRoutesResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type V2SharedDomainResponse struct {
	Entity struct {
		Internal        bool        `json:"internal"`
		Name            string      `json:"name"`
		RouterGroupGUID interface{} `json:"router_group_guid"`
		RouterGroupType interface{} `json:"router_group_type"`
	} `json:"entity"`
	Metadata struct {
		CreatedAt time.Time `json:"created_at"`
		GUID      string    `json:"guid"`
		UpdatedAt time.Time `json:"updated_at"`
		URL       string    `json:"url"`
	} `json:"metadata"`
}

//LoadSharedDomains
func (resource *ResourcesData) LoadSharedDomains(domainGUID string) (*V2SharedDomainResponse, error) {
	path := fmt.Sprintf(`/v2/shared_domains/%s`, domainGUID)
	if resource.TraceLogging {
		fmt.Printf("get routes: %s:\n", path)
	}
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")
	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("response from http call to path: %s was:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V2SharedDomainResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PrettyPrintJSON takes the given JSON string, makes it pretty, and prints it out.
func prettyPrintJSON(jsonUgly string) error {
	jsonPretty := &bytes.Buffer{}
	err := json.Indent(jsonPretty, []byte(jsonUgly), "", "    ")

	if err != nil {
		return err
	}

	fmt.Println(jsonPretty.String())

	return nil
}
