package v3

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ApplicationEntity struct {
	Command     string            `json:"command,omitempty"`
	HealthCheck HealthCheckEntity `json:"health_check"`
	ProcessType string            `json:"type,omitempty"`
}
type DataEntity struct {
	Endpoint          string `json:"endpoint,omitempty"`
	InvocationTimeout int    `json:"invocation_timeout,omitempty"`
	Timeout           int    `json:"timeout,omitempty"`
}
type HealthCheckEntity struct {
	Data            DataEntity `json:"data,omitempty"`
	HealthCheckType string     `json:"type"`
}

var (
	ErrInvocationTimeout = errors.New("could not set invocation timeout to application")
)

// UpdateApplicationProcessWebInformation calls v3 application to set options
// see api documentation http://v3-apidocs.cloudfoundry.org/version/3.67.0/index.html#update-an-app
func (resource *ResourcesData) UpdateApplicationProcessWebInformation(appGUID string, applicationEntity ApplicationEntity) error {
	path := fmt.Sprintf(`/v3/processes/%s`, appGUID)
	appJSON, err := json.Marshal(applicationEntity)
	if err != nil {
		return err
	}

	jsonResult, err := resource.Cli.PatchJSON(path, string(appJSON))

	if err != nil {
		return err
	}

	var applicationResponse ApplicationEntity
	err = json.Unmarshal([]byte(jsonResult), &applicationResponse)

	if err != nil {
		return err
	}

	if applicationResponse.HealthCheck.HealthCheckType != applicationEntity.HealthCheck.HealthCheckType {
		return ErrInvocationTimeout
	}
	return nil
}
