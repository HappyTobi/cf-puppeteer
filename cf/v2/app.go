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

type AppResourcesEntity struct {
	Metadata Metadata `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

type MetaDataEntity struct {
	TotalResults       int                  `json:"total_results"`
	TotalPages         int                  `json:"total_pages"`
	PrevURL            interface{}          `json:"prev_url"`
	NextURL            interface{}          `json:"next_url"`
	AppResourcesEntity []AppResourcesEntity `json:"resources"`
}
type Metadata struct {
	GUID      string    `json:"guid"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Entity struct {
	Name                  string `json:"name"`
	Production            bool   `json:"production"`
	SpaceGUID             string `json:"space_guid"`
	StackGUID             string `json:"stack_guid"`
	Buildpack             string `json:"buildpack"`
	DetectedBuildpack     string `json:"detected_buildpack"`
	DetectedBuildpackGUID string `json:"detected_buildpack_guid"`
	EnvironmentJSON       struct {
		APPLICATIONNAME string `json:"APPLICATION_NAME"`
		OPTIMIZEMEMORY  string `json:"OPTIMIZE_MEMORY"`
	} `json:"environment_json"`
	Memory                   int         `json:"memory"`
	Instances                int         `json:"instances"`
	DiskQuota                int         `json:"disk_quota"`
	State                    string      `json:"state"`
	Version                  string      `json:"version"`
	Command                  string      `json:"command"`
	Console                  bool        `json:"console"`
	Debug                    interface{} `json:"debug"`
	StagingTaskID            string      `json:"staging_task_id"`
	PackageState             string      `json:"package_state"`
	HealthCheckType          string      `json:"health_check_type"`
	HealthCheckTimeout       int         `json:"health_check_timeout"`
	HealthCheckHTTPEndpoint  string      `json:"health_check_http_endpoint"`
	StagingFailedReason      interface{} `json:"staging_failed_reason"`
	StagingFailedDescription interface{} `json:"staging_failed_description"`
	Diego                    bool        `json:"diego"`
	DockerImage              interface{} `json:"docker_image"`
	DockerCredentials        struct {
		Username interface{} `json:"username"`
		Password interface{} `json:"password"`
	} `json:"docker_credentials"`
	PackageUpdatedAt     time.Time `json:"package_updated_at"`
	DetectedStartCommand string    `json:"detected_start_command"`
	EnableSSH            bool      `json:"enable_ssh"`
	Ports                []int     `json:"ports"`
	SpaceURL             string    `json:"space_url"`
	StackURL             string    `json:"stack_url"`
	RoutesURL            string    `json:"routes_url"`
	EventsURL            string    `json:"events_url"`
	ServiceBindingsURL   string    `json:"service_bindings_url"`
	RouteMappingsURL     string    `json:"route_mappings_url"`
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
		ui.Failed("no response / parsable response from %s %s", path, err)
		return nil, err
	}

	if len(metaDataResponseEntity.AppResourcesEntity) == 0 {
		return nil, ErrAppNotFound
	}

	return &metaDataResponseEntity.AppResourcesEntity[0], nil
}
