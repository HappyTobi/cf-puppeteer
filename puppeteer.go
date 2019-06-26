package main

import (
	"code.cloudfoundry.org/cli/cf/api/logs"
	"code.cloudfoundry.org/cli/plugin"
	"context"
	"fmt"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/cf"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"
	"github.com/happytobi/cf-puppeteer/rewind"
	"github.com/happytobi/cf-puppeteer/ui"
	"log"
	"os"
	"strings"
	"time"
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
	var err error
	var curApp, venApp *v2.AppResourcesEntity

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

				// If current app isn't started, then we'll just delete it, and we're done
				if curApp.Entity.State != "STARTED" {
					return appRepo.DeleteApplication(parsedArguments.AppName)
				}

				// Do we have a ven app that will stop a rename? -> normal workflow only if we dont run the add routes mode
				if venApp != nil && parsedArguments.AddRoutes == false {
					// Finally, since the current app claims to be healthy, we'll delete the venerable app, and rename the current over the top
					err = appRepo.DeleteApplication(venName)
					if err != nil {
						return err
					}
				}

				if parsedArguments.AddRoutes == false {
					return appRepo.RenameApplication(parsedArguments.AppName, venName)
				}
				return nil
			},
		},
		// push
		{
			Forward: func() error {
				space, err := appRepo.conn.GetCurrentSpace()
				if err != nil {
					return err
				}

				var puppeteerPush cf.PuppeteerPush = cf.NewApplicationPush(appRepo.conn, appRepo.traceLogging)
				return puppeteerPush.PushApplication(venName, space.Guid, parsedArguments)
			},
			//When upload fails the new application will be deleted and ven app will be renamed
			ReversePrevious: func() error {
				ui.Failed("error while uploading / deploying the application... roll everything back")
				_ = appRepo.DeleteApplication(parsedArguments.AppName)
				return appRepo.RenameApplication(venName, parsedArguments.AppName)
			},
		},
		// start
		{
			Forward: func() error {
				ui.Say("show logs...")
				if parsedArguments.ShowLogs {
					// TODO not working anymore
					_ = appRepo.ShowLogs(parsedArguments.AppName)
				}
				return appRepo.StartApplication(parsedArguments.AppName)
			},
			ReversePrevious: func() error {
				if parsedArguments.ShowCrashLogs {
					//print logs before application delete
					ui.Say("show crash logs")
					_ = appRepo.ShowCrashLogs(parsedArguments.AppName)
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
				//if vendorAppOption was set to stop
				if strings.ToLower(parsedArguments.VendorAppOption) == "stop" {
					return appRepo.StopApplication(venName)
				} else if strings.ToLower(parsedArguments.VendorAppOption) == "delete" {
					return appRepo.DeleteApplication(venName)
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
	if os.Getenv("CF_PUPPETEER_TRACE") == "true" {
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

	_ = appRepo.ListApplications()
}

// GetMetadata get plugin metadata
func (CfPuppeteerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cf-puppeteer",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 1,
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
						"t":                           "push timeout (in seconds)",
						"-env":                        "add environment key value pairs dynamic; can specify multiple times",
						"-vendor-option":              "option to delete or stop vendor application - default is delete",
						"-health-check-type":          "type of health check to perform",
						"-health-check-http-endpoint": "endpoint for the 'http' health check type",
						"-invocation-timeout":         "timeout (in seconds) that controls individual health check invocations",
						"-show-crash-log":             "Show recent logs when applications crashes while the deployment",
						//"-show-app-log": "tail and show application log during application start",
						"-process":     "application process to update",
						"-legacy-push": "use legacy push instead of new v3 api",
						"-no-routes":   "deploy new application without adding routes",
					},
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

func (repo *ApplicationRepo) ShowCrashLogs(appName string) error {
	_, err := repo.conn.CliCommand("logs", "--recent", appName)
	return err
}

func (repo *ApplicationRepo) ListApplications() error {
	_, err := repo.conn.CliCommand("apps")
	return err
}

func (repo *ApplicationRepo) ShowLogs(appName string) error {
	app, err := repo.conn.GetApp(appName)
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

	messages, chanError := cons.TailingLogs(app.Guid, token)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case m := <-messages:
				if m.GetSourceType() != "STG" { // skip STG messages as the cf tool already prints them
					os.Stderr.WriteString(logs.NewNoaaLogMessage(m).ToLog(time.Local) + "\n")
				}
			case e := <-chanError:
				log.Println("error reading logs:", e)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
