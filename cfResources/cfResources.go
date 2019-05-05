package cfResources

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/cli/cf/appfiles"
	"code.cloudfoundry.org/cli/plugin"
)

//https://github.com/cloudfoundry/cloud_controller_ng/wiki/How-to-Create-an-App-Using-V3-of-the-CC-API

//CfResourcesInterface interface for cfResource usage
type CfResourcesInterface interface {
	//Add methods here
	//Step1
	PushApp(appName string, spaceGUID string) (*V3AppResponse, error)
	//Step 3
	CreatePackage(appGUID string) (*V3PackageResponse, error)
	//Step 4 & 5
	UploadApplication(appName string, applicationFiles string, targetURL string) error
	//Step 6
	CreateBuild(packageGUID string) (*V3BuildResponse, error)
	//Step 7
	CheckBuildState(buildGUID string) (*V3BuildResponse, error)
	//Step 8
	GetDropletGUID(buildGUID string) (*V3BuildResponse, error)
	//Step 9
	AssignApp(appGUID string, dropletGUID string) error
	//Step 10
	//CreateRoute do with v2 command
	//Step 11
	RouteMapping(appGUID string, routeGUID string) error
	//Step 12
	StartApp(appGUID string) error
}

//ResourcesData struct to hold important instances to run push
type ResourcesData struct {
	Connection   plugin.CliConnection
	TraceLogging bool
	zipper       appfiles.Zipper
}

//NewResources create a new instance of CFResources
func NewResources(conn plugin.CliConnection, traceLogging bool) *ResourcesData {
	return &ResourcesData{
		Connection:   conn,
		TraceLogging: traceLogging,
		zipper:       &appfiles.ApplicationZipper{},
	}
}

//V3Apps application struct
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
	} `json:"environment_variables,omitempty"`
}

//V3AppResponse application response struct
type V3AppResponse struct {
	GUID string `json:"guid"`
}

//PushApp push app with v3 api to cloudfoundry
func (resource *ResourcesData) PushApp(appName string, spaceGUID string) (*V3AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps`)

	var v3App V3Apps
	v3App.Name = appName
	v3App.Relationships.Space.Data.GUID = spaceGUID

	//TODO move to function
	appJSON, err := json.Marshal(v3App)
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
		fmt.Printf("response from http call to path: %s was:\n", path)
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

//StartApp start app on cf
func (resource *ResourcesData) StartApp(appGUID string) error {
	path := fmt.Sprintf(`/v3/apps/%s/actions/start`, appGUID)
	_, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json")
	return err
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
	GUID  string `json:"guid"`
	Links struct {
		Upload struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"upload"`
	} `json:"links"`
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
		fmt.Printf("response from http call to path: %s was:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V3PackageResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//UploadApplication upload a zip file to the created package
func (resource *ResourcesData) UploadApplication(appName string, applicationFiles string, targetURL string) error {
	//TODO
	zipFile := applicationFiles
	if !resource.zipper.IsZipFile(applicationFiles) {
		zipFileName := fmt.Sprintf("%s%s.zip", os.TempDir(), appName)
		newZipFile, err := os.Create(zipFileName)

		if err != nil {
			return err
		}
		defer newZipFile.Close()

		err = resource.zipper.Zip(applicationFiles, newZipFile)
		if err != nil {
			return err
		}
		zipFile = zipFileName
	}

	/*zipFile, err := resource.zipUploadFile(appName, applicationFiles)
	if err != nil {
		return err
	}*/

	file, err := os.Open(zipFile)
	defer file.Close()

	token, _ := resource.Connection.AccessToken()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("bits", filepath.Base(zipFile))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	defer writer.Close()

	req, err := http.NewRequest("POST", targetURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("post error: %s\n", err)
		panic(err)
	}

	defer resp.Body.Close()
	message, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(message))
	//defer os.Remove(zipFile.Name())
	return err
}

//CheckBuildState check the pushed application is staged or not
func (resource *ResourcesData) CheckBuildState(buildGUID string) (*V3BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds/%s`, buildGUID)

	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json")

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("Cloud Foundry API response while checking build state %s:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V3BuildResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

//GetDropletGUID get dropletGUID for uploaded and staged build
func (resource *ResourcesData) GetDropletGUID(buildGUID string) (*V3BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds/%s`, buildGUID)

	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("Cloud Foundry API response while getting build information (droplet) %s:\n", path)
		prettyPrintJSON(jsonResp)
	}

	var response V3BuildResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

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
	GUID    string `json:"guid"`
	State   string `json:"state"`
	Droplet struct {
		GUID string `json:"guid"`
	} `json:"droplet"`
}

//CreateBuild with packagedGUID
func (resource *ResourcesData) CreateBuild(packageGUID string) (*V3BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds`)
	var v3buildPackage V3BuildPackage
	v3buildPackage.Package.GUID = packageGUID
	v3buildPackage.Lifecycle.LifecycleType = "buildpack"
	v3buildPackage.Lifecycle.LifecycleData.Buildpacks = append(v3buildPackage.Lifecycle.LifecycleData.Buildpacks, "")

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

//V3AppsDroplet -> maybe use V3Apps only
type V3AppsDroplet struct {
	Data struct {
		GUID string `json:"guid"`
	} `json:"data"`
}

//AssignApp to a created droplet guid
func (resource *ResourcesData) AssignApp(appGUID string, dropletGUID string) error {
	path := fmt.Sprintf(`/v3/apps/%s/relationships/current_droplet`, appGUID)

	var v3AppsDroplet V3AppsDroplet
	v3AppsDroplet.Data.GUID = dropletGUID

	appJSON, err := json.Marshal(v3AppsDroplet)
	if err != nil {
		return err
	}

	_, err = resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "PATCH", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	return err
}

type V3RouteMapping struct {
	Relationships struct {
		App struct {
			GUID string `json:"guid"`
		} `json:"app"`
		Route struct {
			GUID string `json:"guid"`
		} `json:"route"`
	} `json:"relationships"`
}

//RouteMapping map route to application
func (resource *ResourcesData) RouteMapping(appGUID string, routeGUID string) error {
	path := fmt.Sprintf(`v3/route_mappings`)

	var v3RouteMapping V3RouteMapping
	v3RouteMapping.Relationships.App.GUID = appGUID
	v3RouteMapping.Relationships.Route.GUID = routeGUID

	appJSON, err := json.Marshal(v3RouteMapping)
	if err != nil {
		return err
	}

	_, err = resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	return err
}

//TODO optimize an fix zip issue
// drop own implementation switch to appfiles/zipper
// see push.go file from cf  - GatherFiles
//zipUploadFile upload the application files
func (resource *ResourcesData) zipUploadFile(appName string, fileName string) (string, error) {
	zipFileName := fmt.Sprintf("%s%s.zip", os.TempDir(), appName)
	if resource.TraceLogging {
		fmt.Printf("try to create zip file: %s from passed file / folder %s \n", zipFileName, filepath.Base(fileName))
	}

	newZipFile, err := os.Create(zipFileName)

	if err != nil {
		return "", err
	}

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	fileToZip, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return "", err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return "", err
	}

	header.Name = fileName
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(writer, fileToZip)
	fmt.Printf("return zip file: %s\n", zipFileName)
	return zipFileName, err
}
