package v3

import (
	"code.cloudfoundry.org/cli/cf/appfiles"
	"code.cloudfoundry.org/cli/plugin"
	"errors"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"
	"github.com/happytobi/cf-puppeteer/ui"
	"time"
)

//Push interface with all v3 actions
type Push interface {
	PushApplication(appName string, venAppName string, appPath string, serviceNames []string, spaceGuid string, buildpacks []string, applicationStack string, environmentVariables []string, manifestPath string, routes []map[string]string, v2Resources v2.Resources, healthCheckType string, healthCheckHttpEndpoint string, process string, invocationTimeout int) error
}

//ResourcesData internal struct with connection an tracing options etc
type ResourcesData struct {
	zipper     appfiles.Zipper
	Cli        cli.Calls
	httpClient cli.HttpCalls
	Connection plugin.CliConnection
}

//NewV3Push constructor
func NewV3Push(conn plugin.CliConnection, traceLogging bool) *ResourcesData {
	return &ResourcesData{
		zipper:     &appfiles.ApplicationZipper{},
		Cli:        cli.NewCli(conn, traceLogging),
		httpClient: cli.NewHttpClient(conn, traceLogging, 30, false),
		Connection: conn,
	}
}

var (
	ErrAppNotFound = errors.New("application not found")
)

//PushApplication call all methods to push a complete application
func (resource *ResourcesData) PushApplication(appName string, venAppName string, appPath string, serviceNames []string, spaceGuid string, buildpacks []string, applicationStack string, environmentVariables []string, manifestPath string, routes []map[string]string, v2Resources v2.Resources, healthCheckType string, healthCheckHttpEndpoint string, process string, invocationTimeout int) error {
	appResponse, err := resource.PushApp(appName, spaceGuid, buildpacks, applicationStack, environmentVariables)
	if err != nil {
		return err
	}

	err = resource.AssignAppManifest(appResponse.Links.Self.Href, manifestPath)
	if err != nil {
		return err
	}

	createPackageResponse, err := resource.CreatePackage(appResponse.GUID)
	if err != nil {
		return err
	}

	domains, err := resource.GetDomain(routes)
	if err != nil {
		return err
	}

	for _, route := range *domains {
		routeResponse, err := v2Resources.CreateRoute(spaceGuid, route.DomainGUID, route.Host)
		if err != nil {
			return err
		}
		err = resource.RouteMapping(appResponse.GUID, routeResponse.Metadata.GUID)
		if err != nil {
			return err
		}

		ui.Say("add routes to application")
	}

	ui.Ok()

	venApp, err := v2Resources.GetAppMetadata(venAppName)
	if err != v2.ErrAppNotFound && err != nil {
		ui.Failed("metadata error %s", err)
		return err
	}

	var venRoutes []string

	if venApp != nil {
		venRoutes, err = resource.GetRoutesApp(venApp.Metadata.GUID)
		if err != nil {
			return err
		}
	}

	ui.Say("map vendor routes to new application")
	for _, route := range venRoutes {
		err = resource.RouteMapping(appResponse.GUID, route)
		ui.LoadingIndication()
		if err != nil {
			return err
		}
	}
	ui.Say("")
	ui.Ok()

	//map services
	serviceGUIDs, err := v2Resources.FindServiceInstances(serviceNames, spaceGuid)
	if err != nil {
		return err
	}

	err = resource.CreateServiceBinding(appResponse.GUID, serviceGUIDs)
	if err != nil {
		return err
	}

	createPackageResponse, err = resource.UploadApplication(appName, appPath, createPackageResponse.Links.Upload.Href)
	if err != nil {
		return err
	}

	duration, _ := time.ParseDuration("1s")

	ui.Say("start uploading application")
	for createPackageResponse.State != "FAILED" &&
		createPackageResponse.State != "READY" &&
		createPackageResponse.State != "EXPIRED" {
		ui.LoadingIndication()
		time.Sleep(duration)
		createPackageResponse, err = resource.CheckPackageState(createPackageResponse.GUID)
		if err != nil {
			return nil
		}
	}
	ui.Say("")
	ui.Ok()

	ui.Say("create build")
	buildResponse, err := resource.CreateBuild(createPackageResponse.GUID, buildpacks)
	if err != nil {
		return err
	}
	ui.Ok()

	//set timeouts
	ui.Say("set timeout parameters")
	err = resource.SetHealthCheck(healthCheckType, healthCheckHttpEndpoint, invocationTimeout, process, appResponse.GUID)
	if err != nil {
		return err
	}
	ui.Ok()

	ui.Say("stage application")
	for buildResponse.State != "FAILED" &&
		buildResponse.State != "STAGED" {
		time.Sleep(duration)
		buildResponse, err = resource.CheckBuildState(buildResponse.GUID)
		ui.LoadingIndication()
		if err != nil {
			return nil
		}
	}
	ui.Say("")
	ui.Ok()

	dropletResponse, err := resource.GetDropletGUID(buildResponse.GUID)
	err = resource.AssignApp(appResponse.GUID, dropletResponse.GUID)
	return nil
}

// SetHealthCheckV3 sets the health check for the specified application using the given health check configuration
func (resource *ResourcesData) SetHealthCheck(healthCheckType string, healthCheckHTTPEndpoint string, invocationTimeout int, process string, GUID string) error {
	// Without a health check type, the CF command is not valid. Therefore, leave if the type is not specified
	if healthCheckType == "" {
		return nil
	}

	// load application by guid
	appProcessEntity, err := resource.GetApplicationProcessWebInformation(GUID)
	if err != nil {
		return err
	}

	applicationEntity := ApplicationEntity{}
	applicationEntity.Command = appProcessEntity.Command
	applicationEntity.HealthCheck.HealthCheckType = healthCheckType
	applicationEntity.HealthCheck.Data.Timeout = 120 //use parsed argument

	if healthCheckType == "http" && healthCheckHTTPEndpoint != "" {
		applicationEntity.HealthCheck.Data.Endpoint = healthCheckHTTPEndpoint
		if invocationTimeout >= 0 {
			applicationEntity.HealthCheck.Data.InvocationTimeout = invocationTimeout
		}
	} else if process != "" && (healthCheckType == "process" || healthCheckType == "port") {
		applicationEntity.ProcessType = process
	}

	err = resource.UpdateApplicationProcessWebInformation(appProcessEntity.GUID, applicationEntity)
	return err
}
