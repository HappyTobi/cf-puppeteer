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
	CreateRoute(spaceGUID string, domainGUID string, host string) (*V2RouteResponse, error)
	FindServiceInstances(serviceNames []string, spaceGUID string) ([]string, error)
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

type V2Route struct {
	DomainGUID string `json:"domain_guid"`
	SpaceGUID  string `json:"space_guid"`
	Port       int    `json:"port,omitempty"`
	Host       string `json:"host,omitempty"`
}

//V2RouteResponse
type V2RouteResponse struct {
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
func (resource *ResourcesData) CreateRoute(spaceGUID string, domainGUID string, host string) (*V2RouteResponse, error) {
	path := fmt.Sprintf(`/v2/routes`)

	var v2Route V2Route
	v2Route.DomainGUID = domainGUID
	v2Route.SpaceGUID = spaceGUID
	v2Route.Host = host

	appJSON, err := json.Marshal(v2Route)
	if err != nil {
		return nil, err
	}

	if resource.TraceLogging {
		fmt.Printf("send POST to route: %s with body:\n", path)
		prettyPrintJSON(string(appJSON))
	}

	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("response from call to path: %s was:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V2RouteResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type V2ServicesInstanceResponse struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	PrevURL      interface{} `json:"prev_url"`
	NextURL      interface{} `json:"next_url"`
	Resources    []struct {
		Metadata struct {
			GUID      string    `json:"guid"`
			URL       string    `json:"url"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"metadata"`
		Entity struct {
			Name        string `json:"name"`
			Credentials struct {
			} `json:"credentials"`
			ServicePlanGUID string      `json:"service_plan_guid"`
			SpaceGUID       string      `json:"space_guid"`
			GatewayData     interface{} `json:"gateway_data"`
			DashboardURL    interface{} `json:"dashboard_url"`
			Type            string      `json:"type"`
			LastOperation   struct {
				Type        string    `json:"type"`
				State       string    `json:"state"`
				Description string    `json:"description"`
				UpdatedAt   time.Time `json:"updated_at"`
				CreatedAt   time.Time `json:"created_at"`
			} `json:"last_operation"`
			Tags                         []interface{} `json:"tags"`
			ServiceGUID                  string        `json:"service_guid"`
			SpaceURL                     string        `json:"space_url"`
			ServicePlanURL               string        `json:"service_plan_url"`
			ServiceBindingsURL           string        `json:"service_bindings_url"`
			ServiceKeysURL               string        `json:"service_keys_url"`
			RoutesURL                    string        `json:"routes_url"`
			ServiceURL                   string        `json:"service_url"`
			SharedFromURL                string        `json:"shared_from_url"`
			SharedToURL                  string        `json:"shared_to_url"`
			ServiceInstanceParametersURL string        `json:"service_instance_parameters_url"`
		} `json:"entity"`
	} `json:"resources"`
}

//FindServiceInstances return guids for all Serviceinsances
func (resource *ResourcesData) FindServiceInstances(serviceNames []string, spaceGUID string) ([]string, error) {
	path := fmt.Sprintf(`/v2/spaces/%s/service_instances?return_user_provided_service_instances=true`, spaceGUID)
	if resource.TraceLogging {
		fmt.Printf("get service instances by name: %s - req path: %s:\n", serviceNames, path)
	}

	var serviceGUIDs []string
	for _, serviceName := range serviceNames {
		serviceQueryPath := fmt.Sprintf(`%s&q=name:%s`, path, serviceName)
		result, _ := resource.Connection.CliCommandWithoutTerminalOutput("curl", serviceQueryPath, "-X", "GET", "-H", "Content-type: application/json")
		jsonResp := strings.Join(result, "")
		if resource.TraceLogging {
			fmt.Printf("response from http call to path: %s for service: %s - was:\n", serviceQueryPath, serviceName)
			prettyPrintJSON(jsonResp)
		}

		var v2ServicesInstanceResponse V2ServicesInstanceResponse
		err := json.Unmarshal([]byte(jsonResp), &v2ServicesInstanceResponse)
		if err != nil {
			return nil, err
		}

		for _, ent := range v2ServicesInstanceResponse.Resources {
			serviceGUIDs = append(serviceGUIDs, ent.Entity.ServiceGUID)
		}

		return serviceGUIDs, nil

	}

	/*var response V2SharedDomainResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}*/

	return serviceGUIDs, nil

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
