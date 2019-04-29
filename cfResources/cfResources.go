package cfResources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

//https://github.com/cloudfoundry/cloud_controller_ng/wiki/How-to-Create-an-App-Using-V3-of-the-CC-API

type CfResourcesInterface interface {
	//Add methods here
	PushApp(appName string, spaceGUID string) (*V3AppResponse, error)
	CreatePackage(appGUID string) (*V3PackageResponse, error)

	CreateBuild(packageGUID string) (*V3BuildResponse, error)
}

type ResourcesData struct {
	Connection   plugin.CliConnection
	TraceLogging bool
}

func NewResources(conn plugin.CliConnection, traceLogging bool) *ResourcesData {
	return &ResourcesData{
		Connection:   conn,
		TraceLogging: traceLogging,
	}
}

type V3Apps struct {
	Name          string `json:"name"`
	Relationships struct {
		Space struct {
			Data struct {
				GUID string `json:"guid"`
			} `json:"data"`
		} `json:"space"`
	} `json:"relationships"`
	EnvironmentVariables struct {
		Vars map[string]string `json:"var"`
	} `json:"environmentVariables,omitempty"`
}

type V3AppResponse struct {
	GUID string `json:"guid"`
}

//PushApp push app with v3 api to cloudfoundry
func (resource *ResourcesData) PushApp(appName string, spaceGUID string) (*V3AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps`)

	var v3App V3Apps
	v3App.Relationships.Space.Data.GUID = spaceGUID

	//TODO move to function
	appJSON, err := json.Marshal(v3App)
	if err != nil {
		return nil, err
	}
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("Cloud Foundry API response to PATCH call on %s:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V3AppResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}

	//TODO add error
	/*if len(applicationProcessResponse.GUID) == 0 {
		return nil, ErrAppNotFound
	}*/

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

//V3Package represents post model of V3Package body
type V3Package struct {
	PackageType   string `json:"type"`
	Relationships struct {
		App struct {
			Data struct {
				GUID string `json:"guid"`
			} `json:"data"`
		} `json:"app"`
	} `json:"relationships"`
}

//V3PackageResponse create package response payload
type V3PackageResponse struct {
	GUID string `json:"guid"`
}

//CreatePackage create a package with v3 cf api
func (resource *ResourcesData) CreatePackage(appGUID string) (*V3PackageResponse, error) {
	path := fmt.Sprintf(`/v3/packages`)
	var v3Package V3Package
	v3Package.PackageType = "bits"
	v3Package.Relationships.App.Data.GUID = appGUID

	//TODO move to function
	appJSON, err := json.Marshal(v3Package)
	if err != nil {
		return nil, err
	}
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("Cloud Foundry API response to PATCH call on %s:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V3PackageResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

/*
func (resource *ResourcesData) uploadApplication(packageGUID string) error {
	path := fmt.Sprintf(`/v3/packages/%s/upload`, packageGUID)

	//TODO BUILD CUSTOM CURL COMMAND
	_, err := resource.conn.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json")

	return err
}

type V3Build struct {
	Package struct {
		guid string `json:"guid"`
	} `json:"package"`
}

func (resource *ResourcesData) stagePackage(packageGUID string) error {
	path := fmt.Sprintf(`/v3/builds`)
	appData := &V3Build{}

	//TODO check alt commands
	appJSON, err := json.Marshal(appData)
	if err != nil {
		return err
	}
	_, err := resource.conn.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d", string(appJSON))

	return err
}*/

//V3BuildPackage represents post model of V3BuildPackage body
type V3BuildPackage struct {
	Package struct {
		GUID string `json:"guid"`
	} `json:"package"`
	Lifecycle struct {
		LifecycleType string `json:"type"`
		LifecycleData struct {
			Buildpacks []string `json:"buildpacks"`
		} `json:"data"`
	} `json:"lifecycle"`
}

//V3BuildResponse represents response ot the created build
type V3BuildResponse struct {
	GUID string `json:"guid"`
}

//check till build is staged
func (resource *ResourcesData) CreateBuild(packageGUID string) (*V3BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds`)
	var v3buildPackage V3BuildPackage
	v3buildPackage.Package.GUID = packageGUID
	v3buildPackage.Lifecycle.LifecycleType = "buildpack"
	v3buildPackage.Lifecycle.LifecycleData.Buildpacks[0] = "" //TODO

	//TODO move to function
	appJSON, err := json.Marshal(v3buildPackage)
	if err != nil {
		return nil, err
	}
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("Cloud Foundry API response to PATCH call on %s:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V3BuildResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//TODO step 8 - 13
