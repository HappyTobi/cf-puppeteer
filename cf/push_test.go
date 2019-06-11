package cf

import (
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer Version Check Test")
}

var _ = Describe("cf push version check", func() {
	var (
		cliConn *pluginfakes.FakeCliConnection
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		//set cliCalls internal used var
		cliCalls = cli.NewCli(cliConn, false)
	})

	Describe("check controller call", func() {
		It("should return version", func() {
			response := []string{
				`{
					"links": {
					   "self": {
						  "href": "https://puppeteer-fake-url.com"
					   },
					   "bits_service": null,
					   "cloud_controller_v2": {
						  "href": "https://puppeteer-fake-url.com/v2",
						  "meta": {
							 "version": "2.134.0"
						  }
					   },
					   "cloud_controller_v3": {
						  "href": "https://puppeteer-fake-url.com/v3",
						  "meta": {
							 "version": "3.69.0"
						  }
					   },
					   "network_policy_v0": {
						  "href": "https://puppeteer-fake-url.com/networking/v0/external"
					   },
					   "network_policy_v1": {
						  "href": "https://puppeteer-fake-url.com/networking/v1/external"
					   },
					   "uaa": {
						  "href": "https://uaa.puppeteer-fake-url.com"
					   },
					   "credhub": null,
					   "routing": {
						  "href": "https://puppeteer-fake-url.com/routing"
					   },
					   "logging": {
						  "href": "wss://doppler.puppeteer-fake-url.com:443"
					   },
					   "log_stream": {
						  "href": "https://log-stream.puppeteer-fake-url.com"
					   }
					}
				 }`,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			v2Ver, v3Ver, err := getCloudControllerAPIVersion()
			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(v3Ver).To(Equal("3.69.0"))
			Expect(v2Ver).To(Equal("2.134.0"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("check useV3Version", func() {
			response := []string{
				`{
					"links": {
					   "self": {
						  "href": "https://puppeteer-fake-url.com"
					   },
					   "bits_service": null,
					   "cloud_controller_v2": {
						  "href": "https://puppeteer-fake-url.com/v2",
						  "meta": {
							 "version": "2.134.0"
						  }
					   },
					   "cloud_controller_v3": {
						  "href": "https://puppeteer-fake-url.com/v3",
						  "meta": {
							 "version": "3.69.0"
						  }
					   },
					   "network_policy_v0": {
						  "href": "https://puppeteer-fake-url.com/networking/v0/external"
					   },
					   "network_policy_v1": {
						  "href": "https://puppeteer-fake-url.com/networking/v1/external"
					   },
					   "uaa": {
						  "href": "https://uaa.puppeteer-fake-url.com"
					   },
					   "credhub": null,
					   "routing": {
						  "href": "https://puppeteer-fake-url.com/routing"
					   },
					   "logging": {
						  "href": "wss://doppler.puppeteer-fake-url.com:443"
					   },
					   "log_stream": {
						  "href": "https://log-stream.puppeteer-fake-url.com"
					   }
					}
				 }`,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			useV3, err := useV3Push()
			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(useV3).To(Equal(true))
			Expect(err).ToNot(HaveOccurred())
		})

		It("check useV3Version - false", func() {
			response := []string{
				`{
					"links": {
					   "self": {
						  "href": "https://puppeteer-fake-url.com"
					   },
					   "bits_service": null,
					   "cloud_controller_v2": {
						  "href": "https://puppeteer-fake-url.com/v2",
						  "meta": {
							 "version": "2.134.0"
						  }
					   },
					   "cloud_controller_v3": {
						  "href": "https://puppeteer-fake-url.com/v3",
						  "meta": {
							 "version": "3.10.0"
						  }
					   },
					   "network_policy_v0": {
						  "href": "https://puppeteer-fake-url.com/networking/v0/external"
					   },
					   "network_policy_v1": {
						  "href": "https://puppeteer-fake-url.com/networking/v1/external"
					   },
					   "uaa": {
						  "href": "https://uaa.puppeteer-fake-url.com"
					   },
					   "credhub": null,
					   "routing": {
						  "href": "https://puppeteer-fake-url.com/routing"
					   },
					   "logging": {
						  "href": "wss://doppler.puppeteer-fake-url.com:443"
					   },
					   "log_stream": {
						  "href": "https://log-stream.puppeteer-fake-url.com"
					   }
					}
				 }`,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			useV3, err := useV3Push()
			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(useV3).To(Equal(false))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
