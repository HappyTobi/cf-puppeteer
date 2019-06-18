package v2

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	"github.com/happytobi/cf-puppeteer/ui"
	"strconv"
	"strings"
)

//Push interface with all v3 actions
type Push interface {
	PushApplication(venAppName string, spaceGUID string, parsedArguments *arguments.ParserArguments) error
}

//ResourcesData internal struct with connection an tracing options etc
type LegacyResourcesData struct {
	Executor cli.Executor
}

//NewV2LegacyPush constructor
func NewV2LegacyPush(conn plugin.CliConnection, traceLogging bool) *LegacyResourcesData {
	return &LegacyResourcesData{
		Executor: cli.NewExecutor(traceLogging),
	}
}

func (resource *LegacyResourcesData) PushApplication(venAppName, spaceGUID string, parsedArguments *arguments.ParserArguments) error {
	ui.Say("use legacy push")
	args := []string{"push", parsedArguments.AppName, "-f", parsedArguments.ManifestPath, "--no-start"}
	if parsedArguments.AppPath != "" {
		args = append(args, "-p", parsedArguments.AppPath)
	}

	if parsedArguments.StackName != "" {
		args = append(args, "-s", parsedArguments.StackName)
	}

	if parsedArguments.InvocationTimeout >= 0 {
		timeoutS := strconv.Itoa(parsedArguments.InvocationTimeout)
		args = append(args, "-t", timeoutS)
	}

	ui.Say("start pushing application with arguments %s", args)
	err := resource.Executor.Execute(args)
	if err != nil {
		return err
	}

	//set all environment variables
	err = resource.setEnvironmentVariables(parsedArguments)
	if err != nil {
		return err
	}

	return nil
}

func (resource *LegacyResourcesData) setEnvironmentVariables(parsedArguments *arguments.ParserArguments) (err error) {
	ui.Say("set passed environment variables")
	varArgs := []string{"set-env", parsedArguments.AppName}
	//set all variables passed by --var
	for _, envPair := range parsedArguments.Envs {
		tmpArgs := make([]string, len(varArgs))
		copy(tmpArgs, varArgs)
		newArgs := strings.Split(envPair, "=")
		tmpArgs = append(tmpArgs, newArgs...)
		err := resource.Executor.Execute(tmpArgs)
		if err != nil {
			return err
		}
	}
	ui.Ok()
	return nil
}
