package arguments

import (
	"errors"
	"flag"
	"fmt"
	"github.com/happytobi/cf-puppeteer/manifest"
	"github.com/happytobi/cf-puppeteer/ui"
	"regexp"
	"strings"
)

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
	VenerableAction         string
	Envs                    []string
	ShowLogs                bool
	ShowCrashLogs           bool
	DockerImage             string
	DockerUserName          string
	Manifest                manifest.Manifest
	MergedEnvs              []string
	LegacyPush              bool
	NoRoute                 bool
	AddRoutes               bool
	NoStart                 bool
}

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprint(*s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var (
	//ErrNoArgument error when zero-downtime-push without a argument called
	ErrNoArgument = errors.New("no valid argument found, use --help / -h for more information")
	//ErrNoManifest error when manifes on push application was not found
	ErrNoManifest = errors.New("a manifest is required to push this application")
	//ErrWrongEnvFormat error when env files was not in right format
	ErrWrongEnvFormat = errors.New("--var would be in wrong format, use the vars like key=value")
	//ErrWrongCombination error when legacy push is used with health check options
	ErrWrongCombination = errors.New("--legacy-push and health check options coiuldn't be combined")
)

// ParseArgs parses the command line arguments
func ParseArgs(args []string) (*ParserArguments, error) {
	flags := flag.NewFlagSet("zero-downtime-push", flag.ContinueOnError)

	var envs stringSlice

	pta := &ParserArguments{}
	flags.StringVar(&pta.ManifestPath, "f", "", "path to an application manifest")
	flags.StringVar(&pta.AppPath, "p", "", "path to application files")
	flags.StringVar(&pta.StackName, "s", "", "name of the stack to use")
	flags.StringVar(&pta.HealthCheckType, "health-check-type", "", "type of health check to perform")
	flags.StringVar(&pta.HealthCheckHTTPEndpoint, "health-check-http-endpoint", "", "endpoint for the 'http' health check type")
	flags.IntVar(&pta.Timeout, "t", 0, "push timeout in seconds (defaults to 60 seconds)")
	flags.IntVar(&pta.InvocationTimeout, "invocation-timeout", -1, "health check invocation timeout in seconds")
	flags.StringVar(&pta.Process, "process", "", "application process to update")
	flags.BoolVar(&pta.ShowLogs, "show-app-log", false, "tail and show application log during application start")
	flags.BoolVar(&pta.ShowCrashLogs, "show-crash-log", false, "Show recent logs when applications crashes while the deployment")
	flags.StringVar(&pta.VendorAppOption, "vendor-option", "delete", "option to delete,stop,none application action on vendor app- default is delete")
	flags.StringVar(&pta.VenerableAction, "venerable-action", "delete", "option to delete,stop,none application action on vendor app- default is delete")
	flags.Var(&envs, "env", "Variable key value pair for adding dynamic environment variables; can specify multiple times")
	flags.BoolVar(&pta.LegacyPush, "legacy-push", false, "use legacy push instead of new v3 api")
	flags.BoolVar(&pta.NoRoute, "no-route", false, "deploy new application without adding routes")
	flags.BoolVar(&pta.AddRoutes, "route-only", false, "only add routes from manifest to the application")
	flags.BoolVar(&pta.NoStart, "no-start", false, "don't start application after deployment")
	//flags.BoolVar(&pta.ShowLogs, "show-app-log", false, "tail and show application log during application start")
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
		return pta, err //ErrManifest
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

	//check that health check works without legacy push only
	if pta.LegacyPush && (pta.HealthCheckType != "" || pta.HealthCheckHTTPEndpoint != "") {
		return nil, ErrWrongCombination
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

	//Merge all environment variables together
	var mergedEnvs []string
	mergedEnvs = append(mergedEnvs, pta.Envs...)
	for k, v := range pta.Manifest.ApplicationManifests[0].Env {
		mergedEnvs = append(mergedEnvs, fmt.Sprintf("%s=%s", k, v))
	}
	pta.MergedEnvs = mergedEnvs

	//print waring for deprecated arguments
	if strings.ToLower(pta.VendorAppOption) != "delete" {
		ui.Warn("deprecated argument used, please use --venerable-action instead - argument will dropped in next version")
		pta.VenerableAction = pta.VendorAppOption
	}

	//no-route set venerable-action to delete as default - but can be overwritten
	if pta.NoRoute && argPassed(flags, "venerable-action") == false {
		pta.VenerableAction = "none"
	}

	return pta, nil
}

//search vor argument in name in passed args
func argPassed(flags *flag.FlagSet, name string) (found bool) {
	found = false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
