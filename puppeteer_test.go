package main_test

import (
	"errors"
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/happytobi/cf-puppeteer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
})
