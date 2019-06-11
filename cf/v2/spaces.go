package v2

import (
	"encoding/json"
	"fmt"
	"time"
)

type ServicesInstanceResponse struct {
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
	/*if resource.TraceLogging {
		fmt.Printf("get service instances by name: %s - req path: %s:\n", serviceNames, path)
	}*/

	var serviceGUIDs []string
	for _, serviceName := range serviceNames {
		serviceQueryPath := fmt.Sprintf(`%s&q=name:%s`, path, serviceName)
		jsonResult, _ := resource.cli.GetJSON(serviceQueryPath)

		var v2ServicesInstanceResponse ServicesInstanceResponse
		err := json.Unmarshal([]byte(jsonResult), &v2ServicesInstanceResponse)
		if err != nil {
			return nil, err
		}

		for _, ent := range v2ServicesInstanceResponse.Resources {
			serviceGUIDs = append(serviceGUIDs, ent.Metadata.GUID)
		}
	}

	return serviceGUIDs, nil
}
