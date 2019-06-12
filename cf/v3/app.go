package v3

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

//Apps application struct
type Apps struct {
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
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
}

//AppResponse application response struct
type AppResponse struct {
	GUID  string `json:"guid"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Space struct {
			Href string `json:"href"`
		} `json:"space"`
		Processes struct {
			Href string `json:"href"`
		} `json:"processes"`
		RouteMappings struct {
			Href string `json:"href"`
		} `json:"route_mappings"`
		Packages struct {
			Href string `json:"href"`
		} `json:"packages"`
		EnvironmentVariables struct {
			Href string `json:"href"`
		} `json:"environment_variables"`
		CurrentDroplet struct {
			Href string `json:"href"`
		} `json:"current_droplet"`
		Droplets struct {
			Href string `json:"href"`
		} `json:"droplets"`
		Tasks struct {
			Href string `json:"href"`
		} `json:"tasks"`
		Start struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"start"`
		Stop struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"stop"`
		Revisions struct {
			Href string `json:"href"`
		} `json:"revisions"`
		DeployedRevisions struct {
			Href string `json:"href"`
		} `json:"deployed_revisions"`
	} `json:"links"`
	Metadata struct {
		Labels struct {
		} `json:"labels"`
		Annotations struct {
		} `json:"annotations"`
	} `json:"metadata"`
}

//AppsDroplet -> maybe use V3Apps only
type AppsDroplet struct {
	Data struct {
		GUID string `json:"guid"`
	} `json:"data"`
}

//RouteMappingResponse
type RouteMappingResponse struct {
	Pagination struct {
		TotalResults int `json:"total_results"`
		TotalPages   int `json:"total_pages"`
		First        struct {
			Href string `json:"href"`
		} `json:"first"`
		Last struct {
			Href string `json:"href"`
		} `json:"last"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		Previous interface{} `json:"previous"`
	} `json:"pagination"`
	Resources []struct {
		GUID      string    `json:"guid"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Weight    int       `json:"weight"`
		Links     struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			App struct {
				Href string `json:"href"`
			} `json:"app"`
			Route struct {
				Href string `json:"href"`
			} `json:"route"`
			Process struct {
				Href string `json:"href"`
			} `json:"process"`
		} `json:"links"`
	} `json:"resources"`
}

//PushApp push app with v3 api to cloudfoundry
func (resource *ResourcesData) PushApp(appName string, spaceGUID string, buildpacks []string, stack string, envVars []string) (*AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps`)

	var app Apps
	app.Name = appName
	app.Relationships.Space.Data.GUID = spaceGUID
	app.Lifecycle.LifecycleType = "buildpack"
	app.Lifecycle.LifecycleData.Stack = stack
	app.Lifecycle.LifecycleData.Buildpacks = buildpacks

	envs := make(map[string]string)
	for _, v := range envVars {
		envPair := strings.Split(v, "=")
		envKey := strings.TrimSpace(envPair[0])
		envVal := strings.TrimSpace(envPair[1])
		envs[envKey] = envVal
	}
	app.EnvironmentVariables = envs

	//TODO move to function
	appJSON, err := json.Marshal(app)
	if err != nil {
		return nil, err
	}

	jsonResult, err := resource.Cli.PostJSON(path, string(appJSON))
	if err != nil {
		return nil, err
	}

	var response AppResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

//GetApp fetch all information about an app with the appGUID
func (resource *ResourcesData) GetApp(appGUID string) (*AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps/%s`, appGUID)

	jsonResult, err := resource.Cli.GetJSON(path)
	if err != nil {
		return nil, err
	}

	var response AppResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//StartApp start app on cf
func (resource *ResourcesData) StartApp(appGUID string) error {
	path := fmt.Sprintf(`/v3/apps/%s/actions/start`, appGUID)
	_, err := resource.Cli.PostJSON(path, "")
	return err
}

//AssignApp to a created droplet guid
func (resource *ResourcesData) AssignApp(appGUID string, dropletGUID string) error {
	path := fmt.Sprintf(`/v3/apps/%s/relationships/current_droplet`, appGUID)

	var appsDroplet AppsDroplet
	appsDroplet.Data.GUID = dropletGUID

	appJSON, err := json.Marshal(appsDroplet)
	if err != nil {
		return err
	}

	_, err = resource.Cli.PatchJSON(path, string(appJSON))

	return err
}

//AssignApp to a created droplet guid
func (resource *ResourcesData) GetRoutesApp(appGUID string) ([]string, error) {
	path := fmt.Sprintf(`/v3/apps/%s/route_mappings`, appGUID)

	jsonResult, err := resource.Cli.GetJSON(path)

	var response RouteMappingResponse

	err = json.Unmarshal([]byte(jsonResult), &response)
	var routeGUIDs []string
	for _, mapResource := range response.Resources {
		count := strings.LastIndex(mapResource.Links.Route.Href, "/")
		if count > 0 {
			routeGUIDs = append(routeGUIDs, mapResource.Links.Route.Href[count+1:])
		}
	}

	/*if resource.TraceLogging {
		fmt.Printf("return used routes from vendor app %s\n", routeGUIDs)
	}*/

	return routeGUIDs, err
}
