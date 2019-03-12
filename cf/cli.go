package cli

import "code.cloudfoundry.org/cli/plugin"

type CFInstance struct {
	conn plugin.CliConnection
}

function NewCFCliInsance(conn plugin.CliConnection) *CFInstance {
	return &CFInstance{
		conn: conn,
	}
}
