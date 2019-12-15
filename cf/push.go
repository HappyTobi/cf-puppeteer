package cf

import (
	"github.com/happytobi/cf-puppeteer/arguments"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"

	v3 "github.com/happytobi/cf-puppeteer/cf/v3"

	"code.cloudfoundry.org/cli/plugin"
)

//ApplicationPushData struct
type ApplicationPushData struct {
	Connection   plugin.CliConnection
	TraceLogging bool
}

//PuppeteerPush push application interface
type PuppeteerPush interface {
	PushApplication(venAppName string, venAppExists bool, spaceGUID string, parsedArguments *arguments.ParserArguments) error
}

//NewApplicationPush generate new cf puppeteer push
func NewApplicationPush(conn plugin.CliConnection, traceLogging bool) *ApplicationPushData {
	return &ApplicationPushData{
		Connection:   conn,
		TraceLogging: traceLogging,
	}
}

//PushApplication push application to cf
func (adp *ApplicationPushData) PushApplication(venAppName string, venAppExists bool, spaceGUID string, parsedArguments *arguments.ParserArguments) error {
	if parsedArguments.LegacyPush == true {
		var legacyPush v2.Push = v2.NewV2LegacyPush(adp.Connection, adp.TraceLogging)
		if parsedArguments.AddRoutes {
			return legacyPush.SwitchRoutesOnly(venAppName, venAppExists, parsedArguments.AppName, parsedArguments.Manifest.ApplicationManifests[0].Routes)
		}
		return legacyPush.PushApplication(venAppName, spaceGUID, parsedArguments)
	}
	//v3 push
	var v2Resources v2.Resources = v2.NewV2Resources(adp.Connection, adp.TraceLogging)
	var push v3.Push = v3.NewV3Push(adp.Connection, adp.TraceLogging)
	if parsedArguments.AddRoutes {
		//TODO loop over applications
		return push.SwitchRoutesOnly(venAppName, parsedArguments.AppName, parsedArguments.Manifest.ApplicationManifests[0].Routes)
	}
	return push.PushApplication(venAppName, spaceGUID, parsedArguments, v2Resources)
}
