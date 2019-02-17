package main_test

import (
	"errors"
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/happytobi/cf-puppeteer"

	plugin_models "code.cloudfoundry.org/cli/plugin/models"
)

func TestPuppeteer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer Suite")
}

var _ = Describe("Flag Parsing", func() {
	It("parses args without appName", func() {
		appName, manifestPath, appPath, timeout, stackName, vars, varsFiles, showLogs, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"-f", "./fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
				"-var", "foo=bar",
				"-var", "baz=bob",
				"-vars-file", "vars.yml",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(appName).Should(Equal("myApp"))
		Expect(manifestPath).To(Equal("./fixtures/manifest.yml"))
		Expect(appPath).To(Equal("app-path"))
		Expect(stackName).To(Equal("stack-name"))
		Expect(vars).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(varsFiles).To(Equal([]string{"vars.yml"}))
		Expect(showLogs).To(Equal(false))
		Expect(timeout).To(Equal(120))
	})

	It("parses a all args without timeout", func() {
		appName, manifestPath, appPath, timeout, stackName, vars, varsFiles, showLogs, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-f", "manifest-path",
				"-p", "app-path",
				"-s", "stack-name",
				"-var", "foo=bar",
				"-var", "baz=bob",
				"-vars-file", "vars.yml",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(appName).To(Equal("appname"))
		Expect(manifestPath).To(Equal("manifest-path"))
		Expect(appPath).To(Equal("app-path"))
		Expect(stackName).To(Equal("stack-name"))
		Expect(vars).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(varsFiles).To(Equal([]string{"vars.yml"}))
		Expect(showLogs).To(Equal(false))
		Expect(timeout).To(Equal(60))
	})

	It("parses a complete set of args", func() {
		appName, manifestPath, appPath, timeout, stackName, vars, varsFiles, showLogs, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-f", "manifest-path",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
				"-var", "foo=bar",
				"-var", "baz=bob",
				"-vars-file", "vars.yml",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(appName).To(Equal("appname"))
		Expect(manifestPath).To(Equal("manifest-path"))
		Expect(appPath).To(Equal("app-path"))
		Expect(stackName).To(Equal("stack-name"))
		Expect(vars).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(varsFiles).To(Equal([]string{"vars.yml"}))
		Expect(showLogs).To(Equal(false))
		Expect(timeout).To(Equal(120))
	})

	It("requires a manifest", func() {
		_, _, _, _, _, _, _, _, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-p", "app-path",
			},
		)
		Expect(err).To(MatchError(ErrNoManifest))
	})
})

var _ = Describe("ApplicationRepo", func() {
	var (
		cliConn *pluginfakes.FakeCliConnection
		repo    *ApplicationRepo
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		repo = NewApplicationRepo(cliConn)
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
				`{"resources":[{"entity":{"state":"STARTED"}}]}`,
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
		})

		It("URL encodes the application name", func() {
			response := []string{
				`{"resources":[{"entity":{"state":"STARTED"}}]}`,
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
			err := repo.PushApplication("appName", "/path/to/a/manifest.yml", "/path/to/the/app", "", 60, []string{}, []string{}, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(2))
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
			err := repo.PushApplication("appName", "/path/to/a/manifest.yml", "", "", 60, []string{}, []string{}, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(2))
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
			err := repo.PushApplication("appName", "/path/to/a/manifest.yml", "/path/to/the/app", "stackName", 60, []string{}, []string{}, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(2))
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
			err := repo.PushApplication("appName", "/path/to/a/manifest.yml", "", "", 60, []string{"foo=bar", "baz=bob"}, []string{"vars.yml"}, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(cliConn.CliCommandCallCount()).To(Equal(2))
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
		})

		It("returns errors from the push", func() {
			cliConn.CliCommandReturns([]string{}, errors.New("bad app"))

			err := repo.PushApplication("appName", "/path/to/a/manifest.yml", "/path/to/the/app", "", 60, []string{}, []string{}, false)
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
