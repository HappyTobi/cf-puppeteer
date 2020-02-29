package arguments

import "fmt"

//Commands list with all commands
func commands() (commands map[string]string) {
	return map[string]string{
		"f":                          "path to an application manifest",
		"p":                          "path to application files",
		"s":                          "name of the stack to use",
		"health-check-type":          "type of health check to perform",
		"health-check-http-endpoint": "endpoint for the 'http' health check type",
		"t":                          "push timeout in seconds (default 60 seconds",
		"invocation-timeout":         "health check invocation timeout in seconds",
		"process":                    "use health check type process",
		"show-crash-log":             "Show recent logs when applications crashes while the deployment",
		"venerable-action":           "option to delete,stop,none application action on venerable app default is delete",
		"env":                        "Variable key value pair for adding dynamic environment variables; can specify multiple times",
		"legacy-push":                "use legacy push instead of new v3 api",
		"no-route":                   "deploy new application without adding routes",
		"route-only":                 "only add routes from manifest to the application",
		"no-start":                   "don't start application after deployment; venerable action is none",
		"vars-file":                  "path to a variable substitution file for manifest",
		"docker-image":               "docker image url",
		"docker-username":            "docker repository username; used with password from env CF_DOCKER_PASSWORD",
	}
}

//UsageDetailsOptionCommands change the args for the options to support posix style
func UsageDetailsOptionCommands() (optionCommands map[string]string) {
	rangeCommands := commands()
	optionCommands = make(map[string]string, len(rangeCommands))
	for k, v := range rangeCommands {
		arg := k
		if len(k) > 1 {
			arg = fmt.Sprintf("-%s", arg)
		}
		optionCommands[arg] = v
	}
	return optionCommands
}
