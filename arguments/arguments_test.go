package arguments

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArgumentsParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arguments Parser TestSuite")
}

var _ = Describe("Flag Parsing", func() {
	var ()

	BeforeEach(func() {

	})

	It("parses args without appName", func() {
		parsedArguments, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"-f", "../fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
				"--env", "foo=bar",
				"--env", "baz=bob=true",
				"--vendor-option", "stop",
				"--invocation-timeout", "2211",
				"--process", "process-name",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).Should(Equal("myApp"))
		Expect(parsedArguments.ManifestPath).To(Equal("../fixtures/manifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob=true"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("stop"))
		Expect(parsedArguments.ShowLogs).To(Equal(false))
		Expect(parsedArguments.ShowCrashLogs).To(Equal(false))
		Expect(parsedArguments.Timeout).To(Equal(120))
		Expect(parsedArguments.InvocationTimeout).To(Equal(2211))
		Expect(parsedArguments.Process).To(Equal("process-name"))
		Expect(parsedArguments.HealthCheckType).To(Equal("http"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/health"))
		Expect(parsedArguments.NoRoute).To(Equal(false))
	})

	It("parses a all args without timeout", func() {
		parsedArguments, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-f", "../fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-env", "foo=bar",
				"-env", "baz=bob",
				"--vendor-option", "stop",
				"--no-routes",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).To(Equal("appname"))
		Expect(parsedArguments.ManifestPath).To(Equal("../fixtures/manifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("stop"))
		Expect(parsedArguments.ShowLogs).To(Equal(false))
		Expect(parsedArguments.ShowCrashLogs).To(Equal(false))
		Expect(parsedArguments.NoRoute).To(Equal(true))
		Expect(parsedArguments.Timeout).To(Equal(2))
		Expect(parsedArguments.InvocationTimeout).To(Equal(-1))
		Expect(parsedArguments.Process).To(Equal(""))
		Expect(parsedArguments.HealthCheckType).To(Equal("http"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/health"))
	})

	It("parses a all args without timeout and no manifest timeout", func() {
		parsedArguments, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-f", "../fixtures/multiManifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-env", "foo=bar",
				"-env", "baz=bob",
				"--vendor-option", "stop",
				"--show-crash-log",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).To(Equal("appname"))
		Expect(parsedArguments.ManifestPath).To(Equal("../fixtures/multiManifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("stop"))
		Expect(parsedArguments.ShowLogs).To(Equal(false))
		Expect(parsedArguments.ShowCrashLogs).To(Equal(true))
		Expect(parsedArguments.Timeout).To(Equal(60))
		Expect(parsedArguments.InvocationTimeout).To(Equal(-1))
		Expect(parsedArguments.Process).To(Equal(""))
		Expect(parsedArguments.HealthCheckType).To(Equal("http"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/health"))
	})

	It("parses a complete set of args", func() {
		parsedArguments, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-f", "../fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
				"-env", "foo=bar",
				"-env", "baz=bob",
				"--invocation-timeout", "2211",
				"--process", "process-name",
				"--health-check-type", "process",
				"--health-check-http-endpoint", "/foo/bar",
				"--show-app-log",
				"--add-routes",
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(parsedArguments.AppName).To(Equal("appname"))
		Expect(parsedArguments.ManifestPath).To(Equal("../fixtures/manifest.yml"))
		Expect(parsedArguments.AppPath).To(Equal("app-path"))
		Expect(parsedArguments.StackName).To(Equal("stack-name"))
		Expect(parsedArguments.Envs).To(Equal([]string{"foo=bar", "baz=bob"}))
		Expect(parsedArguments.VendorAppOption).Should(Equal("delete"))
		Expect(parsedArguments.ShowLogs).To(Equal(true))
		Expect(parsedArguments.ShowCrashLogs).To(Equal(false))
		Expect(parsedArguments.Timeout).To(Equal(120))
		Expect(parsedArguments.InvocationTimeout).To(Equal(2211))
		Expect(parsedArguments.Process).To(Equal("process-name"))
		Expect(parsedArguments.HealthCheckType).To(Equal("process"))
		Expect(parsedArguments.HealthCheckHTTPEndpoint).To(Equal("/foo/bar"))
		Expect(parsedArguments.AddRoutes).To(Equal(true))
	})

	It("parses args without appName and wrong envs format", func() {
		_, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"-f", "../fixtures/manifest.yml",
				"-p", "app-path",
				"-s", "stack-name",
				"-t", "120",
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
			[]string{
				"zero-downtime-push",
				"appname",
				"-p", "app-path",
			},
		)
		Expect(err).To(MatchError(ErrNoManifest))
	})

	It("legacy push with health check option", func() {
		_, err := ParseArgs(
			[]string{
				"zero-downtime-push",
				"appname",
				"-f", "../fixtures/manifest.yml",
				"-p", "app-path",
				"--legacy-push",
				"--health-check-type", "process",
				"--health-check-http-endpoint", "/foo/bar",
			},
		)
		Expect(err).To(MatchError(ErrWrongCombination))
	})
})
