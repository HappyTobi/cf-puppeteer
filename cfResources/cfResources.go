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
				Guid string `json:"guid"`
			} `json:"data"`
		} `json:"space"`
	} `json:"relationships"`
	EnvironmentVariables struct {
		Vars map[string]string `json:"var"`
	} `json:"environmentVariables,omitempty"`
}

type V3AppResponse struct {
	Guid string `json:"guid"`
}

func (resource *ResourcesData) PushApp(appName string, spaceGUID string) (*V3AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps`)
	var v3App V3Apps
	v3App.Relationships.Space.Data.Guid = spaceGUID
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

/*
type V3Package struct {
	packageType   string `json:"type"`
	Relationships struct {
		App struct {
			Data struct {
				guid string `json:"guid"`
			} `json:"data"`
		} `json:"app"`
	} `json:"relationships"`
}

func (resource *ResourcesData) createPackage(appGUID string) error {
	path := fmt.Sprintf(`/v3/packages`)
	appData := &V3Package{}

	//TODO move to function
	appJSON, err := json.Marshal(appData)
	if err != nil {
		return err
	}
	result, err := resource.conn.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	return err
}

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
}

//check till build is staged
func (resource *ResourcesData) checkBuildStage(buildGUID string) error {
	path := fmt.Sprintf(`/v3/builds/%s`, buildGUID)

	_, err := resource.conn.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")
	//check state "STAGED"

	return err
}

//TODO step 8 - 13
*/
