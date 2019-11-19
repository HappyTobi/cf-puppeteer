package v3

import (
	"fmt"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/ui"
	"github.com/pkg/errors"
)

func (resource *ResourcesData) CreateApp(parsedArguments *arguments.ParserArguments) (err error) {
	args := []string{"v3-create-app", parsedArguments.AppName}

	//check if user pushes a docker image
	if len(parsedArguments.DockerImage) > 0 {
		args = append(args, "--app-type", "docker")
	}

	err = resource.Executor.Execute(args)
	if err != nil {
		return errors.Wrap(err, "could not create app")
	}
	return nil
}

func (resource *ResourcesData) PushApp(parsedArguments *arguments.ParserArguments) (err error) {
	args := []string{"v3-push", parsedArguments.AppName, "--no-start"}
	if len(parsedArguments.AppPath) > 0 && len(parsedArguments.DockerImage) == 0 {
		args = append(args, "-p", parsedArguments.AppPath)
	}

	if len(parsedArguments.DockerImage) > 0 {
		args = append(args, "--docker-image", parsedArguments.DockerImage)
		if len(parsedArguments.DockerUserName) > 0 {
			args = append(args, "--docker-username", parsedArguments.DockerUserName)
		}
	}

	if parsedArguments.NoRoute == true {
		args = append(args, "--no-route")
	}

	err = resource.Executor.Execute(args)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not push application with passed args: %v", args))
	}

	for envKey, envVal := range parsedArguments.Envs {
		ui.Say(fmt.Sprintf("set environment-variable %s", envKey))
		args = []string{"v3-set-env", parsedArguments.AppName, envKey, envVal}
		err = resource.Executor.Execute(args)
		if err != nil {
			ui.Failed("could not set environment variable with key: %s to application: %s", envKey, parsedArguments.AppName)
		}
	}

	return nil
}
