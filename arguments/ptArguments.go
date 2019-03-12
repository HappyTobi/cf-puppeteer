package arguments

import (
	"flag"
	"os"
	"regexp"
	"strings"

	"github.com/happytobi/cf-puppeteer/manifest"
)

//Interface how to use the Puppeteer Arguments parser
type Interface interface {
	ParseArgs(args []string) (*PTArguments, error)
	Usage() string
}

//PTArguments struct where all found args will be parsed into
type PTArguments struct {
	appName                 string
	manifestPath            string
	appPath                 string
	healthCheckType         string
	healthCheckHTTPEndpoint string
	timeout                 int
	invocationTimeout       int
	process                 string
	stackName               string
	vendorAppOption         string
	vars                    []string
	varFiles                []string
	envs                    []string
	showLogs                bool
	dockerImage             string
	dockerUserName          string
}

//NewPTArgumentsParser create a new instance of Pupetter argument parser
func NewPTArgumentsParser() *PTArguments {
	return &PTArguments{}
}

//ParseArgs parse all arguments into the PTArguments struct
func (pta *PTArguments) ParseArgs(args []string) (*PTArguments, error) {

	flags := flag.NewFlagSet("zero-downtime-push", flag.ContinueOnError)

	flags.StringVar(&pta.manifestPath, "f", "", "path to an application manifest")
	flags.StringVar(&pta.appPath, "p", "", "path to application files")
	flags.StringVar(&pta.stackName, "s", "", "name of the stack to use")
	flags.StringVar(&pta.healthCheckType, "health-check-type", "", "type of health check to perform")
	flags.StringVar(&pta.healthCheckHTTPEndpoint, "health-check-http-endpoint", "", "endpoint for the 'http' health check type")
	flags.IntVar(&pta.timeout, "t", 0, "push timeout in secounds (defaults to 60 seconds)")
	flags.IntVar(&pta.invocationTimeout, "invocation-timeout", -1, "health check invocation timeout in seconds")
	flags.StringVar(&pta.process, "process", "", "application process to update")
	flags.BoolVar(&pta.showLogs, "show-app-log", false, "tail and show application log during application start")
	flags.StringVar(&pta.vendorAppOption, "vendor-option", "delete", "option to delete or stop vendor application - default is delete")
	/*flags.StringVar(&pta.envs, "env", "Variable key value pair for adding dynamic environment variables; can specity multiple times")
	flags.StringVar(&pta.vars, "var", "Variable key value pair for variable substitution, (e.g., name=app1); can specify multiple times")
	flags.StringVar(&pta.varsFiles, "vars-file", "Path to a variable substitution file for manifest; can specify multiple times")*/
	flags.StringVar(&pta.dockerImage, "docker-image", "", "url to docker image")
	flags.StringVar(&pta.dockerUserName, "docker-username", "", "pass docker username if image came from private repository")
	dockerPass := os.Getenv("CF_DOCKER_PASSWORD")

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

	if pta.manifestPath == "" {
		return pta, ErrNoManifest
	}

	//parse manifest
	parsedManifest, err := manifest.Parse(pta.manifestPath)
	if err != nil {
		return pta, ErrManifest
	}

	/*if *dockerImage != "" && *dockerUserName != "" && dockerPass != "" {
		//TODO use dockerImage stuff and pass to push command
	}*/

	//set timeout
	manifestTimeout := parsedManifest.ApplicationManifests[0].Timeout
	if manifestTimeout > 0 && pta.timeout <= 0 {
		pta.timeout = manifestTimeout
	} else if manifestTimeout <= 0 && pta.timeout <= 0 {
		pta.timeout = 60
	}

	//parse first argument as appName
	appName := args[1]
	if noAppNameProvided {
		appName = parsedManifest.ApplicationManifests[0].Name
	}

	// get health check settings from manifest if nothing else was specified as command line argument
	if pta.healthCheckType != "" || pta.healthCheckHTTPEndpoint != "" || pta.process != "" {
		//TODO check it here?
		//err = repo.CheckAPIV3()
		if err != nil {
			return pta, ErrNoV3ApiAvailable
		}
	}

	//TODO check invocationtimeout argument

	//validate envs format
	if len(pta.envs) > 0 {
		for _, envPair := range pta.envs {
			if strings.Contains(envPair, "=") == false {
				return pta, ErrWrongEnvFormat
			}
		}
	}

	return pta, nil
}
