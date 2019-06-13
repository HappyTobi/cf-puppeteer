package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"errors"
	"flag"
	"fmt"
	"github.com/happytobi/cf-puppeteer/cf"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"
	"github.com/happytobi/cf-puppeteer/manifest"
	"github.com/happytobi/cf-puppeteer/rewind"
	"github.com/happytobi/cf-puppeteer/ui"
	"os"
	"regexp"
	"strings"
)

func fatalIf(err error) {
	if err != nil {
		ui.Failed("error: %s", err)
		os.Exit(1)
	}
}

func main() {
	plugin.Start(&CfPuppeteerPlugin{})
}

type CfPuppeteerPlugin struct{}

func venerableAppName(appName string) string {
	return fmt.Sprintf("%s-venerable", appName)
}

func getActionsForApp(appRepo *ApplicationRepo, parsedArguments *ParserArguments) []rewind.Action {
	venName := venerableAppName(parsedArguments.AppName)
	var err error
	var curApp, venApp *v2.AppResourcesEntity
	var haveVenToCleanup bool

	return []rewind.Action{
		// get info about current app
		{
			Forward: func() error {
				curApp, err = appRepo.v2Resources.GetAppMetadata(parsedArguments.AppName)
				if err != nil {
					if err == v2.ErrAppNotFound {
						curApp = nil
					} else {
						return err
					}
				}
				return nil
			},
		},
		// get info about ven app
		{
			Forward: func() error {
				venApp, err = appRepo.v2Resources.GetAppMetadata(venName)
				if err != nil {
					if err == v2.ErrAppNotFound {
						venApp = nil
					} else {
						return err
					}
				}
				return nil
			},
		},
		// rename any existing app such so that next step can push to a clear space
		{
			Forward: func() error {
				// Unless otherwise specified, go with our start state
				haveVenToCleanup = (venApp != nil)

				// If there is no current app running, that's great, we're done here
				if curApp == nil {
					return nil
				}

				// If current app isn't started, then we'll just delete it, and we're done
				if curApp.Entity.State != "STARTED" {
					return appRepo.DeleteApplication(parsedArguments.AppName)
				}

				// Do we have a ven app that will stop a rename?
				if venApp != nil {
					// Finally, since the current app claims to be healthy, we'll delete the venerable app, and rename the current over the top
					err = appRepo.DeleteApplication(venName)
					if err != nil {
						return err
					}
				}

				// Finally, rename
				haveVenToCleanup = true
				return appRepo.RenameApplication(parsedArguments.AppName, venName)
			},
		},
		// push
		{
			Forward: func() error {
				//add v3 push
				space, err := appRepo.conn.GetCurrentSpace()
				if err != nil {
					return err
				}

				applicationBuildpacks := parsedArguments.Manifest.ApplicationManifests[0].Buildpacks
				applicationStack := parsedArguments.Manifest.ApplicationManifests[0].Stack
				appName := parsedArguments.AppName
				appPath := parsedArguments.AppPath
				serviceNames := parsedArguments.Manifest.ApplicationManifests[0].Services
				spaceGUID := space.Guid
				manifestPath := parsedArguments.ManifestPath
				routes := parsedArguments.Manifest.ApplicationManifests[0].Routes
				healthCheckType := parsedArguments.HealthCheckType
				healthCheckHttpEndpoint := parsedArguments.HealthCheckHTTPEndpoint
				process := parsedArguments.Process
				invocationTimeout := parsedArguments.InvocationTimeout

				//move to own function
				//TODO
				var mergedEnvs []string
				mergedEnvs = append(mergedEnvs, parsedArguments.Envs...)
				for k, v := range parsedArguments.Manifest.ApplicationManifests[0].Env {
					mergedEnvs = append(mergedEnvs, fmt.Sprintf("%s=%s", k, v))
				}

				var puppeteerPush cf.PuppeteerPush = cf.NewApplicationPush(appRepo.conn, appRepo.traceLogging)
				err = puppeteerPush.PushApplication(appName, venName, appPath, serviceNames, spaceGUID, applicationBuildpacks, applicationStack, mergedEnvs, manifestPath, routes, healthCheckType, healthCheckHttpEndpoint, process, invocationTimeout)
				if err != nil {
					return err
				}
				return nil
			},
		},
		// start
		{
			Forward: func() error {
				return appRepo.StartApplication(parsedArguments.AppName)
			},
			ReversePrevious: func() error {
				if !haveVenToCleanup {
					return nil
				}

				// If the app cannot start we'll have a lingering application
				// We delete this application so that the rename can succeed
				appRepo.DeleteApplication(parsedArguments.AppName)

				return appRepo.RenameApplication(venName, parsedArguments.AppName)
			},
		},
		// delete
		{
			Forward: func() error {
				if !haveVenToCleanup {
					return nil
				}

				//if vendorAppOption was set to stop
				if strings.ToLower(parsedArguments.VendorAppOption) == "stop" {
					return appRepo.StopApplication(venName)
				}
				return appRepo.DeleteApplication(venName)
			},
		},
	}
}

func (plugin CfPuppeteerPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// only handle if actually invoked, else it can't be uninstalled cleanly
	if args[0] != "zero-downtime-push" {
		return
	}

	var traceLogging bool
	if os.Getenv("CF_PUPPETEER_TRACE") == "true" {
		traceLogging = true
	}
	appRepo := NewApplicationRepo(cliConnection, traceLogging)
	parsedArguments, err := ParseArgs(appRepo, args)
	fatalIf(err)

	fatalIf((&rewind.Actions{
		Actions:              getActionsForApp(appRepo, parsedArguments),
		RewindFailureMessage: "Oh no. Something's gone wrong. I've tried to roll back but you should check to see if everything is OK.",
	}).Execute())

	fmt.Println()
	fmt.Println("A new version of your application has successfully been pushed!")
	fmt.Println()

	_ = appRepo.ListApplications()
}

// GetMetadata get plugin metadata
func (CfPuppeteerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cf-puppeteer",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "zero-downtime-push",
				HelpText: "Perform a zero-downtime push of an application over the top of an old one",
				UsageDetails: plugin.Usage{
					Usage: "$ cf zero-downtime-push [<App-Name>] -f <Manifest.yml> [options]",
					Options: map[string]string{
						"f":             "path to application manifest",
						"p":             "path to application files",
						"s":             "name of the stack to use",
						"t":             "push timeout (in secounds)",
						"-show-app-log": "tail and show application log during application start",
						"env":           "add environment key value pairs dynamic; can specify multiple times",
						//"var":                         "variable key value pair for variable substitution; can specify multiple times",
						//"vars-file":                   "Path to a variable substitution file for manifest; can specify multiple times",
						"-vendor-option":              "option to delete or stop vendor application - default is delete",
						"-health-check-type":          "type of health check to perform",
						"-health-check-http-endpoint": "endpoint for the 'http' health check type",
						"-invocation-timeout":         "timeout (in seconds) that controls individual health check invocations",
						//"-process":                    "application process to update",
						//"v":                           "print additional details on the deployment process",
					},
				},
			},
		},
	}
}

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprint(*s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

//ParserArguments struct where all arguments will be parsed into
type ParserArguments struct {
	AppName                 string
	ManifestPath            string
	AppPath                 string
	HealthCheckType         string
	HealthCheckHTTPEndpoint string
	Timeout                 int
	InvocationTimeout       int
	Process                 string
	StackName               string
	VendorAppOption         string
	Vars                    []string
	VarsFiles               []string
	Envs                    []string
	ShowLogs                bool
	DockerImage             string
	DockerUserName          string
	Manifest                manifest.Manifest
}

// ParseArgs parses the command line arguments
func ParseArgs(repo *ApplicationRepo, args []string) (*ParserArguments, error) {
	flags := flag.NewFlagSet("zero-downtime-push", flag.ContinueOnError)

	var envs stringSlice
	var vars stringSlice
	var varsFiles stringSlice

	pta := &ParserArguments{}
	flags.StringVar(&pta.ManifestPath, "f", "", "path to an application manifest")
	flags.StringVar(&pta.AppPath, "p", "", "path to application files")
	flags.StringVar(&pta.StackName, "s", "", "name of the stack to use")
	flags.StringVar(&pta.HealthCheckType, "health-check-type", "", "type of health check to perform")
	flags.StringVar(&pta.HealthCheckHTTPEndpoint, "health-check-http-endpoint", "", "endpoint for the 'http' health check type")
	flags.IntVar(&pta.Timeout, "t", 0, "push timeout in secounds (defaults to 60 seconds)")
	flags.IntVar(&pta.InvocationTimeout, "invocation-timeout", -1, "health check invocation timeout in seconds")
	//flags.StringVar(&pta.Process, "process", "", "application process to update")
	flags.BoolVar(&pta.ShowLogs, "show-app-log", false, "tail and show application log during application start")
	flags.StringVar(&pta.VendorAppOption, "vendor-option", "delete", "option to delete or stop vendor application - default is delete")
	flags.Var(&envs, "env", "Variable key value pair for adding dynamic environment variables; can specity multiple times")
	//flags.Var(&vars, "var", "Variable key value pair for variable substitution, (e.g., name=app1); can specify multiple times")
	//flags.Var(&varsFiles, "vars-file", "Path to a variable substitution file for manifest; can specify multiple times")
	//flags.StringVar(&pta.DockerImage, "docker-image", "", "url to docker image")
	//flags.StringVar(&pta.DockerUserName, "docker-username", "", "pass docker username if image came from private repository")
	//dockerPass := os.Getenv("CF_DOCKER_PASSWORD")

	//first check if argument was passed
	if len(args) < 2 {
		return pta, ErrNoArgument
	}

	//default start index of parameters is 2 because 1 is the appName
	argumentStartIndex := 2
	//if first argument is not the appName we have to read the appName out of the manifest
	noAppNameProvided, _ := regexp.MatchString("^-[a-z]{0,3}", args[1])
	//noAppNameProvided := strings.Contains(args[1], "-")
	if noAppNameProvided {
		argumentStartIndex = 1
	}

	err := flags.Parse(args[argumentStartIndex:])
	if err != nil {
		return pta, err
	}

	if pta.ManifestPath == "" {
		return pta, ErrNoManifest
	}

	//parse manifest
	parsedManifest, err := manifest.Parse(pta.ManifestPath)
	if err != nil {
		return pta, ErrManifest
	}
	pta.Manifest = parsedManifest

	/*if *dockerImage != "" && *dockerUserName != "" && dockerPass != "" {
	    //TODO use dockerImage stuff and pass to push command
	}*/

	//set timeout
	manifestTimeout := parsedManifest.ApplicationManifests[0].Timeout
	if manifestTimeout > 0 && pta.Timeout <= 0 {
		pta.Timeout = manifestTimeout
	} else if manifestTimeout <= 0 && pta.Timeout <= 0 {
		pta.Timeout = 60
	}

	//parse first argument as appName
	pta.AppName = args[1]
	if noAppNameProvided {
		pta.AppName = parsedManifest.ApplicationManifests[0].Name
	}

	// get health check settings from manifest if nothing else was specified in the command line
	if pta.HealthCheckType == "" {
		if parsedManifest.ApplicationManifests[0].HealthCheckType == "" {
			pta.HealthCheckType = "port"
		} else {
			pta.HealthCheckType = parsedManifest.ApplicationManifests[0].HealthCheckType
		}

	}
	if pta.HealthCheckHTTPEndpoint == "" {
		pta.HealthCheckHTTPEndpoint = parsedManifest.ApplicationManifests[0].HealthCheckHTTPEndpoint
	}

	//validate envs format
	if len(envs) > 0 {
		for _, envPair := range envs {
			if strings.Contains(envPair, "=") == false {
				return pta, ErrWrongEnvFormat
			}
		}
	}

	pta.Envs = envs
	pta.Vars = vars
	pta.VarsFiles = varsFiles

	return pta, nil
}

//all custom errors
var (
	//ErrNoArgument error when zero-downtime-push without a argument called
	ErrNoArgument = errors.New("no valid argument found, use --help / -h for more information")
	//ErrNoManifest error when manifes on push application was not found
	ErrNoManifest = errors.New("a manifest is required to push this application")
	//ErrManifest error when manifes could not be parsed
	ErrManifest = errors.New("could not parse manifest")
	//ErrWrongEnvFormat error when env files was not in right format
	ErrWrongEnvFormat = errors.New("--var would be in wrong format, use the vars like key=value")
	//ErrAppNotFound application not found error
	ErrAppNotFound = errors.New("application not found")
)

type ApplicationRepo struct {
	conn         plugin.CliConnection
	traceLogging bool
	v2Resources  v2.Resources
}

func NewApplicationRepo(conn plugin.CliConnection, traceLogging bool) *ApplicationRepo {
	return &ApplicationRepo{
		conn:         conn,
		traceLogging: traceLogging,
		v2Resources:  v2.NewV2Resources(conn, traceLogging),
	}
}

func (repo *ApplicationRepo) RenameApplication(oldName, newName string) error {
	_, err := repo.conn.CliCommand("rename", oldName, newName)
	return err
}

func (repo *ApplicationRepo) StopApplication(appName string) error {
	_, err := repo.conn.CliCommand("stop", appName)
	return err
}

func (repo *ApplicationRepo) StartApplication(appName string) error {
	_, err := repo.conn.CliCommand("start", appName)
	return err
}

func (repo *ApplicationRepo) DeleteApplication(appName string) error {
	_, err := repo.conn.CliCommand("delete", appName, "-f")
	return err
}

func (repo *ApplicationRepo) ListApplications() error {
	_, err := repo.conn.CliCommand("apps")
	return err
}
