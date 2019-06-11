package cfResources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/happytobi/cf-puppeteer/ui"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/cf/appfiles"
	"code.cloudfoundry.org/cli/plugin"
)

//CfResourcesInterface interface for cfResource usage
type CfResourcesInterface interface {
	PushApp(appName string, spaceGUID string, buildpacks []string, stack string, envVars []string) (*V3AppResponse, error)
	CreatePackage(appGUID string) (*V3PackageResponse, error)
	UploadApplication(appName string, applicationFiles string, targetURL string) (*V3PackageResponse, error)
	CreateBuild(packageGUID string, buildpacks []string) (*V3BuildResponse, error)
	CheckBuildState(buildGUID string) (*V3BuildResponse, error)
	GetDropletGUID(buildGUID string) (*V3BuildResponse, error)
	AssignApp(appGUID string, dropletGUID string) error
	RouteMapping(appGUID string, routeGUID string) error
	StartApp(appGUID string) error
	//Stuff
	AssignAppManifest(appLink string, manifestPath string) error
	CheckPackageState(packageGUID string) (*V3PackageResponse, error)
	GetDomain(domains []map[string]string) (*[]V3Routes, error)
	GetApp(appGUID string) (*V3AppResponse, error)
	GetRoutesApp(appGUID string) ([]string, error)
	CreateServiceBinding(appGUID string, serviceInstanceGUID []string) error
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
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	/*EnvironmentVariables struct {
		Vars map[string]string `json:"var,omitempty"`
	} `json:"environment_variables,omitempty"`*/
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
func (resource *ResourcesData) PushApp(appName string, spaceGUID string, buildpacks []string, stack string, envVars []string) (*V3AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps`)

	var v3App V3Apps
	v3App.Name = appName
	v3App.Relationships.Space.Data.GUID = spaceGUID
	v3App.Lifecycle.LifecycleType = "buildpack"
	v3App.Lifecycle.LifecycleData.Stack = stack
	v3App.Lifecycle.LifecycleData.Buildpacks = buildpacks

	envs := make(map[string]string)
	for _, v := range envVars {
		envPair := strings.Split(v, "=")
		envKey := strings.TrimSpace(envPair[0])
		envVal := strings.TrimSpace(envPair[1])
		envs[envKey] = envVal
	}
	v3App.EnvironmentVariables = envs

	//TODO move to function
	appJSON, err := json.Marshal(v3App)
	if err != nil {
		return nil, err
	}

	if resource.TraceLogging {
		ui.Say("send POST to route: %s with body:", path)
		prettyPrintJSON(string(appJSON))
	}
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
		string(appJSON))

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		ui.Say("response from http call to path: %s was", path)
		prettyPrintJSON(jsonResp)
	}

	var response V3AppResponse
	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

//GetApp fetch all informations about an app with the appGUID
func (resource *ResourcesData) GetApp(appGUID string) (*V3AppResponse, error) {
	path := fmt.Sprintf(`/v3/apps/%s`, appGUID)

	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("response from get call GetApp to path: %s was: \n", path)
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

//V3DomainResponse reponse while loading domains
type V3DomainResponse struct {
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
		GUID          string    `json:"guid"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Name          string    `json:"name"`
		Internal      bool      `json:"internal"`
		Relationships struct {
			Organization struct {
				Data interface{} `json:"data"`
			} `json:"organization"`
			SharedOrganizations struct {
				Data []interface{} `json:"data"`
			} `json:"shared_organizations"`
		} `json:"relationships"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"resources"`
}

type V3Routes struct {
	Host       string
	DomainGUID string
}

//TODO optimize code
//ADD TCP option
//GetDomain create a package with v3 cf api
func (resource *ResourcesData) GetDomain(domains []map[string]string) (*[]V3Routes, error) {
	path := fmt.Sprintf(`/v3/domains`)

	if resource.TraceLogging {
		fmt.Printf("call api %s and search for domain name: %s", path, domains)
	}

	response, err := resource.getDomain(path)
	if err != nil {
		return nil, err
	}

	domainGUID := make(map[string]V3Routes)

	for _, domainRes := range response.Resources {
		for _, routes := range domains {
			for _, domain := range routes {
				_, exists := domainGUID[domain]
				if strings.Contains(domain, domainRes.Name) && !exists {
					hostName := strings.ReplaceAll(domain, domainRes.Name, "")
					hostName = strings.TrimRight(hostName, ".")
					newRoute := &V3Routes{
						Host:       hostName,
						DomainGUID: domainRes.GUID,
					}
					domainGUID[domain] = *newRoute
				}
			}
		}
	}

	for response.Pagination.Next.Href != "" && len(domainGUID) <= 0 {
		response, err := resource.getDomain(response.Pagination.Next.Href)
		if err != nil {
			return nil, err
		}

		for _, domainRes := range response.Resources {
			for _, routes := range domains {
				for _, domain := range routes {
					_, exists := domainGUID[domain]
					if strings.Contains(domain, domainRes.Name) && !exists {
						hostName := strings.ReplaceAll(domain, domainRes.Name, "")
						hostName = strings.TrimRight(hostName, ".")
						newRoute := &V3Routes{
							Host:       hostName,
							DomainGUID: domainRes.GUID,
						}
						domainGUID[domain] = *newRoute
					}
				}
			}
		}
	}

	var domainsFound []V3Routes
	for _, v := range domainGUID {
		domainsFound = append(domainsFound, v)
	}

	if resource.TraceLogging {
		fmt.Printf("domainGUID found return: %s \n", domainsFound)
	}

	return &domainsFound, err
}

func (resource *ResourcesData) getDomain(path string) (*V3DomainResponse, error) {
	var response V3DomainResponse
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")
	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("response from get call to path: %s was: \n", path)
		prettyPrintJSON(jsonResp)
	}

	err = json.Unmarshal([]byte(jsonResp), &response)
	if err != nil {
		return nil, err
	}

	sort.Slice(response.Resources, func(i, j int) bool {
		return response.Resources[i].Name < response.Resources[j].Name
	})
	return &response, nil
}

//UploadApplication upload a zip file to the created package
func (resource *ResourcesData) UploadApplication(appName string, applicationFiles string, targetURL string) (*V3PackageResponse, error) {
	if !resource.zipper.IsZipFile(applicationFiles) {

		zipFileName := fmt.Sprintf("%s%s.zip", os.TempDir(), appName)
		newZipFile, err := os.Create(zipFileName)

		if err != nil {
			return nil, err
		}
		defer newZipFile.Close()

		err = resource.zipper.Zip(applicationFiles, newZipFile)

		if resource.TraceLogging {
			fmt.Printf("zip will be created with from folder: %s - zip will be written as: %s \n", applicationFiles, zipFileName)
		}
		applicationFiles = zipFileName
	}

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
func (resource *ResourcesData) CreateBuild(packageGUID string, buildpacks []string) (*V3BuildResponse, error) {
	path := fmt.Sprintf(`/v3/builds`)
	var v3buildPackage V3BuildPackage
	v3buildPackage.Package.GUID = packageGUID
	v3buildPackage.Lifecycle.LifecycleType = "buildpack"

	for _, buildpack := range buildpacks {
		v3buildPackage.Lifecycle.LifecycleData.Buildpacks = append(v3buildPackage.Lifecycle.LifecycleData.Buildpacks, buildpack)
	}

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

type V3RouteMappingResponse struct {
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

//AssignApp to a created droplet guid
func (resource *ResourcesData) GetRoutesApp(appGUID string) ([]string, error) {
	path := fmt.Sprintf(`/v3/apps/%s/route_mappings`, appGUID)
	result, err := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")
	if err != nil {
		fmt.Printf("Error while calling the apply manifest url %s - error: %s \n", path, err)
		return nil, err
	}

	var response V3RouteMappingResponse
	jsonResp := strings.Join(result, "")

	if resource.TraceLogging {
		fmt.Printf("response from get call to path: %s was: \n", path)
		prettyPrintJSON(jsonResp)
	}

	err = json.Unmarshal([]byte(jsonResp), &response)
	var routeGUIDs []string
	for _, mapResource := range response.Resources {
		count := strings.LastIndex(mapResource.Links.Route.Href, "/")
		if count > 0 {
			routeGUIDs = append(routeGUIDs, mapResource.Links.Route.Href[count+1:])
		}
	}

	if resource.TraceLogging {
		fmt.Printf("return used routes from vendor app %s\n", routeGUIDs)
	}
	return routeGUIDs, err
}

//AssignAppManifest assign an appManifest
func (resource *ResourcesData) AssignAppManifest(appLink string, manifestPath string) error {
	path := fmt.Sprintf(`%s/actions/apply_manifest`, appLink)

	if resource.TraceLogging {
		fmt.Printf("Apply manifest to path %s and use manifestFilePath %s:\n", path, manifestPath)
	}

	file, err := os.Open(manifestPath)
	if err != nil {
		fmt.Printf("could not read manifest from path %s error: %s", manifestPath, err)
		panic(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Printf("manifest file stat error %s", err)
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
	if err != nil {
		fmt.Printf("Error while calling the apply manifest url %s - error: %s \n", path, err)
		return err
	}

	result, _ := ioutil.ReadAll(res.Body)
	jsonResp := string(result)

	if resource.TraceLogging {
		fmt.Printf("response while appling manifest to path %s - status %s\n", path, res.Status)
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

	/*
				210003
				{
		    "description": "The host is taken: proto-business-hydra-node-red-dev",
		    "error_code": "CF-RouteHostTaken",
		    "code": 210003
		}
	*/

	return err
}

type V3ServiceBinding struct {
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
		var v3ServiceBinding V3ServiceBinding
		v3ServiceBinding.Type = "app"
		v3ServiceBinding.Relationships.App.Data.GUID = appGUID
		v3ServiceBinding.Relationships.ServiceInstance.Data.GUID = serviceGUID
		appJSON, err := json.Marshal(v3ServiceBinding)
		if err != nil {
			return err
		}

		if resource.TraceLogging {
			fmt.Printf("post to : %s with appGUID: %s serviceGUID: %s\n", path, appGUID, serviceGUID)
			prettyPrintJSON(string(appJSON))
		}
		result, _ := resource.Connection.CliCommandWithoutTerminalOutput("curl", path, "-X", "POST", "-H", "Content-type: application/json", "-d",
			string(appJSON))
		jsonResp := strings.Join(result, "")

		if resource.TraceLogging {
			fmt.Printf("service binding created path: %s for service\n", path)
			prettyPrintJSON(jsonResp)
		}
	}

	return nil

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
