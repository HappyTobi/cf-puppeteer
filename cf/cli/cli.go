package cli

import (
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	print "github.com/happytobi/cf-puppeteer/cf/utils"
	"github.com/happytobi/cf-puppeteer/ui"
)

//Calls cli curl calls
type Calls interface {
	GetJSON(path string) (string, error)
	PostJSON(path string, jsonBody string) (string, error)
	PatchJSON(path string, jsonBody string) (string, error)
}

//Connection cli connection object
type Connection struct {
	cf           plugin.CliConnection
	traceLogging bool
}

//NewCli ff
func NewCli(conn plugin.CliConnection, traceLogging bool) *Connection {
	return &Connection{
		cf:           conn,
		traceLogging: traceLogging,
	}
}

//GetJSON make an get call to an url
func (conn *Connection) GetJSON(path string) (string, error) {
	result, err := conn.cf.CliCommandWithoutTerminalOutput("curl", path, "-X", "GET", "-H", "Content-type: application/json")
	if err != nil {
		ui.Failed("error while calling %s - error: %s", path, err)
		return "", err
	}

	jsonResp := strings.Join(result, "")

	if conn.traceLogging {
		ui.Say("response from GET call - path: %s was: %s %s", path, jsonResp, print.PrettyJSON(jsonResp))
	}

	return jsonResp, nil
}

//PostJSON post to path with json body
func (conn *Connection) PostJSON(path string, jsonBody string) (string, error) {
	args := []string{"curl", path, "-X", "POST", "-H", "Content-type: application/json"}
	if jsonBody != "" {
		bodyArgs := []string{"-d", jsonBody}
		args = append(args, bodyArgs...)
	}

	result, err := conn.cf.CliCommandWithoutTerminalOutput(args...)
	if err != nil {
		return "", err
	}

	jsonResp := strings.Join(result, "")
	if conn.traceLogging {
		ui.Say("response from POST call - path: %s was: %s", path, print.PrettyJSON(jsonResp))
	}
	return jsonResp, nil
}

//PatchJSON post to path with json body
func (conn *Connection) PatchJSON(path string, jsonBody string) (string, error) {
	args := []string{"curl", path, "-X", "PATCH", "-H", "Content-type: application/json", "-d", jsonBody}

	result, err := conn.cf.CliCommandWithoutTerminalOutput(args...)
	if err != nil {
		return "", err
	}
	jsonResp := strings.Join(result, "")
	if conn.traceLogging {
		ui.Say("response from PATCH call - path: %s was: %s", path, print.PrettyJSON(jsonResp))
	}
	return jsonResp, nil
}
