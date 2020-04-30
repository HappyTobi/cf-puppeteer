package v2_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	v2 "github.com/happytobi/cf-puppeteer/cf/v2"
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
		resourcesData *v2.LegacyResourcesData
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		resourcesData = v2.NewV2LegacyPush(cliConn, false)
	})
	Describe("Domain actions with curl v2 api", func() {
		It("use v2 domain api", func() {
			response := []string{`{
			   "total_results": 18,
			   "total_pages": 2,
			   "prev_url": null,
			   "next_url": "https://api.example.org/v2/domains?page=1&per_page=2",
			   "resources": [
				  {
					 "metadata": {
						"guid": "e97d1675-894e-4808-8b58-82d1805x7368",
						"url": "/v2/shared_domains/e97d1675-894e-4808-8b58-82d1805x7368",
						"created_at": "2016-10-20T07:40:24Z",
						"updated_at": "2019-02-12T10:33:49Z"
					 },
					 "entity": {
						"name": "test-domain.com",
						"internal": false,
						"router_group_guid": null,
						"router_group_type": null
					 }
				  },
				  {
					 "metadata": {
						"guid": "aa23b15e-dc54-437e-a651-a29415b66d1h",
						"url": "/v2/shared_domains/aa23b15e-dc54-437e-a651-a29415b66d1h",
						"created_at": "2016-10-20T07:40:24Z",
						"updated_at": "2018-02-09T06:19:46Z"
					 },
					 "entity": {
						"name": "example.com",
						"internal": false,
						"router_group_guid": null,
						"router_group_type": null
					 }
				  },
 				  {
					 "metadata": {
						"guid": "aa23b15e-dc54-437e-a651-a29415b66d1m",
						"url": "/v2/shared_domains/aa23b15e-dc54-437e-a651-a29415b66d1m",
						"created_at": "2016-10-20T07:40:24Z",
						"updated_at": "2018-02-09T06:19:46Z"
					 },
					 "entity": {
						"name": "foo.example.com",
						"internal": false,
						"router_group_guid": null,
						"router_group_type": null
					 }
				  },
				  {
					 "metadata": {
						"guid": "aa23b15e-dc54-437e-a651-a29415b66d9m",
						"url": "/v2/shared_domains/aa23b15e-dc54-437e-a651-a29415b66d9m",
						"created_at": "2016-10-20T07:40:24Z",
						"updated_at": "2018-02-09T06:19:46Z"
					 },
					 "entity": {
						"name": "internal.emea.github.com",
						"internal": false,
						"router_group_guid": null,
						"router_group_type": null
					 }
				  }
				]}
			`}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			var routes = []map[string]string{0: {"route": "url.test-domain.com"}, 1: {"route": "foo.example.com"}, 2: {"route": "my.foo.example.com"}, 3: {"route": "puppeteer.internal.emea.github.com"}}
			domainResponse, err := resourcesData.GetDomain(routes)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(4).To(Equal(len(*domainResponse)))

			checkMap := make(map[string]string, len(*domainResponse))
			for _, value := range *domainResponse {
				checkMap[value.Host] = value.Domain
			}

			Expect(checkMap["my"]).To(Equal("foo.example.com"))
			Expect(checkMap["puppeteer"]).To(Equal("internal.emea.github.com"))
			Expect(checkMap["url"]).To(Equal("test-domain.com"))

			Expect(err).ToNot(HaveOccurred())
		})

		It("use v2 domain api with multiple subdomains ", func() {
			response := []string{`{
			   "total_results": 18,
			   "total_pages": 2,
			   "prev_url": null,
			   "next_url": "https://api.example.org/v2/domains?page=1&per_page=2",
			   "resources": [
				  {
					 "metadata": {
						"guid": "aa23b15e-dc54-437e-a651-a29415b66d9m",
						"url": "/v2/shared_domains/aa23b15e-dc54-437e-a651-a29415b66d9m",
						"created_at": "2016-10-20T07:40:24Z",
						"updated_at": "2018-02-09T06:19:46Z"
					 },
					 "entity": {
						"name": "staging.product.com",
						"internal": false,
						"router_group_guid": null,
						"router_group_type": null
					 }
				  },
 				  {
					 "metadata": {
						"guid": "aa23b15e-dc54-437e-a651-5673241243",
						"url": "/v2/shared_domains/aa23b15e-dc54-437e-a651-5673241243",
						"created_at": "2016-10-20T07:40:24Z",
						"updated_at": "2018-02-09T06:19:46Z"
					 },
					 "entity": {
						"name": "cfapp.io",
						"internal": false,
						"router_group_guid": null,
						"router_group_type": null
					 }
				  }
				]}
			`}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			var routes = []map[string]string{0: {"route": " staging-p.cfapp.io"}, 1: {"route": "staging-product.cfapp.io"}, 2: {"route": "staging.product.com"}}
			domainResponse, err := resourcesData.GetDomain(routes)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(3).To(Equal(len(*domainResponse)))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
