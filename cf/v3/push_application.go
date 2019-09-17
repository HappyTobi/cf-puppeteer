package v3

import (
	"code.cloudfoundry.org/cli/cf/appfiles"
	"code.cloudfoundry.org/cli/plugin"

	"errors"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"
	"github.com/happytobi/cf-puppeteer/manifest"
	"github.com/happytobi/cf-puppeteer/ui"
	"github.com/jinzhu/copier"
)

//Push interface with all v3 actions
type Push interface {
	PushApplication(venAppName string, spaceGUID string, parsedArguments *arguments.ParserArguments, v2Resources v2.Resources) error
	SwitchRoutesOnly(venAppName string, appName string, spaceGUID string, routes []map[string]string, v2Resources v2.Resources) error
}

//ResourcesData internal struct with connection an tracing options etc
type ResourcesData struct {
	zipper     appfiles.Zipper
	Cli        cli.Calls
	httpClient cli.HttpCalls
	Connection plugin.CliConnection
	Executor   cli.Executor
}

//NewV3Push constructor
func NewV3Push(conn plugin.CliConnection, traceLogging bool) *ResourcesData {
	return &ResourcesData{
		zipper:     &appfiles.ApplicationZipper{},
		Cli:        cli.NewCli(conn, traceLogging),
		httpClient: cli.NewHttpClient(conn, traceLogging, 30, true),
		Connection: conn,
		Executor:   cli.NewExecutor(traceLogging),
	}

}

var (
	ErrAppNotFound = errors.New("application not found")
)

//PushApplication call all methods to push a complete application
func (resource *ResourcesData) PushApplication(venAppName, spaceGUID string, parsedArguments *arguments.ParserArguments, v2Resources v2.Resources) error {

	/*buildpacks := parsedArguments.Manifest.ApplicationManifests[0].Buildpacks
	applicationStack := parsedArguments.Manifest.ApplicationManifests[0].Stack
	appName := parsedArguments.AppName
	appPath := parsedArguments.AppPath
	serviceNames := parsedArguments.Manifest.ApplicationManifests[0].Services
	manifestPath := parsedArguments.ManifestPath
	routes := parsedArguments.Manifest.ApplicationManifests[0].Routes
	healthCheckType := parsedArguments.HealthCheckType
	healthCheckHttpEndpoint := parsedArguments.HealthCheckHTTPEndpoint
	process := parsedArguments.Process
	invocationTimeout := parsedArguments.InvocationTimeout
	timeout := parsedArguments.Timeout*/

	err := resource.CreateApp(parsedArguments)
	if err != nil {
		return err
	}
	ui.Ok()

	err = resource.AssignAppManifest(parsedArguments)
	if err != nil {
		return err
	}
	ui.Ok()

	err = resource.PushApp(parsedArguments)
	if err != nil {
		return err
	}
	ui.Ok()

	err = resource.SetHealthCheck(parsedArguments.AppName, parsedArguments.HealthCheckType, parsedArguments.HealthCheckHTTPEndpoint, parsedArguments.InvocationTimeout, parsedArguments.Process)
	if err != nil {
		return err
	}
	ui.Ok()

	return nil
}

//SwitchRoutes switch route interface method to provide switch routes only option
func (resource *ResourcesData) SwitchRoutesOnly(venAppName string, appName string, spaceGUID string, routes []map[string]string, v2Resources v2.Resources) (err error) {
	appResource, err := v2Resources.GetAppMetadata(appName)
	if err != nil {
		return err
	}

	return resource.SwitchRoutes(venAppName, appResource.Metadata.GUID, routes, spaceGUID, v2Resources)
}

//GenerateNoRouteYml generate temp manifest without routes to skip route creation
func (resource *ResourcesData) GenerateNoRouteYml(appName string, originalManifest manifest.Manifest) (newManifestPath string, err error) {
	manifestPathTemp := resource.GenerateTempFile(appName, "yml")
	//Clone manifest to change them without side effects
	newTempManifest := manifest.Manifest{}
	err = copier.Copy(&newTempManifest, &originalManifest)
	if err != nil {
		return "", err
	}
	//clean up manifest
	newTempManifest.ApplicationManifests[0].NoRoute = true
	newTempManifest.ApplicationManifests[0].Routes = []map[string]string{}

	_, err = manifest.WriteYmlFile(manifestPathTemp, originalManifest)
	if err != nil {
		return "", err
	}
	return manifestPathTemp, nil
}

//SwitchRoutes add new routes and switch "old" one from vendor app to the one
func (resource *ResourcesData) SwitchRoutes(venAppName string, pushedAppGUID string, routes []map[string]string, spaceGUID string, v2Resources v2.Resources) (err error) {
	domains, err := resource.GetDomain(routes)
	if err != nil {
		return err
	}

	for _, route := range *domains {
		routeResponse, err := v2Resources.CreateRoute(spaceGUID, route.DomainGUID, route.Host)
		if err != nil {
			return err
		}
		err = resource.RouteMapping(pushedAppGUID, routeResponse.Metadata.GUID)
		if err != nil {
			return err
		}
		ui.Say("route with host %s added", route.Host)
	}
	ui.Say("routes from manifest added to application")
	ui.Ok()

	ui.Say("switch routes - defined in manifest and uses by vendor app to app")
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

	ui.Say("map all routes to new application")
	for _, route := range venRoutes {
		err = resource.RouteMapping(pushedAppGUID, route)
		ui.LoadingIndication()
		if err != nil {
			return err
		}
	}
	ui.Say("")
	ui.Ok()
	return nil
}

// SetHealthCheckV3 sets the health check for the specified application using the given health check configuration
func (resource *ResourcesData) SetHealthCheck(appName string, healthCheckType string, healthCheckHTTPEndpoint string, invocationTimeout int, process string) (err error) {
	if healthCheckType == "" {
		return nil
	}

	args := []string{"v3-set-health-check", appName}

	if healthCheckType == "http" && healthCheckHTTPEndpoint != "" {

		args = append(args, healthCheckType, "--endpoint", healthCheckHTTPEndpoint)
		if invocationTimeout >= 0 {
			args = append(args, "--invocation-timeout", string(invocationTimeout))
		}
	} else if process != "" && (healthCheckType == "process" || healthCheckType == "port") {
		args = append(args, healthCheckType, "--process", process)
	}

	ui.Say("apply health check timeouts")
	err = resource.Executor.Execute(args)
	if err != nil {
		ui.Failed("could not set health check timeouts", err)
		return err
	}
	ui.Ok()
	return nil
}
