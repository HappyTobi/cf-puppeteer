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
					"guid": "2f35885d-0c9d-4423-83ad-fd05066f8576",
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
			Expect(pushResponse.GUID).To(Equal("2f35885d-0c9d-4423-83ad-fd05066f8576"))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("createPackage with curl v3 api", func() {
		It("create package for created application", func() {
			response := []string{
				`
				{
					"guid": "44f7c078-0934-470f-9883-4fcddc5b8f13",
					"type": "bits",
					"data": {
					  "checksum": {
						"type": "sha256",
						"value": null
					  },
					  "error": null
					},
					"state": "PROCESSING_UPLOAD",
					"created_at": "2015-11-13T17:02:56Z",
					"updated_at": "2016-06-08T16:41:26Z",
					"metadata": {
					  "labels": { },
					  "annotations": { }
					},
					"links": {
					  "self": {
						"href": "https://api.example.org/v3/packages/44f7c078-0934-470f-9883-4fcddc5b8f13"
					  },
					  "upload": {
						"href": "https://api.example.org/v3/packages/44f7c078-0934-470f-9883-4fcddc5b8f13/upload",
						"method": "POST"
					  },
					  "download": {
						"href": "https://api.example.org/v3/packages/44f7c078-0934-470f-9883-4fcddc5b8f13/download",
						"method": "GET"
					  },
					  "app": {
						"href": "https://api.example.org/v3/apps/1d3bf0ec-5806-43c4-b64e-8364dba1086a"
					  }
					}
				  }
				  `,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			//pass packageGUID
			pushResponse, err := cf.CreatePackage("2f35885d-0c9d-4423-83ad-fd05066f8576")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(pushResponse.GUID).To(Equal("44f7c078-0934-470f-9883-4fcddc5b8f13"))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("uploadApplication with curl v3 api", func() {
		It("create package for created< application", func() {
			response := []string{
				`
				{
					"guid": "44f7c078-0934-470f-9883-4fcddc5b8f13",
					"type": "bits",
					"data": {
					  "checksum": {
						"type": "sha256",
						"value": null
					  },
					  "error": null
					},
					"state": "PROCESSING_UPLOAD",
					"created_at": "2015-11-13T17:02:56Z",
					"updated_at": "2016-06-08T16:41:26Z",
					"metadata": {
					  "labels": { },
					  "annotations": { }
					},
					"links": {
					  "self": {
						"href": "https://api.example.org/v3/packages/44f7c078-0934-470f-9883-4fcddc5b8f13"
					  },
					  "upload": {
						"href": "https://api.example.org/v3/packages/44f7c078-0934-470f-9883-4fcddc5b8f13/upload",
						"method": "POST"
					  },
					  "download": {
						"href": "https://api.example.org/v3/packages/44f7c078-0934-470f-9883-4fcddc5b8f13/download",
						"method": "GET"
					  },
					  "app": {
						"href": "https://api.example.org/v3/apps/1d3bf0ec-5806-43c4-b64e-8364dba1086a"
					  }
					}
				  }
				  `,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			//pass packageGUID
			pushResponse, err := cf.CreatePackage("2f35885d-0c9d-4423-83ad-fd05066f8576")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(pushResponse.GUID).To(Equal("44f7c078-0934-470f-9883-4fcddc5b8f13"))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("createBuild with curl v3 api", func() {
		It("create build", func() {
			response := []string{
				`
				{
					"guid": "585bc3c1-3743-497d-88b0-403ad6b56d16",
					"created_at": "2016-03-28T23:39:34Z",
					"updated_at": "2016-06-08T16:41:26Z",
					"created_by": {
						"guid": "3cb4e243-bed4-49d5-8739-f8b45abdec1c",
						"name": "bill",
						"email": "bill@example.com"
					},
					"state": "STAGING",
					"error": null,
					"lifecycle": {
						"type": "buildpack",
						"data": {
							"buildpacks": [ "ruby_buildpack" ],
							"stack": "cflinuxfs3"
						}
					},
					"package": {
						"guid": "8e4da443-f255-499c-8b47-b3729b5b7432"
					},
					"droplet": null,
					"metadata": {
						"labels": { },
						"annotations": { }
					},
					"links": {
						"self": {
							"href": "https://api.example.org/v3/builds/585bc3c1-3743-497d-88b0-403ad6b56d16"
						},
						"app": {
							"href": "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396"
						}
					}
				}
				
				  `,
			}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)
			//pass packageGUID
			pushResponse, err := cf.CreateBuild("2f35885d-0c9d-4423-83ad-fd05066f8576")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(pushResponse.GUID).To(Equal("585bc3c1-3743-497d-88b0-403ad6b56d16"))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
