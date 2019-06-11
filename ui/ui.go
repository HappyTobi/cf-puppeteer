package ui

import (
	"code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
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

//Failed message see cf/terminal
func Failed(message string, args ...interface{}) {
	ui.Failed(message, args...)
}

//LoadingIndication message see cf/terminal
func LoadingIndication() {
	ui.LoadingIndication()
}
