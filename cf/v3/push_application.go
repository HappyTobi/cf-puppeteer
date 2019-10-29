package v3

import (
	"code.cloudfoundry.org/cli/cf/appfiles"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"
	"github.com/happytobi/cf-puppeteer/manifest"
	"github.com/happytobi/cf-puppeteer/ui"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"strconv"
)

//Push interface with all v3 actions
type Push interface {
	PushApplication(venAppName string, spaceGUID string, parsedArguments *arguments.ParserArguments, v2Resources v2.Resources) error
	SwitchRoutesOnly(venAppName string, appName string, routes []map[string]string) error
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

//PushApplication call all methods to push a complete application
func (resource *ResourcesData) PushApplication(venAppName, spaceGUID string, parsedArguments *arguments.ParserArguments, v2Resources v2.Resources) error {

	ui.Say("create application %s", parsedArguments.AppName)
	err := resource.CreateApp(parsedArguments)
	if err != nil {
		return err
	}

	ui.Say("apply manifest file")
	manifestPath := parsedArguments.ManifestPath
	if parsedArguments.NoRoute {
		newManifestPath, err := resource.GenerateNoRouteYml(parsedArguments.AppName, parsedArguments.Manifest)
		if err != nil {
			return errors.Wrap(err, "could not generate a new temp manifest without routes")
		}
		ui.Say("use no route manifest")
		manifestPath = newManifestPath
	}

	err = resource.AssignAppManifest(manifestPath)
	if err != nil {
		return err
	}

	ui.Say("push application")
	err = resource.PushApp(parsedArguments)
	if err != nil {
		return err
	}

	ui.Say("set health-check with type: %s for application %s", parsedArguments.HealthCheckType, parsedArguments.AppName)
	err = resource.SetHealthCheck(parsedArguments.AppName, parsedArguments.HealthCheckType, parsedArguments.HealthCheckHTTPEndpoint, parsedArguments.InvocationTimeout, parsedArguments.Process)
	if err != nil {
		return err
	}
	ui.Ok()

	return nil
}

//SwitchRoutes switch route interface method to provide switch routes only option
func (resource *ResourcesData) SwitchRoutesOnly(venAppName string, appName string, routes []map[string]string) (err error) {
	return resource.SwitchRoutes(venAppName, appName, routes)
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
func (resource *ResourcesData) SwitchRoutes(venAppName string, appName string, routes []map[string]string) (err error) {
	domains, err := resource.GetDomain(routes)
	if err != nil {
		return err
	}

	ui.Say("map routes to new application %s", appName)
	for _, route := range *domains {
		err = resource.MapRoute(appName, route.Host, route.Domain)
		if err != nil {
			//loop through
			ui.Warn("could not map route %s.%s to application", route.Host, route.Domain, appName)
		}
	}

	ui.Say("remove routes from venerable application %s", venAppName)
	for _, route := range *domains {
		err = resource.UnMapRoute(venAppName, route.Host, route.Domain)
		if err != nil {
			//loop through
			ui.Warn("could not remove route %s.%s from application", route.Host, route.Domain, venAppName)
		}
	}

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
			args = append(args, "--invocation-timeout", strconv.Itoa(invocationTimeout))
		}
	} else if process != "" && healthCheckType == "process" {
		args = append(args, healthCheckType, "--process", process)
	} else if healthCheckType == "port" {
		args = append(args, healthCheckType)
	}

	err = resource.Executor.Execute(args)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not set healthcheck with type: %s - endpoint: %s - invocationTimeout %v", healthCheckType, healthCheckHTTPEndpoint, invocationTimeout))
	}
	return nil
}
