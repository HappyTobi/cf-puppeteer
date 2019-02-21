package main

import (
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
	"github.com/happytobi/cf-puppeteer/manifest"
	"github.com/happytobi/cf-puppeteer/rewind"
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

func getActionsForApp(appRepo *ApplicationRepo, appName, manifestPath, appPath, stackName string, timeout int, vars []string, varsFiles []string, envs []string, showLogs bool) []rewind.Action {
	venName := venerableAppName(appName)
	var err error
	var curApp, venApp *AppEntity
	var haveVenToCleanup bool

	return []rewind.Action{
		// get info about current app
		{
			Forward: func() error {
				curApp, err = appRepo.GetAppMetadata(appName)
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
				if curApp.State != "STARTED" {
					return appRepo.DeleteApplication(appName)
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
				return appRepo.RenameApplication(appName, venName)
			},
		},
		// push
		{
			Forward: func() error {
				return appRepo.PushApplication(appName, manifestPath, appPath, stackName, timeout, vars, varsFiles, envs, showLogs)
			},
			ReversePrevious: func() error {
				if !haveVenToCleanup {
					return nil
				}

				// If the app cannot start we'll have a lingering application
				// We delete this application so that the rename can succeed
				appRepo.DeleteApplication(appName)

				return appRepo.RenameApplication(venName, appName)
			},
		},
		// delete
		{
			Forward: func() error {
				if !haveVenToCleanup {
					return nil
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

	appRepo := NewApplicationRepo(cliConnection)
	appName, manifestPath, appPath, timeout, stackName, vars, varsFiles, envs, showLogs, err := ParseArgs(args)
	fatalIf(err)

	fatalIf((&rewind.Actions{
		Actions:              getActionsForApp(appRepo, appName, manifestPath, appPath, stackName, timeout, vars, varsFiles, envs, showLogs),
		RewindFailureMessage: "Oh no. Something's gone wrong. I've tried to roll back but you should check to see if everything is OK.",
	}).Execute())

	fmt.Println()
	fmt.Println("A new version of your application has successfully been pushed!")
	fmt.Println()

	_ = appRepo.ListApplications()
}

func (CfPuppeteerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cf-puppeteer",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 12,
		},
		Commands: []plugin.Command{
			{
				Name:     "zero-downtime-push",
				HelpText: "Perform a zero-downtime push of an application over the top of an old one",
				UsageDetails: plugin.Usage{
					Usage: "$ cf zero-downtime-push application-to-replace \\ \n \t-f path/to/new_manifest.yml \\ \n \t-p path/to/new/path",
				},
			},
		},
	}
}

type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprint(*s)
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

//ParseArgs parse all cmd arguments
func ParseArgs(args []string) (string, string, string, int, string, []string, []string, []string, bool, error) {
	flags := flag.NewFlagSet("zero-downtime-push", flag.ContinueOnError)

	var envs StringSlice
	var vars StringSlice
	var varsFiles StringSlice

	manifestPath := flags.String("f", "", "path to an application manifest")
	appPath := flags.String("p", "", "path to application files")
	stackName := flags.String("s", "", "name of the stack to use")
	timeout := flags.Int("t", 0, "push timout in secounds defualt 60s")
	showLogs := flags.Bool("show-app-log", false, "tail and show application log during application start")
	flags.Var(&envs, "env", "Variable key value pair for adding dynamic environment variables; can specity multiple times")
	flags.Var(&vars, "var", "Variable key value pair for variable substitution, (e.g., name=app1); can specify multiple times")
	flags.Var(&varsFiles, "vars-file", "Path to a variable substitution file for manifest; can specify multiple times")

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
		return "", "", "", *timeout, "", []string{}, []string{}, []string{}, false, err
	}

	if *manifestPath == "" {
		return "", "", "", *timeout, "", []string{}, []string{}, []string{}, false, ErrNoManifest
	}

	//parse manifest
	manifest, err := manifest.ParseManifest(*manifestPath)
	if err != nil {
		return "", "", "", *timeout, "", []string{}, []string{}, []string{}, false, ErrManifest
	}

	//set timeout
	manifestTimeout := manifest.ApplicationManifests[0].Timeout
	if manifestTimeout > 0 && *timeout <= 0 {
		*timeout = manifestTimeout
	} else if manifestTimeout <= 0 && *timeout <= 0 {
		*timeout = 60
	}

	//parse first argument as appName
	appName := args[1]
	if noAppNameProvided {
		appName = manifest.ApplicationManifests[0].Name

	}

	//validate envs format
	if len(envs) > 0 {
		for _, envPair := range envs {
			if strings.Contains(envPair, "=") == false {
				return "", "", "", *timeout, "", []string{}, []string{}, []string{}, false, ErrWrongEnvFormat
			}
		}
	}

	return appName, *manifestPath, *appPath, *timeout, *stackName, vars, varsFiles, envs, *showLogs, nil
}

var (
	ErrNoManifest     = errors.New("a manifest is required to push this application")
	ErrManifest       = errors.New("could not parse manifest")
	ErrWrongEnvFormat = errors.New("--var would be in wrong format, use the vars like key=value")
)

type ApplicationRepo struct {
	conn plugin.CliConnection
}

func NewApplicationRepo(conn plugin.CliConnection) *ApplicationRepo {
	return &ApplicationRepo{
		conn: conn,
	}
}

func (repo *ApplicationRepo) RenameApplication(oldName, newName string) error {
	_, err := repo.conn.CliCommand("rename", oldName, newName)
	return err
}

func (repo *ApplicationRepo) PushApplication(appName, manifestPath, appPath, stackName string, timeout int, vars []string, varsFiles []string, envs []string, showLogs bool) error {
	args := []string{"push", appName, "-f", manifestPath, "--no-start"}

	if appPath != "" {
		args = append(args, "-p", appPath)
	}

	if stackName != "" {
		args = append(args, "-s", stackName)
	}

	/* always append timout */
	timeoutS := strconv.Itoa(timeout)
	args = append(args, "-t", timeoutS)

	for _, varPair := range vars {
		args = append(args, "--var", varPair)
	}

	for _, varsFile := range varsFiles {
		args = append(args, "--vars-file", varsFile)
	}

	_, err := repo.conn.CliCommand(args...)
	if err != nil {
		return err
	}

	envErr := repo.SetEnvironmentVariables(appName, envs)
	if envErr != nil {
		return envErr
	}

	if showLogs {
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

	_, err = repo.conn.CliCommand("start", appName)
	if err != nil {
		return err
	}

	return nil
}

//SetEnvironmentVariable set passed envs with set-env to set variables dynamically
func (repo *ApplicationRepo) SetEnvironmentVariables(appName string, envs []string) error {
	varArgs := []string{"set-env", appName}
	//set all variables passed by --var
	for _, envPair := range envs {
		tmpArgs := make([]string, len(varArgs))
		copy(tmpArgs, varArgs)
		newArgs := strings.Split(envPair, "=")
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

type AppEntity struct {
	State string `json:"state"`
}

var (
	ErrAppNotFound = errors.New("application not found")
)

// GetAppMetadata returns metadata about an app with appName
func (repo *ApplicationRepo) GetAppMetadata(appName string) (*AppEntity, error) {
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

	output := struct {
		Resources []struct {
			Entity AppEntity `json:"entity"`
		} `json:"resources"`
	}{}
	err = json.Unmarshal([]byte(jsonResp), &output)

	if err != nil {
		return nil, err
	}

	if len(output.Resources) == 0 {
		return nil, ErrAppNotFound
	}

	return &output.Resources[0].Entity, nil
}
