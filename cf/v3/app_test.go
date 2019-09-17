package v3_test

/*
import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
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
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, true)}

	})

	Describe("PushApplication with curl v3 api", func() {
		It("use v3 push api", func() {
			response := []string{`{
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
				}`}
			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			var envVars = []string{"foo=bar", "x=y"}
			var buildPacks = []string{"java_buildpack"}
			pushResponse, err := resourcesData.PushApp("my_app", "5EB2DB8D-2808-4871-9AF2-0A0977881A4B", buildPacks, "cflinuxfs3", envVars)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			Expect(pushResponse.GUID).To(Equal("1cb006ee-fb05-47e1-b541-c34179ddc446"))
			Expect(err).ToNot(HaveOccurred())
		})
		//TODO parse more data and check them
		It("use v3 push api without a buildpack", func() {
			response := []string{`{
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
				}`}
			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			var envVars = []string{"foo=bar", "x=y"}
			var buildPacks = []string{""}
			pushResponse, err := resourcesData.PushApp("my_app", "5EB2DB8D-2808-4871-9AF2-0A0977881A4B", buildPacks, "cflinuxfs3", envVars)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			Expect(pushResponse.GUID).To(Equal("1cb006ee-fb05-47e1-b541-c34179ddc446"))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})*/
