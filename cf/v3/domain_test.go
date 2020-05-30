package v3_test

import (
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v3 "github.com/happytobi/cf-puppeteer/cf/v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfDomainActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Resources")
}

var _ = Describe("cf-domain test", func() {
	var (
		cliConn       *pluginfakes.FakeCliConnection
		resourcesData *v3.ResourcesData
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, false)}

	})
	Describe("Domain actions with curl v3 api", func() {
		It("use v3 domain api", func() {
			response := []string{`{
				"pagination": {
				  "total_results": 3,
				  "total_pages": 2,
				  "first": {
					"href": "https://api.example.org/v3/domains?page=1&per_page=2"
				  },
				  "last": {
					"href": "https://api.example.org/v3/domains?page=2&per_page=2"
				  },
				  "next": {
					"href": "https://api.example.org/v3/domains?page=2&per_page=2"
				  },
				  "previous": null
				},
				"resources": [
				  {
					"guid": "3a5d3d89-3f89-4f05-8188-8a2b298c79d5",
					"created_at": "2019-03-08T01:06:19Z",
					"updated_at": "2019-03-08T01:06:19Z",
					"name": "test-domain.com",
					"internal": false,
					"metadata": {
					  "labels": {},
					  "annotations": {}
					},
					"relationships": {
					  "organization": {
						"data": null
					  },
					  "shared_organizations": {
						"data": []
					  }
					},
					"links": {
					  "self": {
						"href": "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5"
					  }
					}
				  },
					{
					"guid": "3a5d3d89-3f89-4f05-8188-8a2b298c79d7",
					"created_at": "2019-03-08T01:06:19Z",
					"updated_at": "2019-03-08T01:06:19Z",
					"name": "example.com",
					"internal": false,
					"metadata": {
					  "labels": {},
					  "annotations": {}
					},
					"relationships": {
					  "organization": {
						"data": null
					  },
					  "shared_organizations": {
						"data": []
					  }
					},
					"links": {
					  "self": {
						"href": "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5"
					  }
					}
				  }
				]
			  }
			`}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			var routes = []map[string]string{0: {"key": "route", "value": "url.test-domain.com"}, 1: {"key": "route", "value": "foo.example.com"}, 2: {"key": "route", "value": "boo.example.com/api"}}
			domainResponse, err := resourcesData.GetDomain(routes)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			checkMap := make(map[string]string, len(*domainResponse))
			checkPath := make(map[string]string, len(*domainResponse))
			for _, value := range *domainResponse {
				checkMap[value.Host] = value.Domain
				checkPath[value.Host] = value.Path
			}

			Expect(checkMap["foo"]).To(Equal("example.com"))
			Expect(checkMap["url"]).To(Equal("test-domain.com"))
			Expect(checkMap["boo"]).To(Equal("example.com"))

			Expect(checkPath["foo"]).To(Equal(""))
			Expect(checkPath["url"]).To(Equal(""))
			Expect(checkPath["boo"]).To(Equal("api"))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
