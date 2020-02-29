package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/cf"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"
	"github.com/happytobi/cf-puppeteer/rewind"
	"github.com/happytobi/cf-puppeteer/ui"
	"os"
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

func getActionsForApp(appRepo *ApplicationRepo, parsedArguments *arguments.ParserArguments) []rewind.Action {
	venName := venerableAppName(parsedArguments.AppName)
	puppeteerPush := cf.NewApplicationPush(appRepo.conn, appRepo.traceLogging)
	var err error
	var curApp *v2.AppResourcesEntity
	var venApp *v2.AppResourcesEntity

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
				// If there is no current app running, that's great, we're done here
				if curApp == nil {
					return nil
				}

				// If current app isn't started, then we'll just delete it, and we're done => only if route switching was not used
				if curApp.Entity.State != "STARTED" && parsedArguments.AddRoutes == false {
					return appRepo.v2Resources.DeleteApplication(parsedArguments.AppName)
				}

				// Do we have a ven app that will stop a rename? -> normal workflow only if we dont run the add routes mode
				if venApp != nil && parsedArguments.AddRoutes == false {
					// Finally, since the current app claims to be healthy, we'll delete the venerable app, and rename the current over the top
					err = appRepo.v2Resources.DeleteApplication(venName)
					if err != nil {
						return err
					}
				}

				if parsedArguments.AddRoutes == false {
					return appRepo.v2Resources.RenameApplication(parsedArguments.AppName, venName)
				}
				return nil
			},
		},
		// push
		{
			Forward: func() error {
				venAppExists := venApp != nil
				space, err := appRepo.conn.GetCurrentSpace()
				if err != nil {
					return err
				}
				if parsedArguments.AddRoutes == false {
					return puppeteerPush.PushApplication(venName, venAppExists, space.Guid, parsedArguments)
				}
				return nil
			},
			//When upload fails the new application will be deleted and ven app will be renamed
			ReversePrevious: func() error {
				ui.FailedMessage("error while uploading / deploying the application... roll everything back")
				_ = appRepo.v2Resources.DeleteApplication(parsedArguments.AppName)
				_ = appRepo.v2Resources.RenameApplication(venName, parsedArguments.AppName)
				return nil
			},
		},
		// start
		{
			Forward: func() error {
				if parsedArguments.NoStart == false {
					return appRepo.v2Resources.StartApplication(parsedArguments.AppName)
				}
				return nil
			},
			ReversePrevious: func() error {
				if parsedArguments.ShowCrashLogs {
					//print logs before application delete
					ui.Say("show crash logs")
					_ = appRepo.v2Resources.ShowCrashLogs(parsedArguments.AppName)
				}

				// If the app cannot start we'll have a lingering application
				// We delete this application so that the rename can succeed
				_ = appRepo.v2Resources.DeleteApplication(parsedArguments.AppName)
				return appRepo.v2Resources.RenameApplication(venName, parsedArguments.AppName)
			},
		},
		//switch routes because new application was started correct
		{
			Forward: func() error {
				//switch route only is application was started and route switch option was set
				ui.Say("check if routes should be added or switched from existing one")
				if parsedArguments.NoStart == false && parsedArguments.NoRoute == false {
					venAppExists := venApp != nil
					return puppeteerPush.SwitchRoutes(venName, venAppExists, parsedArguments.AppName, parsedArguments.Manifest.ApplicationManifests[0].Routes, parsedArguments.LegacyPush)
				}
				ui.Say("nothing to do")
				return nil
			},
			ReversePrevious: func() error {
				// If the app cannot start we'll have a lingering application
				// We delete this application so that the rename can succeed
				_ = appRepo.v2Resources.DeleteApplication(parsedArguments.AppName)
				return appRepo.v2Resources.RenameApplication(venName, parsedArguments.AppName)
			},
		},
		//check vor venerable application again -> because venerable action was set correct and ven app could exist now.
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
		// delete
		{
			Forward: func() error {
				//if venerableAction was set to stop
				if strings.ToLower(parsedArguments.VenerableAction) == "stop" && venApp != nil {
					return appRepo.v2Resources.StopApplication(venName)
				} else if strings.ToLower(parsedArguments.VenerableAction) == "delete" && venApp != nil {
					return appRepo.v2Resources.DeleteApplication(venName)
				}
				//do nothing with the ven app
				return nil
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
	if os.Getenv("CF_TRACE") == "true" {
		traceLogging = true
	}
	appRepo := NewApplicationRepo(cliConnection, traceLogging)
	parsedArguments, err := arguments.ParseArgs(args)
	fatalIf(err)

	fatalIf((&rewind.Actions{
		Actions:              getActionsForApp(appRepo, parsedArguments),
		RewindFailureMessage: "Oh no. Something's gone wrong. I've tried to roll back but you should check to see if everything is OK.",
	}).Execute())

	ui.Say("")
	ui.Say("A new version of your application has successfully been pushed!")
	ui.Say("")

	_ = appRepo.v2Resources.ListApplications()
}

// GetMetadata get plugin metadata
func (CfPuppeteerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cf-puppeteer",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 3,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "zero-downtime-push",
				HelpText: "Perform a zero-downtime push of an application over the top of an old one",
				UsageDetails: plugin.Usage{
					Usage:   "$ cf zero-downtime-push [<App-Name>] -f <Manifest.yml> [options]",
					Options: arguments.UsageDetailsOptionCommands(),
				},
			},
		},
	}
}

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
