package v3_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v3 "github.com/happytobi/cf-puppeteer/cf/v3"
	"testing"

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
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, true)}

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

			var routes = []map[string]string{0: {"key": "route", "value": "url.test-domain.com"}, 1: {"key": "route", "value": "foo.example.com"}}
			domainResponse, err := resourcesData.GetDomain(routes)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			Expect((*domainResponse)[0].Host).To(Equal("foo"))
			Expect((*domainResponse)[0].DomainGUID).To(Equal("3a5d3d89-3f89-4f05-8188-8a2b298c79d7"))

			Expect((*domainResponse)[1].Host).To(Equal("url"))
			Expect((*domainResponse)[1].DomainGUID).To(Equal("3a5d3d89-3f89-4f05-8188-8a2b298c79d5"))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
