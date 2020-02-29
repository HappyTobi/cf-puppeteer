package arguments

import (
	"errors"
	"flag"
	"fmt"
	"github.com/happytobi/cf-puppeteer/cf/utils/env"
	"github.com/happytobi/cf-puppeteer/manifest"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//ParserArguments struct where all arguments will be parsed into
type ParserArguments struct {
	AppName                 string
	ManifestPath            string
	NoRouteManifestPath     string
	AppPath                 string
	HealthCheckType         string
	HealthCheckHTTPEndpoint string
	Timeout                 int
	InvocationTimeout       int
	Process                 string
	StackName               string
	VenerableAction         string
	Envs                    map[string]string
	ShowCrashLogs           bool
	DockerImage             string
	DockerUserName          string
	Manifest                manifest.Manifest
	LegacyPush              bool
	NoRoute                 bool
	AddRoutes               bool
	NoStart                 bool
	VarsFile                string
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
	ErrNoManifest = errors.New("a application manifest is required to push an application")
	//ErrWrongEnvFormat error when env files was not in right format
	ErrWrongEnvFormat = errors.New("environment variables passed in wrong format, pass the variables like key=value")
	//ErrWrongCombination error when legacy push is used with health check options
	ErrWrongCombination = errors.New("--legacy-push and health check options couldn't be combined")
	//ErrWrongDockerCombination error when private docker image repo will be pushed without a pass
	ErrWrongPrivateDockerRepoCombination = errors.New("--docker-username have to be used in combination with env CF_DOCKER_PASSWORD and --docker-image")
	//Error manifest error when a wildcard was in the path directive
	ErrNoWildcardSupport = errors.New("wildcard expressions within the path directive in the application manifest are not supported - delete this path directive and pass the artifact path by using the -p option")
)

// ParseArgs parses the command line arguments
func ParseArgs(args []string) (*ParserArguments, error) {
	flags := flag.NewFlagSet("zero-downtime-push", flag.ContinueOnError)

	var envs stringSlice
	commands := commands()

	pta := &ParserArguments{}
	flags.StringVar(&pta.ManifestPath, "f", "", commands["f"])
	flags.StringVar(&pta.AppPath, "p", "", commands["p"])
	flags.StringVar(&pta.StackName, "s", "", commands["s"])
	flags.StringVar(&pta.HealthCheckType, "health-check-type", "port", commands["health-check-type"])
	flags.StringVar(&pta.HealthCheckHTTPEndpoint, "health-check-http-endpoint", "", commands["health-check-http-endpoint"])
	flags.IntVar(&pta.Timeout, "t", 0, commands["t"])
	flags.IntVar(&pta.InvocationTimeout, "invocation-timeout", -1, commands["invocation-timeout"])
	flags.StringVar(&pta.Process, "process", "", commands["process"])
	flags.BoolVar(&pta.ShowCrashLogs, "show-crash-log", false, commands["show-crash-log"])
	flags.StringVar(&pta.VenerableAction, "venerable-action", "delete", commands["venerable-action"])
	flags.Var(&envs, "env", commands["env"])
	flags.BoolVar(&pta.LegacyPush, "legacy-push", false, commands["legacy-push"])
	flags.BoolVar(&pta.NoRoute, "no-route", false, commands["no-route"])
	flags.BoolVar(&pta.AddRoutes, "route-only", false, commands["route-only"])
	flags.BoolVar(&pta.NoStart, "no-start", false, commands["no-start"])
	flags.StringVar(&pta.VarsFile, "vars-file", "", commands["vars-file"])
	flags.StringVar(&pta.DockerImage, "docker-image", "", commands["docker-image"])
	flags.StringVar(&pta.DockerUserName, "docker-username", "", commands["docker-username"])
	dockerPass := os.Getenv("CF_DOCKER_PASSWORD")

	//first check if argument was passed
	if len(args) < 2 {
		return pta, ErrNoArgument
	}

	//default start index of parameters is 2 because 1 is the appName
	argumentStartIndex := 2
	//if first argument is not the appName we have to read the appName out of the manifest
	noAppNameProvided, _ := regexp.MatchString("^-[a-z]{0,3}", args[1])

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
	parsedManifest, noRouteManifestPath, err := manifest.ParseApplicationManifest(pta.ManifestPath, pta.VarsFile)
	if err != nil {
		return pta, err //ErrManifest
	}

	if strings.ContainsAny(parsedManifest.ApplicationManifests[0].Path, "*") && pta.LegacyPush == false {
		return pta, ErrNoWildcardSupport
	}

	pta.Manifest = parsedManifest
	pta.NoRouteManifestPath = noRouteManifestPath

	//check if a docker image shouldbe pushed and verify passed args combination
	if len(pta.DockerUserName) > 0 && (len(dockerPass) == 0 || len(pta.DockerImage) == 0) {
		return nil, ErrWrongPrivateDockerRepoCombination
	}

	//set timeout
	manifestTimeout, _ := strconv.Atoi(parsedManifest.ApplicationManifests[0].Timeout)
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
	if pta.LegacyPush && ((argPassed(flags, "health-check-type") && pta.HealthCheckType != "") || (argPassed(flags, "health-check-http-endpoint") && pta.HealthCheckHTTPEndpoint != "")) {
		return nil, ErrWrongCombination
	}

	// get health check settings from manifest if nothing else was specified in the command line
	if argPassed(flags, "health-check-type") == false {
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
		//convert variables to use them later in set-ent
		pta.Envs = env.Convert(envs)
	}

	//no-route set venerable-action to delete as default - but can be overwritten
	if (pta.NoRoute || pta.NoStart) && argPassed(flags, "venerable-action") == false {
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
