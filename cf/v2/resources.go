package v2

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/happytobi/cf-puppeteer/cf/cli"
)

type Resources interface {
	GetAppMetadata(appName string) (*AppResourcesEntity, error)
	RenameApplication(oldName string, newName string) (err error)
	StopApplication(appName string) (err error)
	StartApplication(appName string) (err error)
	DeleteApplication(appName string) (err error)
	ShowCrashLogs(appName string) (err error)
	ListApplications() (err error)
}

//ResourcesData internal struct with connection an tracing options etc
type ResourcesData struct {
	cli        cli.Calls
	connection plugin.CliConnection
}

//NewV2Resources constructor
func NewV2Resources(conn plugin.CliConnection, traceLogging bool) *ResourcesData {
	return &ResourcesData{
		cli:        cli.NewCli(conn, traceLogging),
		connection: conn,
	}
}
