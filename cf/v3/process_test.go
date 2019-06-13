package v3_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"fmt"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v3 "github.com/happytobi/cf-puppeteer/cf/v3"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfProcessActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Resources")
}

var _ = Describe("cf-process test", func() {
	var (
		cliConn       *pluginfakes.FakeCliConnection
		resourcesData *v3.ResourcesData
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, true)}

	})
	Describe("Fetch process with curl v3 api", func() {
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

			applicationEntity := v3.ApplicationEntity{}
			applicationEntity.Command = command
			applicationEntity.HealthCheck.Data.Endpoint = "/health"
			applicationEntity.HealthCheck.Data.InvocationTimeout = 60
			applicationEntity.HealthCheck.HealthCheckType = "http"

			err := resourcesData.UpdateApplicationProcessWebInformation(appGUID, applicationEntity)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			args := cliConn.CliCommandWithoutTerminalOutputArgsForCall(0)
			fmt.Printf("%v", args)
			Expect(args).To(Equal([]string{"curl", "/v3/processes/999", "-X", "PATCH", "-H", "Content-type: application/json", "-d", "{\"command\":\"JAVA_OPTS=FOOBAR\",\"health_check\":{\"data\":{\"endpoint\":\"/health\",\"invocation_timeout\":60},\"type\":\"http\"}}"}))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
