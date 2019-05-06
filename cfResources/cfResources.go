package cfResources

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/cf/appfiles"
	"code.cloudfoundry.org/cli/plugin"
)

//CfResourcesInterface interface for cfResource usage
type CfResourcesInterface interface {
	//Step1
	PushApp(appName string, spaceGUID string, buildpacks []string, stack string) (*V3AppResponse, error)
	//Step 3
	CreatePackage(appGUID string) (*V3PackageResponse, error)
	//Step 4 & 5
	UploadApplication(appName string, applicationFiles string, targetURL string) (*V3PackageResponse, error)
	//Step 6
	CreateBuild(packageGUID string) (*V3BuildResponse, error)
	//Step 7
	CheckBuildState(buildGUID string) (*V3BuildResponse, error)
	//Step 8
	GetDropletGUID(buildGUID string) (*V3BuildResponse, error)
	//Step 9
	AssignApp(appGUID string, dropletGUID string) error
	//Step 10
	//CreateRoute?
	//Step 11
	RouteMapping(appGUID string, routeGUID string) error
	//Step 12
	StartApp(appGUID string) error
	//Stuff
	AssignAppManifest(appLink string, manifestPath string) error
	CheckPackageState(packageGUID string) (*V3PackageResponse, error)
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

//V3AppResponse application response struct
type V3AppResponse struct {
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

//PushApp push app with v3 api to cloudfoundry
func (resource *ResourcesData) PushApp(appName string, spaceGUID string, buildpacks []string, stack string) (*V3AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps`)

	var v3App V3Apps
	v3App.Name = appName
	v3App.Relationships.Space.Data.GUID = spaceGUID
	v3App.Lifecycle.LifecycleType = "buildpack"
	v3App.Lifecycle.LifecycleData.Stack = stack
	v3App.Lifecycle.LifecycleData.Buildpacks = buildpacks

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

	return &response, nil
}

//StartApp start app on cf
func (resource *ResourcesData) StartApp(appGUID string) error {
	path := fmt.Sprintf(`/v3/apps/%s/actions/start`, appGUID)
	_, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json")
	return err
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
	State string `json:"state"`
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

//CheckPackageState create a package with v3 cf api
func (resource *ResourcesData) CheckPackageState(packageGUID string) (*V3PackageResponse, error) {
	path := fmt.Sprintf(`/v3/packages/%s`, packageGUID)

	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("response from get call to path: %s was: \n", path)
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
func (resource *ResourcesData) UploadApplication(appName string, applicationFiles string, targetURL string) (*V3PackageResponse, error) {
	/*if !resource.zipper.IsZipFile(applicationFiles) {
	zipFileName := fmt.Sprintf("%s%s.zip", os.TempDir(), appName)
	newZipFile, err := os.Create(zipFileName)
	*/
	file, err := os.Open(applicationFiles)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("bits", filepath.Base(applicationFiles))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()

	request, err := http.NewRequest(
		"POST",
		targetURL,
		bytes.NewReader(body.Bytes()),
	)

	token, _ := resource.Connection.AccessToken()
	request.Header = http.Header{}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "client.userAgent")
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	timeoutDuration, _ := time.ParseDuration("5m")

	tr := &http.Transport{
		/*TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //config.SkipSSLValidation,
		},*/
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			Timeout:   timeoutDuration,
		}).DialContext,
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("post error: %s\n", err)
		panic(err)
	}

	//defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	jsonResp := string(result)

	if resource.TraceLogging {
		fmt.Printf("Cloud Foundry API response while uploading the artifact %s:\n", targetURL)
		prettyPrintJSON(jsonResp)
	}

	var response V3PackageResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}

	return &response, err

}

//CheckBuildState check the pushed application is staged or not
func (resource *ResourcesData) CheckBuildState(buildGUID string) (*V3BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds/%s`, buildGUID)

	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")

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
			Stack      string   `json:"stack"`
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
	v3buildPackage.Lifecycle.LifecycleData.Buildpacks = append(v3buildPackage.Lifecycle.LifecycleData.Buildpacks, "java_buildpack")
	v3buildPackage.Lifecycle.LifecycleData.Stack = "cflinuxfs3"

	appJSON, err := json.Marshal(v3buildPackage)
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
		fmt.Printf("Cloud Foundry API response to POST call on %s:\n", path)
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

//AssignAppManifest assign an appManifest
func (resource *ResourcesData) AssignAppManifest(appLink string, manifestPath string) error {
	path := fmt.Sprintf(`%s/actions/apply_manifest`, appLink)

	file, err := os.Open(manifestPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)
	file.Read(buffer)

	request, err := http.NewRequest(
		"POST",
		path,
		bytes.NewReader(buffer),
	)

	token, _ := resource.Connection.AccessToken()

	request.Header = http.Header{}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "client.userAgent")
	request.Header.Set("Content-Type", "application/x-yaml")
	request.Header.Set("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(request)

	result, _ := ioutil.ReadAll(res.Body)
	jsonResp := string(result)

	if resource.TraceLogging {
		fmt.Printf("Cloud Foundry API response while appling manifest %s:\n", path)
		prettyPrintJSON(jsonResp)
	}

	defer res.Body.Close()

	return nil
}

//REMOVE?
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

//RouteMapping map route to application REMOVE?
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
