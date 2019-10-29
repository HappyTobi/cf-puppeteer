package v2

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	"github.com/happytobi/cf-puppeteer/ui"
	"github.com/pkg/errors"
	"strconv"
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

	if parsedArguments.NoStart {
		args = append(args, "--no-start")
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
	for envKey, envVal := range parsedArguments.Envs {
		tmpArgs := make([]string, len(parsedArguments.Envs))
		copy(tmpArgs, varArgs)
		tmpArgs = append(tmpArgs, fmt.Sprintf("%s %s", envKey, envVal))
		err := resource.Executor.Execute(tmpArgs)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not set env-variable with key %s to application %s", envKey, parsedArguments.AppName))
		}
	}
	ui.Ok()
	return nil
}
