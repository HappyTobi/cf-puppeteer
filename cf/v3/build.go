package v3

import (
	"encoding/json"
	"fmt"
)

//BuildPackage represents post model of V3BuildPackage body
type BuildPackage struct {
	Package struct {
		GUID string `json:"guid"`
	} `json:"package"`
	Lifecycle struct {
		LifecycleType string `json:"type"`
		LifecycleData struct {
			Buildpacks []string `json:"buildpacks"`
			Stack      string   `json:"stack"`
		} `json:"data"`
	} `json:"lifecycle"`
}

//BuildResponse represents response ot the created build
type BuildResponse struct {
	GUID    string `json:"guid"`
	State   string `json:"state"`
	Droplet struct {
		GUID string `json:"guid"`
	} `json:"droplet"`
}

//CreateBuild with packagedGUID
func (resource *ResourcesData) CreateBuild(packageGUID string, buildpacks []string) (*BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds`)
	var buildPackage BuildPackage
	buildPackage.Package.GUID = packageGUID
	buildPackage.Lifecycle.LifecycleType = "buildpack"

	for _, buildpack := range buildpacks {
		buildPackage.Lifecycle.LifecycleData.Buildpacks = append(buildPackage.Lifecycle.LifecycleData.Buildpacks, buildpack)
	}

	buildPackage.Lifecycle.LifecycleData.Stack = "cflinuxfs3"

	appJSON, err := json.Marshal(buildPackage)
	if err != nil {
		return nil, err
	}

	jsonResult, err := resource.Cli.PostJSON(path, string(appJSON))
	if err != nil {
		return nil, err
	}

	var response BuildResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//CheckBuildState check the pushed application is staged or not
func (resource *ResourcesData) CheckBuildState(buildGUID string) (*BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds/%s`, buildGUID)
	jsonResult, err := resource.Cli.GetJSON(path)
	if err != nil {
		return nil, err
	}

	var response BuildResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//GetDropletGUID get dropletGUID for uploaded and staged build
func (resource *ResourcesData) GetDropletGUID(buildGUID string) (*BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds/%s`, buildGUID)
	jsonResult, err := resource.Cli.GetJSON(path)

	var response BuildResponse
	err = json.Unmarshal([]byte(jsonResult), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
