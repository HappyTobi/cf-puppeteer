package main_test

import (
	"errors"
	"fmt"
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/happytobi/cf-puppeteer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	plugin_models "code.cloudfoundry.org/cli/plugin/models"
)

func TestPuppeteer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer Suite")
}

var _ = Describe("Flag Parsing", func() {
	var (
		cliConn *pluginfakes.FakeCliConnection
		repo    *ApplicationRepo
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		repo = NewApplicationRepo(cliConn, false)
	})

	It("parses args without appName", func() {
		parsedArguments, err := ParseArgs(
			repo, []string{
				"zero-downtime-push",
				"-f", "./fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
				"-var", "foo=bar",
				"-var", "baz=bob",
				"-vars-file", "vars.yml",
				"-env", "foo=bar",
				"-env", "baz=bob=true",
				"--vendor-option", "stop",
				"--invocation-timeout", "2211",
				"--process", "process-name",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).Should(Equal("myApp"))
		Expect(parsedArguments.ManifestPath).To(Equal("./fixtures/manifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Vars).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VarsFiles).To(Equal([]string{"vars.yml"}))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob=true"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("stop"))
		Expect(parsedArguments.ShowLogs).To(Equal(false))
		Expect(parsedArguments.Timeout).To(Equal(120))
		Expect(parsedArguments.InvocationTimeout).To(Equal(2211))
		Expect(parsedArguments.Process).To(Equal("process-name"))
		Expect(parsedArguments.HealthCheckType).To(Equal("http"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/health"))
	})

	It("parses a all args without timeout", func() {
		parsedArguments, err := ParseArgs(
			repo, []string{
				"zero-downtime-push",
				"appname",
				"-f", "./fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-var", "foo=bar",
				"-var", "baz=bob",
				"-vars-file", "vars.yml",
				"-env", "foo=bar",
				"-env", "baz=bob",
				"--vendor-option", "stop",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).To(Equal("appname"))
		Expect(parsedArguments.ManifestPath).To(Equal("./fixtures/manifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Vars).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VarsFiles).To(Equal([]string{"vars.yml"}))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("stop"))
		Expect(parsedArguments.ShowLogs).To(Equal(false))
		Expect(parsedArguments.Timeout).To(Equal(2))
		Expect(parsedArguments.InvocationTimeout).To(Equal(-1))
		Expect(parsedArguments.Process).To(Equal(""))
		Expect(parsedArguments.HealthCheckType).To(Equal("http"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/health"))
	})

	It("parses a all args without timeout and no manifest timeout", func() {
		parsedArguments, err := ParseArgs(
			repo, []string{
				"zero-downtime-push",
				"appname",
				"-f", "./fixtures/multiManifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-var", "foo=bar",
				"-var", "baz=bob",
				"-vars-file", "vars.yml",
				"-env", "foo=bar",
				"-env", "baz=bob",
				"--vendor-option", "stop",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).To(Equal("appname"))
		Expect(parsedArguments.ManifestPath).To(Equal("./fixtures/multiManifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Vars).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VarsFiles).To(Equal([]string{"vars.yml"}))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("stop"))
		Expect(parsedArguments.ShowLogs).To(Equal(false))
		Expect(parsedArguments.Timeout).To(Equal(60))
		Expect(parsedArguments.InvocationTimeout).To(Equal(-1))
		Expect(parsedArguments.Process).To(Equal(""))
		Expect(parsedArguments.HealthCheckType).To(Equal("http"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/health"))
	})

	It("parses a complete set of args", func() {
		parsedArguments, err := ParseArgs(
			repo, []string{
				"zero-downtime-push",
				"appname",
				"-f", "./fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
				"-var", "foo=bar",
				"-var", "baz=bob",
				"-vars-file", "vars.yml",
				"-env", "foo=bar",
				"-env", "baz=bob",
				"--invocation-timeout", "2211",
				"--process", "process-name",
				"--health-check-type", "process",
				"--health-check-http-endpoint", "/foo/bar",
				"--show-app-log",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).To(Equal("appname"))
		Expect(parsedArguments.ManifestPath).To(Equal("./fixtures/manifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Vars).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VarsFiles).To(Equal([]string{"vars.yml"}))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("delete"))
		Expect(parsedArguments.ShowLogs).To(Equal(true))
		Expect(parsedArguments.Timeout).To(Equal(120))
		Expect(parsedArguments.InvocationTimeout).To(Equal(2211))
		Expect(parsedArguments.Process).To(Equal("process-name"))
		Expect(parsedArguments.HealthCheckType).To(Equal("process"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/foo/bar"))
	})

	It("parses args without appName and wrong envs format", func() {
		_, err := ParseArgs(
			repo, []string{
				"zero-downtime-push",
				"-f", "./fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
				"-var", "foo=bar",
				"-var", "baz bob",
				"-vars-file", "vars.yml",
				"-env", "foo=bar",
				"-env", "baz bob",
				"--invocation-timeout", "2211",
				"--process", "process-name",
			},
		)
		Expect(err).To(MatchError(ErrWrongEnvFormat))
	})

	It("requires a manifest", func() {
		_, err := ParseArgs(
			repo, []string{
				"zero-downtime-push",
				"appname",
				"-p", "app-path",
			},
		)
		Expect(err).To(MatchError(ErrNoManifest))
	})
})

var _ = Describe("CheckAllV3Commands", func() {
	var (
		cliConn *pluginfakes.FakeCliConnection
		repo    *ApplicationRepo
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		repo = NewApplicationRepo(cliConn, false)
	})

	Describe("checkAPIV3", func() {
		It("available CfsV3Api", func() {
			response := []string{
				`{}`,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			err := repo.CheckAPIV3()

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args).To(Equal([]string{"curl", "/v3", "-X", "GET"}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("not available CfV3Api", func() {
			response := []string{
				`{
                    "description": "Unknown request",
                    "error_code": "CF-NotFound",
                    "code": 10000
                 }`,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			err := repo.CheckAPIV3()

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args).To(Equal([]string{"curl", "/v3", "-X", "GET"}))

			Expect(err).To(MatchError("cf api v3 is not available"))
		})

		It("check application process web informations", func() {
			response := []string{
				`{
                    "guid": "999",
                    "type": "web",
                    "command": "helloWorld=comman",
                    "instances": 1,
                    "memory_in_mb": 128
                    }`,
			}

			appGUID := "999"

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			applicationProcessesEntityV3, err := repo.GetApplicationProcessWebInformation(appGUID)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args).To(Equal([]string{"curl", "/v3/apps/999/processes/web", "-X", "GET"}))

			Expect(applicationProcessesEntityV3).ToNot(BeNil())
			Expect(applicationProcessesEntityV3.GUID).To(Equal("999"))
			Expect(applicationProcessesEntityV3.Command).To(Equal("helloWorld=comman"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("check application process web informations - app not available", func() {
			response := []string{
				`{
                    "errors": []
                    }`,
			}

			appGUID := "999"

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			_, err := repo.GetApplicationProcessWebInformation(appGUID)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args).To(Equal([]string{"curl", "/v3/apps/999/processes/web", "-X", "GET"}))

			Expect(err).To(MatchError("application not found"))
			Expect(err).To(HaveOccurred())
		})

		It("update application with invocation timeout setting", func() {
			response := []string{
				`{
                    "command": "JAVA_OPTS=FOOBAR",
                    "created_at": "2019-02-25T14:09:01Z",
                    "disk_in_mb": 1024,
                    "guid": "6ca30711-72d2-415b-8ed3-6870b7e56741",
                    "health_check": {
                        "data": {
                            "endpoint": "/health",
                            "invocation_timeout": 60
                        },
                        "type": "http"
                    }
                }`,
			}

			appGUID := "999"
			command := "JAVA_OPTS=FOOBAR"

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			applicationEntity := ApplicationEntityV3{}
			applicationEntity.Command = command
			applicationEntity.HealthCheck.Data.Endpoint = "/health"
			applicationEntity.HealthCheck.Data.InvocationTimeout = 60
			applicationEntity.HealthCheck.HealthCheckType = "http"

			err := repo.UpdateApplicationProcessWebInformation(appGUID, applicationEntity)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			fmt.Printf("%v", args)
			Expect(args).To(Equal([]string{"curl", "/v3/processes/999", "-X", "PATCH", "-H", "Content-type: application/json", "-d", "{\"command\":\"JAVA_OPTS=FOOBAR\",\"health_check\":{\"data\":{\"endpoint\":\"/health\",\"invocation_timeout\":60},\"type\":\"http\"}}"}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("update application with process but without command setting", func() {
			response := []string{
				`{
                    "created_at": "2019-02-25T14:09:01Z",
                    "disk_in_mb": 1024,
                    "guid": "6ca30711-72d2-415b-8ed3-6870b7e56741",
                    "health_check": {
                        "type": "process"
                    },
                    "type": "web"
                }`,
			}

			appGUID := "999"

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			applicationEntity := ApplicationEntityV3{}
			applicationEntity.HealthCheck.HealthCheckType = "process"
			applicationEntity.ProcessType = "web"

			err := repo.UpdateApplicationProcessWebInformation(appGUID, applicationEntity)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args).To(Equal([]string{"curl", "/v3/processes/999", "-X", "PATCH", "-H", "Content-type: application/json", "-d", "{\"health_check\":{\"data\":{},\"type\":\"process\"},\"type\":\"web\"}"}))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})

var _ = Describe("ApplicationRepo", func() {
	var (
		cliConn *pluginfakes.FakeCliConnection
		repo    *ApplicationRepo
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		repo = NewApplicationRepo(cliConn, false)
	})

	Describe("RenameApplication", func() {
		It("renames the application", func() {
			err := repo.RenameApplication("old-name", "new-name")
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Expect(args).To(Equal([]string{"rename", "old-name", "new-name"}))
		})

		It("returns an error if one occurs", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("no app"))

			err := repo.RenameApplication("old-name", "new-name")
			Expect(err).To(MatchError("no app"))
		})
	})

	Describe("GetAppMetadata", func() {

		It("returns an error if the cli returns an error", func() {
			cliConn.CliCommandWithoutTerminalOutputReturns([]string{}, errors.New("you shall not curl"))
			_, err := repo.GetAppMetadata("app-name")

			Expect(err).To(MatchError("you shall not curl"))
		})

		It("returns an error if the cli response is invalid JSON", func() {
			response := []string{
				"}notjson{",
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			_, err := repo.GetAppMetadata("app-name")

			Expect(err).To(HaveOccurred())
		})

		It("returns app data if the app exists", func() {
			response := []string{
				`{"resources":[
                    {
                        "metadata": {
                            "guid": "6ca30711-72d2-415b-8ed3-6870b7e56741"
                         },
                        "entity":
                            {
                                "state":"STARTED"
                            }
                    }]
                }`,
			}
			spaceGUID := "4"

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			cliConn.GetCurrentSpaceReturns(
				plugin_models.Space{
					SpaceFields: plugin_models.SpaceFields{
						Guid: spaceGUID,
					},
				},
				nil,
			)

			result, err := repo.GetAppMetadata("app-name")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args).To(Equal([]string{"curl", "v2/apps?q=name:app-name&q=space_guid:4"}))

			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
			Expect(result.Metadata.GUID).To(Equal("6ca30711-72d2-415b-8ed3-6870b7e56741"))
		})

		It("URL encodes the application name", func() {
			response := []string{
				`{"resources":[
                    {
                        "metadata": {
                            "guid": "6ca30711-72d2-415b-8ed3-6870b7e56741"
                         },
                        "entity":
                            {
                                "state":"STARTED"
                            }
                    }]
                }`,
			}
			spaceGUID := "4"

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			cliConn.GetCurrentSpaceReturns(
				plugin_models.Space{
					SpaceFields: plugin_models.SpaceFields{
						Guid: spaceGUID,
					},
				},
				nil,
			)

			result, err := repo.GetAppMetadata("app name")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args).To(Equal([]string{"curl", "v2/apps?q=name:app+name&q=space_guid:4"}))

			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
		})

		It("returns nil if the app does not exist", func() {
			response := []string{
				`{"resources":[]}`,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			result, err := repo.GetAppMetadata("app-name")

			Expect(err).To(Equal(ErrAppNotFound))
			Expect(result).To(BeNil())
		})

	})

	Describe("PushApplication", func() {
		It("pushes an application with both a manifest and a path", func() {
			parsedArguments := &ParserArguments{
				AppName:      "appName",
				ManifestPath: "/path/to/a/manifest.yml",
				AppPath:      "/path/to/the/app",
				Timeout:      60,
			}

			err := repo.PushApplication(parsedArguments)
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Expect(args).To(Equal([]string{
				"push",
				"appName",
				"-f", "/path/to/a/manifest.yml",
				"--no-start",
				"-p", "/path/to/the/app",
				"-t", "60",
			}))
		})

		It("pushes an application with only a manifest", func() {
			parsedArguments := &ParserArguments{
				AppName:      "appName",
				ManifestPath: "/path/to/a/manifest.yml",

				Timeout: 60,
			}

			err := repo.PushApplication(parsedArguments)
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Expect(args).To(Equal([]string{
				"push",
				"appName",
				"-f", "/path/to/a/manifest.yml",
				"--no-start",
				"-t", "60",
			}))
		})

		It("pushes an application with a stack", func() {
			parsedArguments := &ParserArguments{
				AppName:      "appName",
				ManifestPath: "/path/to/a/manifest.yml",
				AppPath:      "/path/to/the/app",
				StackName:    "stackName",
				Timeout:      60,
			}

			err := repo.PushApplication(parsedArguments)

			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Expect(args).To(Equal([]string{
				"push",
				"appName",
				"-f", "/path/to/a/manifest.yml",
				"--no-start",
				"-p", "/path/to/the/app",
				"-s", "stackName",
				"-t", "60",
			}))
		})

		It("pushes an application with variables", func() {
			parsedArguments := &ParserArguments{
				AppName:      "appName",
				ManifestPath: "/path/to/a/manifest.yml",
				Vars:         []string{"foo=bar", "baz=bob"},
				VarsFiles:    []string{"vars.yml"},
				Envs:         []string{"foo=bar", "bar=foo=true"},
				Timeout:      60,
			}

			err := repo.PushApplication(parsedArguments)
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(3))
			args := cliConn.CliCommandArgsForCall(0)
			Expect(args).To(Equal([]string{
				"push",
				"appName",
				"-f", "/path/to/a/manifest.yml",
				"--no-start",
				"-t", "60",
				"--var", "foo=bar",
				"--var", "baz=bob",
				"--vars-file", "vars.yml",
			}))

			args2 := cliConn.CliCommandArgsForCall(1)
			Expect(args2).To(Equal([]string{
				"set-env",
				"appName",
				"foo",
				"bar",
			}))

			args3 := cliConn.CliCommandArgsForCall(2)
			Expect(args3).To(Equal([]string{
				"set-env",
				"appName",
				"bar",
				"foo=true",
			}))
		})

		It("returns errors from the push", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("bad app"))

			parsedArguments := &ParserArguments{
				AppName:      "appName",
				ManifestPath: "/path/to/a/manifest.yml",
				AppPath:      "/path/to/the/app",
			}

			err := repo.PushApplication(parsedArguments)

			Expect(err).To(MatchError("bad app"))
		})
	})

	Describe("DeleteApplication", func() {
		It("deletes all trace of an application", func() {
			err := repo.DeleteApplication("app-name")
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Expect(args).To(Equal([]string{
				"delete", "app-name",
				"-f",
			}))
		})

		It("returns errors from the delete", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("bad app"))

			err := repo.DeleteApplication("app-name")
			Expect(err).To(MatchError("bad app"))
		})
	})

	Describe("ListApplications", func() {
		It("lists all the applications", func() {
			err := repo.ListApplications()
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(1))
			args := cliConn.CliCommandArgsForCall(0)
			Expect(args).To(Equal([]string{"apps"}))
		})

		It("returns errors from the list", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("bad apps"))

			err := repo.ListApplications()
			Expect(err).To(MatchError("bad apps"))
		})
	})
})
