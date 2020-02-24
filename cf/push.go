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
	SwitchRoutes(venAppName string, venAppExists bool, appName string, routes []map[string]string, legacyPush bool) error
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
		return legacyPush.PushApplication(parsedArguments)
	}
	//v3 push
	var v2Resources v2.Resources = v2.NewV2Resources(adp.Connection, adp.TraceLogging)
	var push v3.Push = v3.NewV3Push(adp.Connection, adp.TraceLogging)
	return push.PushApplication(venAppName, spaceGUID, parsedArguments, v2Resources)
}

//handle route switch
func (adp *ApplicationPushData) SwitchRoutes(venAppName string, venAppExists bool, appName string, routes []map[string]string, legacyPush bool) error {
	if legacyPush {
		var legacyPush v2.Push = v2.NewV2LegacyPush(adp.Connection, adp.TraceLogging)
		return legacyPush.SwitchRoutesOnly(venAppName, venAppExists, appName, routes)
	}
	var push v3.Push = v3.NewV3Push(adp.Connection, adp.TraceLogging)
	return push.SwitchRoutesOnly(venAppName, appName, routes)
}
