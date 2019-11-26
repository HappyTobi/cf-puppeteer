package v3_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/happytobi/cf-puppeteer/arguments"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v3 "github.com/happytobi/cf-puppeteer/cf/v3"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfAppActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Resources")
}

var _ = Describe("cf-app test", func() {
	var (
		cliConn       *pluginfakes.FakeCliConnection
		resourcesData *v3.ResourcesData
		fakeExecutor  *cli.FakeExecutor
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		fakeExecutor = &cli.FakeExecutor{}
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, false), Executor: fakeExecutor.NewFakeExecutor()}
	})

	Describe("CreateApp v3", func() {
		It("push buildpack", func() {
			arguments := &arguments.ParserArguments{
				AppName: "myTestApp",
			}
			err := resourcesData.CreateApp(arguments)

			argumentsOutput := []string{"v3-create-app", "myTestApp"}
			Expect(fakeExecutor.ExecutorCallCount()).To(Equal(1))
			Expect(fakeExecutor.ExecutorArgumentsOutput()[0]).To(Equal(argumentsOutput))
			Expect(err).ToNot(HaveOccurred())
		})

		It("push docker image", func() {
			arguments := &arguments.ParserArguments{
				AppName:        "myTestApp",
				DockerImage:    "myDockerImage",
				DockerUserName: "mySecretDockerUser",
			}
			err := resourcesData.CreateApp(arguments)

			argumentsOutput := []string{"v3-create-app", "myTestApp", "--app-type", "docker"}
			Expect(fakeExecutor.ExecutorCallCount()).To(Equal(1))
			Expect(fakeExecutor.ExecutorArgumentsOutput()[0]).To(Equal(argumentsOutput))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("PushApp v3", func() {
		It("push buildpack with routes", func() {
			arguments := &arguments.ParserArguments{
				AppName: "myTestApp",
			}
			err := resourcesData.PushApp(arguments)

			argumentsOutput := []string{"v3-push", "myTestApp", "--no-start"}
			Expect(fakeExecutor.ExecutorCallCount()).To(Equal(1))
			Expect(fakeExecutor.ExecutorArgumentsOutput()[0]).To(Equal(argumentsOutput))
			Expect(err).ToNot(HaveOccurred())
		})

		It("push buildpack without routes", func() {
			arguments := &arguments.ParserArguments{
				AppName: "myTestApp",
				NoRoute: true,
			}
			err := resourcesData.PushApp(arguments)

			argumentsOutput := []string{"v3-push", "myTestApp", "--no-start", "--no-route"}
			Expect(fakeExecutor.ExecutorCallCount()).To(Equal(1))
			Expect(fakeExecutor.ExecutorArgumentsOutput()[0]).To(Equal(argumentsOutput))
			Expect(err).ToNot(HaveOccurred())
		})

		It("push buildpack with routes and envs", func() {
			envsMap := make(map[string]string, 2)
			envsMap["key"] = "value"
			envsMap["newKey"] = "newValue"

			arguments := &arguments.ParserArguments{
				AppName: "myTestApp",
				Envs:    envsMap,
			}
			err := resourcesData.PushApp(arguments)

			Expect(fakeExecutor.ExecutorCallCount()).To(Equal(3)) //called for each env
			Expect(err).ToNot(HaveOccurred())
		})

		It("push docker image", func() {
			arguments := &arguments.ParserArguments{
				AppName:        "myTestApp",
				DockerImage:    "myDockerImage",
				DockerUserName: "mySecretDockerUser",
			}
			err := resourcesData.PushApp(arguments)

			argumentsOutput := []string{"v3-push",
				"myTestApp",
				"--no-start",
				"--docker-image",
				"myDockerImage",
				"--docker-username",
				"mySecretDockerUser"}
			Expect(fakeExecutor.ExecutorCallCount()).To(Equal(1))
			Expect(fakeExecutor.ExecutorArgumentsOutput()[0]).To(Equal(argumentsOutput))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
