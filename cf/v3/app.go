package v3

import (
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/ui"
)

/* TODO add type to docker */
func (resource *ResourcesData) CreateApp(parsedArguments *arguments.ParserArguments) (err error) {
	args := []string{"v3-create-app", parsedArguments.AppName}
	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}
	return nil
}

func (resource *ResourcesData) PushApp(parsedArguments *arguments.ParserArguments) (err error) {
	args := []string{"v3-push", parsedArguments.AppName, "--no-start"}
	if parsedArguments.AppPath != "" {
		args = append(args, "-p", parsedArguments.AppPath)
	}

	if parsedArguments.NoRoute == true {
		args = append(args, "--no-route")
	}

	err = resource.Executor.Execute(args)
	if err != nil {
		return err
	}

	for _, env := range parsedArguments.Envs {
		ui.Say("set environment-variable")
		args = []string{"v3-set-env", parsedArguments.AppName, env}
		err = resource.Executor.Execute(args)
		if err != nil {
			ui.Failed("could not set environment variable %s to application %s", env, parsedArguments.AppName)
		}
	}

	return nil
}
