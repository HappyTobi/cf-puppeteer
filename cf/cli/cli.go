package cli

import (
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/happytobi/cf-puppeteer/ui"
)

type Calls interface {
	GetJSON(path string) (string, error)
}

type CliConnection struct {
	conn         plugin.CliConnection
	traceLogging bool
}

//NewCli ff
func NewCli(conn plugin.CliConnection, traceLogging bool) *CliConnection {
	return &CliConnection{
		conn:         conn,
		traceLogging: traceLogging,
	}
}

//GetJSONCall make an get call to an url
func (conn *CliConnection) GetJSON(path string) (string, error) {
	result, err := conn.conn.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")
	if err != nil {
		return "", err
	}
	jsonResp := strings.Join(result, "")
	if conn.traceLogging {
		ui.Say("response from get call to path: %s was: %s\n", path /*prettyPrintJSON(jsonResp) */, jsonResp)
	}
	return jsonResp, nil
}
