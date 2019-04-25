package cfResources_test

import (
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/happytobi/cf-puppeteer/cfResources"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Resources")
}

var _ = Describe("cfResources", func() {
	var (
		cliConn *pluginfakes.FakeCliConnection
		cf      cfResources.CfResourcesInterface
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		cf = cfResources.NewResources(cliConn, false)
	})
	Describe("PushApplication with curl v3 api", func() {
		It("push simple application with v3 api", func() {
			response := []string{
				`{
					"guid": "1cb006ee-fb05-47e1-b541-c34179ddc446",
					"name": "my_app",
					"state": "STOPPED",
					"created_at": "2016-03-17T21:41:30Z",
					"updated_at": "2016-06-08T16:41:26Z",
					"lifecycle": {
					  "type": "buildpack",
					  "data": {
						"buildpacks": ["java_buildpack"],
						"stack": "cflinuxfs3"
					  }
					},
					"relationships": {
					  "space": {
						"data": {
						  "guid": "2f35885d-0c9d-4423-83ad-fd05066f8576"
						}
					  }
					},
					"links": {
					  "self": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446"
					  },
					  "space": {
						"href": "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576"
					  },
					  "processes": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/processes"
					  },
					  "route_mappings": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/route_mappings"
					  },
					  "packages": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/packages"
					  },
					  "environment_variables": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables"
					  },
					  "current_droplet": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets/current"
					  },
					  "droplets": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets"
					  },
					  "tasks": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/tasks"
					  },
					  "start": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start",
						"method": "POST"
					  },
					  "stop": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/stop",
						"method": "POST"
					  },
					  "revisions": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/revisions"
					  },
					  "deployed_revisions": {
						"href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/revisions/deployed"
					  }
					},
					"metadata": {
					  "labels": {},
					  "annotations": {}
					}
				  }
				  `,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			pushResponse, err := cf.PushApp("my_app", "2f35885d-0c9d-4423-83ad-fd05066f8576")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(pushResponse.Guid).To(Equal("1cb006ee-fb05-47e1-b541-c34179ddc446"))
			Expect(err).ToNot(HaveOccurred())

		})
	})
})
