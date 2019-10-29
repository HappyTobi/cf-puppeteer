package cf

import (
	"encoding/json"
	"fmt"
	"github.com/happytobi/cf-puppeteer/manifest"
	"os"

	"github.com/happytobi/cf-puppeteer/arguments"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"

	v3 "github.com/happytobi/cf-puppeteer/cf/v3"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/blang/semver"
	"github.com/happytobi/cf-puppeteer/cf/cli"
)

//ApplicationPushData struct
type ApplicationPushData struct {
	Connection   plugin.CliConnection
	TraceLogging bool
}

//PuppeteerPush push application interface
type PuppeteerPush interface {
	PushApplication(venAppName string, spaceGUID string, manifest manifest.Application, parsedArguments *arguments.ParserArguments) error
}

var cliCalls cli.Calls

//NewApplicationPush generate new cf puppeteer push
func NewApplicationPush(conn plugin.CliConnection, traceLogging bool) *ApplicationPushData {
	cliCalls = cli.NewCli(conn, traceLogging)

	return &ApplicationPushData{
		Connection:   conn,
		TraceLogging: traceLogging,
	}
}

//PushApplication push application to cf
func (adp *ApplicationPushData) PushApplication(venAppName string, spaceGUID string, manifest manifest.Application, parsedArguments *arguments.ParserArguments) error {
	v3Push, err := useV3Push()
	if err != nil {
		//fatal exit
		os.Exit(1)
	}

	if v3Push && parsedArguments.LegacyPush != true {
		var v2Resources v2.Resources = v2.NewV2Resources(adp.Connection, adp.TraceLogging)
		var push v3.Push = v3.NewV3Push(adp.Connection, adp.TraceLogging)
		if parsedArguments.AddRoutes {
			//TODO loop over applications
			return push.SwitchRoutesOnly(venAppName, manifest.Name, manifest.Routes)
		}
		return push.PushApplication(venAppName, spaceGUID, manifest, parsedArguments, v2Resources)
	}

	var legacyPush v2.Push = v2.NewV2LegacyPush(adp.Connection, adp.TraceLogging)
	return legacyPush.PushApplication(venAppName, spaceGUID, parsedArguments)
}

func useV3Push() (bool, error) {
	_, v3ServerVersion, err := getCloudControllerAPIVersion()
	if err != nil {
		return false, err
	}
	v3SerSemVer, err := semver.Make(v3ServerVersion)
	if err != nil {
		return false, nil
	}

	expectedRange, err := semver.ParseRange(fmt.Sprintf(">=%s", v3.MinControllerVersion))
	//check if we can use the v3 push
	if expectedRange(v3SerSemVer) {
		return true, nil
	}
	return false, nil
}

// cloudControllerResponse
type cloudControllerResponse struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		BitsService       interface{} `json:"bits_service"`
		CloudControllerV2 struct {
			Href string `json:"href"`
			Meta struct {
				Version string `json:"version"`
			} `json:"meta"`
		} `json:"cloud_controller_v2"`
		CloudControllerV3 struct {
			Href string `json:"href"`
			Meta struct {
				Version string `json:"version"`
			} `json:"meta"`
		} `json:"cloud_controller_v3"`
		NetworkPolicyV0 struct {
			Href string `json:"href"`
		} `json:"network_policy_v0"`
		NetworkPolicyV1 struct {
			Href string `json:"href"`
		} `json:"network_policy_v1"`
		Uaa struct {
			Href string `json:"href"`
		} `json:"uaa"`
		Credhub interface{} `json:"credhub"`
		Routing struct {
			Href string `json:"href"`
		} `json:"routing"`
		Logging struct {
			Href string `json:"href"`
		} `json:"logging"`
		LogStream struct {
			Href string `json:"href"`
		} `json:"log_stream"`
		AppSSH struct {
			Href string `json:"href"`
			Meta struct {
				HostKeyFingerprint string `json:"host_key_fingerprint"`
				OauthClient        string `json:"oauth_client"`
			} `json:"meta"`
		} `json:"app_ssh"`
	} `json:"links"`
}

func getCloudControllerAPIVersion() (string, string, error) {
	callResp, err := cliCalls.GetJSON("/")
	var response cloudControllerResponse
	err = json.Unmarshal([]byte(callResp), &response)
	if err != nil {
		return "", "", err
	}
	return response.Links.CloudControllerV2.Meta.Version, response.Links.CloudControllerV3.Meta.Version, nil
}
