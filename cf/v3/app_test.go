package v3_test

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

	Describe("GetApplication with curl v3 api", func() {
		It("use v3 get app api", func() {
			response := []string{`{
			  "guid": "5EB2DB8D-2808-4871-9AF2-0A0977881A4B",
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
			getAppResponse, err := resourcesData.GetApp("5EB2DB8D-2808-4871-9AF2-0A0977881A4B")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			Expect(getAppResponse.GUID).To(Equal("5EB2DB8D-2808-4871-9AF2-0A0977881A4B"))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetRoutesApp with curl v3 api", func() {
		It("use v3 api to get routes", func() {
			response := []string{`{
			  "guid": "5EB2DB8D-2808-4871-9AF2-0A0977881A4B",
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
			err := resourcesData.StartApp("5EB2DB8D-2808-4871-9AF2-0A0977881A4B")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("AssignApp droplet to application with curl v3 api", func() {
		It("use v3 assignApp api", func() {
			response := []string{`{
			  "data": {
				"guid": "F6E3E504-99C2-41BD-B78D-262F3E70A1F8"
			  },
			  "links": {
				"self": {
				  "href": "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/relationships/current_droplet"
				},
				"related": {
				  "href": "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/droplets/current"
				}
			  }
			}
			`}
			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			err := resourcesData.AssignApp("5EB2DB8D-2808-4871-9AF2-0A0977881A4B", "F6E3E504-99C2-41BD-B78D-262F3E70A1F8")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetRoutesApp with curl v3 api", func() {
		It("use v3 api to get routes", func() {
			response := []string{`
						  {
			  "pagination": {
				"total_results": 3,
				"total_pages": 2,
				"first": {
				  "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5/route_mappings?page=1&per_page=2"
				},
				"last": {
				  "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5/route_mappings?page=2&per_page=2"
				},
				"next": {
				  "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5/route_mappings?page=2&per_page=2"
				},
				"previous": null
			  },
			  "resources": [
				{
				  "guid": "89323d4e-2e84-43e7-83e9-adbf50a20c0e",
				  "created_at": "2016-02-17T01:50:05Z",
				  "updated_at": "2016-06-08T16:41:26Z",
				  "weight": 65,
				  "links": {
					"self": {
					  "href": "https://api.example.org/v3/route_mappings/89323d4e-2e84-43e7-83e9-adbf50a20c0e"
					},
					"app": {
					  "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5"
					},
					"route": {
					  "href": "https://api.example.org/v2/routes/9612ecbd-36f1-4ada-927a-efae9310b3a1"
					},
					"process": {
					  "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5/processes/web"
					}
				  }
				},
				{
				  "guid": "9f4970a8-9e6f-4984-b0a5-5e4a8af91665",
				  "created_at": "2016-02-17T01:50:05Z",
				  "updated_at": "2016-06-08T16:41:26Z",
				  "weight": 35,
				  "links": {
					"self": {
					  "href": "https://api.example.org/v3/route_mappings/9f4970a8-9e6f-4984-b0a5-5e4a8af91665"
					},
					"app": {
					  "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5"
					},
					"route": {
					  "href": "https://api.example.org/v2/routes/a32332f0-fb30-4e9e-9b78-348b8b6b98b6"
					},
					"process": {
					  "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5/processes/admin-web"
					}
				  }
				}
			  ]
			}
			`}
			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			routes, err := resourcesData.GetRoutesApp("5EB2DB8D-2808-4871-9AF2-0A0977881A4B")

			Expect(routes[0]).To(Equal("9612ecbd-36f1-4ada-927a-efae9310b3a1"))
			Expect(routes[1]).To(Equal("a32332f0-fb30-4e9e-9b78-348b8b6b98b6"))
			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
