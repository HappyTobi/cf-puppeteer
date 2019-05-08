package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/cf/api/logs"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/happytobi/cf-puppeteer/cfResources"
	"github.com/happytobi/cf-puppeteer/manifest"
	"github.com/happytobi/cf-puppeteer/rewind"
	"github.com/happytobi/cf-puppeteer/v2cfResources"
)

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
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
	var curApp, venApp, temp *AppResourcesEntity
	var haveVenToCleanup bool

	return []rewind.Action{
		// get info about current app
		{
			Forward: func() error {
				curApp, err = appRepo.GetAppMetadata(parsedArguments.AppName)
				if err != ErrAppNotFound {
					return err
				}
				curApp = nil
				return nil
			},
		},
		// get info about ven app
		{
			Forward: func() error {
				venApp, err = appRepo.GetAppMetadata(venName)
				if err != ErrAppNotFound {
					return err
				}
				venApp = nil
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

				appResponse, err := appRepo.cf.PushApp(parsedArguments.AppName, space.Guid, applicationBuildpacks, applicationStack)
				if err != nil {
					return err
				}

				appRepo.cf.AssignAppManifest(appResponse.Links.Self.Href, parsedArguments.ManifestPath)

				createPackageResponse, err := appRepo.cf.CreatePackage(appResponse.GUID)
				if err != nil {
					return err
				}

				domains, err := appRepo.cf.GetDomain(parsedArguments.Manifest.ApplicationManifests[0].Routes)
				if err != nil {
					return err
				}

				for _, route := range *domains {
					routeResponse, err := appRepo.cf2.CreateRoute(space.Guid, route.DomainGUID, route.Host)
					if err != nil {
						return err
					}
					err = appRepo.cf.RouteMapping(appResponse.GUID, routeResponse.Metadata.GUID)
					if err != nil {
						return err
					}
					fmt.Printf("route generated and added to application - host: %s - domain: %s \n", route.Host, route.DomainGUID)
				}

				createPackageResponse, err = appRepo.cf.UploadApplication(parsedArguments.AppName, parsedArguments.AppPath, createPackageResponse.Links.Upload.Href)
				if err != nil {
					return err
				}

				duration, _ := time.ParseDuration("5s")
				fmt.Printf("Upload application [")
				for createPackageResponse.State != "FAILED" &&
					createPackageResponse.State != "READY" &&
					createPackageResponse.State != "EXPIRED" {
					time.Sleep(duration)
					createPackageResponse, err = appRepo.cf.CheckPackageState(createPackageResponse.GUID)
					fmt.Printf("=")
					if err != nil {
						return nil
					}
				}

				fmt.Printf("] - done \n")

				buildResponse, err := appRepo.cf.CreateBuild(createPackageResponse.GUID)
				if err != nil {
					return err
				}

				fmt.Printf("waiting while creating the application [")
				for buildResponse.State != "FAILED" &&
					buildResponse.State != "STAGED" {
					time.Sleep(duration)
					buildResponse, err = appRepo.cf.CheckBuildState(buildResponse.GUID)
					fmt.Printf("=")
					if err != nil {
						return nil
					}
				}
				fmt.Printf("] - done \n")

				dropletResponse, err := appRepo.cf.GetDropletGUID(buildResponse.GUID)
				fmt.Printf("get droplet guid %s \n", dropletResponse)

				err = appRepo.cf.AssignApp(appResponse.GUID, dropletResponse.GUID)
				fmt.Printf("app assigned \n")

				/*err = appRepo.cf.StartApp(appResponse.GUID)
				if err != nil {
					return err
				}*/
				return nil
			},
		},
		{
			Forward: func() error {
				temp, err = appRepo.GetAppMetadata(parsedArguments.AppName)
				if err != nil {
					return err
				}
				return appRepo.SetHealthCheckV3(parsedArguments, temp.Metadata.GUID)
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
	if os.Getenv("CF_PUPPETEER_TRACE") != "true" {
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
						"f":                           "path to application manifest",
						"p":                           "path to application files",
						"s":                           "name of the stack to use",
						"t":                           "push timeout (in secounds)",
						"-show-app-log":               "tail and show application log during application start",
						"env":                         "add environment key value pairs dynamic; can specity multiple times",
						"var":                         "variable key value pair for variable substitution; can specify multiple times",
						"vars-file":                   "Path to a variable substitution file for manifest; can specify multiple times",
						"-vendor-option":              "option to delete or stop vendor application - default is delete",
						"-health-check-type":          "type of health check to perform",
						"-health-check-http-endpoint": "endpoint for the 'http' health check type",
						"-invocation-timeout":         "timeout (in seconds) that controls individual health check invocations",
						"-process":                    "application process to update",
						"v":                           "print additional details on the deployment process",
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
	flags.StringVar(&pta.Process, "process", "", "application process to update")
	flags.BoolVar(&pta.ShowLogs, "show-app-log", false, "tail and show application log during application start")
	flags.StringVar(&pta.VendorAppOption, "vendor-option", "delete", "option to delete or stop vendor application - default is delete")
	flags.Var(&envs, "env", "Variable key value pair for adding dynamic environment variables; can specity multiple times")
	flags.Var(&vars, "var", "Variable key value pair for variable substitution, (e.g., name=app1); can specify multiple times")
	flags.Var(&varsFiles, "vars-file", "Path to a variable substitution file for manifest; can specify multiple times")
	flags.StringVar(&pta.DockerImage, "docker-image", "", "url to docker image")
	flags.StringVar(&pta.DockerUserName, "docker-username", "", "pass docker username if image came from private repository")
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
		pta.HealthCheckType = parsedManifest.ApplicationManifests[0].HealthCheckType
	}
	if pta.HealthCheckHTTPEndpoint == "" {
		pta.HealthCheckHTTPEndpoint = parsedManifest.ApplicationManifests[0].HealthCheckHTTPEndpoint
	}

	if pta.HealthCheckType != "" || pta.HealthCheckHTTPEndpoint != "" || pta.Process != "" {
		err = repo.CheckAPIV3()
		if err != nil {
			return pta, ErrNoV3ApiAvailable
		}
	}

	//TODO check invocationtimeout argument

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
	//ErrNoV3ApiAvailable error message when cf api v3 is not available
	ErrNoV3ApiAvailable = errors.New("cf api v3 is not available")
	//ErrInvokationTimeout
	ErrInvokationTimeout = errors.New("could not set invocation timeout to application")
	//ErrWrongInvocationTimeoutArgs
	ErrWrongInvocationTimeoutArgs = errors.New("wrong combination of timeout arguments passed")
)

type ApplicationRepo struct {
	conn         plugin.CliConnection
	traceLogging bool
	cf           cfResources.CfResourcesInterface
	cf2          v2cfResources.V2CfResourcesInterface
}

func NewApplicationRepo(conn plugin.CliConnection, traceLogging bool) *ApplicationRepo {
	return &ApplicationRepo{
		conn:         conn,
		traceLogging: traceLogging,
		cf:           cfResources.NewResources(conn, traceLogging),
		cf2:          v2cfResources.NewResources(conn, traceLogging),
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

// SetHealthCheckV3 sets the health check for the specified application using the given health check configuration
func (repo *ApplicationRepo) SetHealthCheckV3(parsedArguments *ParserArguments, GUID string) error {
	// Without a health check type, the CF command is not valid. Therefore, leave if the type is not specified
	if parsedArguments.HealthCheckType == "" {
		return nil
	}

	// load application by guid
	appProcesEntity, err := repo.GetApplicationProcessWebInformation(GUID)

	applicationEntity := ApplicationEntityV3{}
	applicationEntity.Command = appProcesEntity.Command
	applicationEntity.HealthCheck.HealthCheckType = parsedArguments.HealthCheckType

	if parsedArguments.HealthCheckType == "http" && parsedArguments.HealthCheckHTTPEndpoint != "" {
		applicationEntity.HealthCheck.Data.Endpoint = parsedArguments.HealthCheckHTTPEndpoint
		if parsedArguments.InvocationTimeout >= 0 {
			applicationEntity.HealthCheck.Data.InvocationTimeout = parsedArguments.InvocationTimeout
		}
	} else if parsedArguments.Process != "" && (parsedArguments.HealthCheckType == "process" || parsedArguments.HealthCheckType == "port") {
		applicationEntity.ProcessType = parsedArguments.Process
	} else {
		return ErrWrongInvocationTimeoutArgs
	}

	fmt.Println("")
	fmt.Printf("Updating health check setting for application %v", parsedArguments.AppName)
	fmt.Println("")
	err = repo.UpdateApplicationProcessWebInformation(appProcesEntity.GUID, applicationEntity)
	return err
}

// PushApplication executes the Cloud Foundry push command for the specified application.
// It returns any error that prevents a successful completion of the operation.
func (repo *ApplicationRepo) PushApplication(parsedArguments *ParserArguments) error {
	args := []string{"push", parsedArguments.AppName, "-f", parsedArguments.ManifestPath, "--no-start"}

	if parsedArguments.AppPath != "" {
		args = append(args, "-p", parsedArguments.AppPath)
	}

	if parsedArguments.StackName != "" {
		args = append(args, "-s", parsedArguments.StackName)
	}

	/* always append timeout */
	timeout := strconv.Itoa(parsedArguments.Timeout)
	args = append(args, "-t", timeout)
	for _, varPair := range parsedArguments.Vars {
		args = append(args, "--var", varPair)
	}

	for _, varsFile := range parsedArguments.VarsFiles {
		args = append(args, "--vars-file", varsFile)
	}

	_, err := repo.conn.CliCommand(args...)
	if err != nil {
		return err
	}

	envErr := repo.setEnvironmentVariables(parsedArguments.AppName, parsedArguments.Envs)
	if envErr != nil {
		return envErr
	}

	if parsedArguments.ShowLogs {
		app, err := repo.conn.GetApp(parsedArguments.AppName)
		if err != nil {
			return err
		}
		dopplerEndpoint, err := repo.conn.DopplerEndpoint()
		if err != nil {
			return err
		}
		token, err := repo.conn.AccessToken()
		if err != nil {
			return err
		}

		cons := consumer.New(dopplerEndpoint, nil, nil)
		defer cons.Close()

		messages, errors := cons.TailingLogs(app.Guid, token)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			for {
				select {
				case m := <-messages:
					if m.GetSourceType() != "STG" { // skip STG messages as the cf tool already prints them
						os.Stderr.WriteString(logs.NewNoaaLogMessage(m).ToLog(time.Local) + "\n")
					}
				case e := <-errors:
					log.Println("error reading logs:", e)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return nil
}

// setEnvironmentVariables sets passed envs with set-env to set variables dynamically
func (repo *ApplicationRepo) setEnvironmentVariables(appName string, envs []string) error {
	varArgs := []string{"set-env", appName}
	//set all variables passed by --var
	for _, envPair := range envs {
		tmpArgs := make([]string, len(varArgs))
		copy(tmpArgs, varArgs)
		newArgs := strings.SplitN(envPair, "=", 2)
		tmpArgs = append(tmpArgs, newArgs...)
		_, err := repo.conn.CliCommand(tmpArgs...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *ApplicationRepo) DeleteApplication(appName string) error {
	_, err := repo.conn.CliCommand("delete", appName, "-f")
	return err
}

func (repo *ApplicationRepo) ListApplications() error {
	_, err := repo.conn.CliCommand("apps")
	return err
}

type MetaDataEntity struct {
	AppResourcesEntity []AppResourcesEntity `json:"resources"`
}
type Metadata struct {
	GUID string `json:"guid"`
}
type Entity struct {
	Name  string `json:"name"`
	State string `json:"state"`
}
type AppResourcesEntity struct {
	Metadata Metadata `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

// GetAppMetadata fetches information on the application specified by appName
func (repo *ApplicationRepo) GetAppMetadata(appName string) (*AppResourcesEntity, error) {
	space, err := repo.conn.GetCurrentSpace()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf(`v2/apps?q=name:%s&q=space_guid:%s`, url.QueryEscape(appName), space.Guid)
	result, err := repo.conn.CliCommandWithoutTerminalOutput("curl", path)

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	var metaDataResponseEntity MetaDataEntity
	err = json.Unmarshal([]byte(jsonResp), &metaDataResponseEntity)

	if err != nil {
		return nil, err
	}

	if len(metaDataResponseEntity.AppResourcesEntity) == 0 {
		return nil, ErrAppNotFound
	}

	return &metaDataResponseEntity.AppResourcesEntity[0], nil
}

/**
V3 Entities
*/

type ApplicationProcessesEntityV3 struct {
	GUID    string `json:"guid"`
	Command string `json:"command,omitempty"`
}

type ApplicationEntityV3 struct {
	Command     string              `json:"command,omitempty"`
	HealthCheck HealthCheckEntityV3 `json:"health_check"`
	ProcessType string              `json:"type,omitempty"`
}
type DataEnvityV3 struct {
	Endpoint          string `json:"endpoint,omitempty"`
	InvocationTimeout int    `json:"invocation_timeout,omitempty"`
}
type HealthCheckEntityV3 struct {
	Data            DataEnvityV3 `json:"data,omitempty"`
	HealthCheckType string       `json:"type"`
}

// CheckAPIV3 call v3 url to check availablility
func (repo *ApplicationRepo) CheckAPIV3() error {
	response, err := repo.conn.CliCommandWithoutTerminalOutput("curl", "/v3", "-X", "GET")
	result := strings.Join(response, "")

	if err != nil || strings.Contains(result, "error") {
		return ErrNoV3ApiAvailable
	}
	return nil
}

// GetApplicationProcessWebInformation call v3 process api
func (repo *ApplicationRepo) GetApplicationProcessWebInformation(appGUID string) (*ApplicationProcessesEntityV3, error) {
	path := fmt.Sprintf(`/v3/apps/%s/processes/web`, appGUID)
	result, err := repo.conn.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET")

	if err != nil {
		return nil, err
	}

	jsonResp := strings.Join(result, "")

	if repo.traceLogging {
		fmt.Printf("Cloud Foundry API response to GET call on %s:\n", path)
		PrettyPrintJSON(jsonResp)
	}

	var applicationProcessResponse ApplicationProcessesEntityV3
	err = json.Unmarshal([]byte(jsonResp), &applicationProcessResponse)
	if err != nil {
		return nil, err
	}

	if len(applicationProcessResponse.GUID) == 0 {
		return nil, ErrAppNotFound
	}

	return &applicationProcessResponse, nil
}

// UpdateApplicationProcessWebInformation calls v3 application to set options
// see api documentation http://v3-apidocs.cloudfoundry.org/version/3.67.0/index.html#update-an-app
func (repo *ApplicationRepo) UpdateApplicationProcessWebInformation(appGUID string, applicationEntity ApplicationEntityV3) error {
	path := fmt.Sprintf(`/v3/processes/%s`, appGUID)
	appJSON, err := json.Marshal(applicationEntity)
	if err != nil {
		return err
	}

	result, err := repo.conn.CliCommandWithoutTerminalOutput("curl", path, "-X", "PATCH", "-H", "Content-type: application/json",
		"-d", string(appJSON))

	if err != nil {
		return err
	}

	jsonResp := strings.Join(result, "")

	if repo.traceLogging {
		fmt.Printf("Cloud Foundry API response to PATCH call on %s:\n", path)
		PrettyPrintJSON(jsonResp)
	}

	var applicationResponse ApplicationEntityV3
	err = json.Unmarshal([]byte(jsonResp), &applicationResponse)

	if err != nil {
		return err
	}

	if applicationResponse.HealthCheck.HealthCheckType != applicationEntity.HealthCheck.HealthCheckType {
		return ErrInvokationTimeout
	}

	return nil
}

// PrettyPrintJSON takes the given JSON string, makes it pretty, and prints it out.
func PrettyPrintJSON(jsonUgly string) error {
	jsonPretty := &bytes.Buffer{}
	err := json.Indent(jsonPretty, []byte(jsonUgly), "", "    ")

	if err != nil {
		return err
	}

	fmt.Println(jsonPretty.String())

	return nil
}
