package ui

import (
	"code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"fmt"
	"os"
)

var ui terminal.UI

func init() {
	i18n.T = func(translationID string, args ...interface{}) string {
		return translationID
	}

	ui = terminal.NewUI(
		os.Stdin,
		os.Stdout,
		terminal.NewTeePrinter(os.Stdout),
		trace.NewLogger(os.Stdout, false, "", ""))
}

//Say message see cf/terminal
func Say(message string, args ...interface{}) {
	ui.Say(message, args...)
}

//Ok message see cf/terminal
func Ok() {
	ui.Ok()
}

//InfoMessage print out message in green / ok colored mode
func InfoMessage(message string) {
	ui.Say(terminal.SuccessColor(message))
}

//Failed message see cf/terminal
func Failed(message string, args ...interface{}) {
	ui.Failed(message, args...)
}

//FailedMessage print out message in error color without the "FAILED" message
func FailedMessage(message string) {
	ui.Say(terminal.FailureColor(message))
}

//Warn message see cf/terminal
func Warn(message string, args ...interface{}) {
	ui.Warn(message, args...)
}

func DebugMessage(message string, args ...interface{}) {
	traceEnv := os.Getenv("CF_TRACE")
	if traceEnv == "true" || (traceEnv != "false" && len(traceEnv) > 0) {
		//check env for CF_TRACE
		message = fmt.Sprintf(message, args...)
		ui.Say(terminal.AdvisoryColor(message))
	}
}

//LoadingIndication message see cf/terminal
func LoadingIndication() {
	ui.LoadingIndication()
}
