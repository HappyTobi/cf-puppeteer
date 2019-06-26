package v3

import (
	"encoding/json"
	"fmt"
)

type ServiceBinding struct {
	Type          string `json:"type"`
	Relationships struct {
		App struct {
			Data struct {
				GUID string `json:"guid"`
			} `json:"data"`
		} `json:"app"`
		ServiceInstance struct {
			Data struct {
				GUID string `json:"guid"`
			} `json:"data"`
		} `json:"service_instance"`
	} `json:"relationships"`
}

//CreateServiceBinding
func (resource *ResourcesData) CreateServiceBinding(appGUID string, serviceInstanceGUID []string) error {
	path := fmt.Sprintf(`/v3/service_bindings`)

	for _, serviceGUID := range serviceInstanceGUID {
		var v3ServiceBinding ServiceBinding
		v3ServiceBinding.Type = "app"
		v3ServiceBinding.Relationships.App.Data.GUID = appGUID
		v3ServiceBinding.Relationships.ServiceInstance.Data.GUID = serviceGUID
		appJSON, err := json.Marshal(v3ServiceBinding)
		if err != nil {
			return err
		}

		_, _ = resource.Cli.PostJSON(path, string(appJSON))
	}
	return nil
}
