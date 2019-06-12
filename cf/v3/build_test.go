package v3_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	v3 "github.com/happytobi/cf-puppeteer/cf/v3"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfBuildActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Resources")
}

var _ = Describe("cf-build test", func() {
	var (
		cliConn       *pluginfakes.FakeCliConnection
		resourcesData *v3.ResourcesData
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, true)}

	})
	Describe("build actions with curl v3 api", func() {
		It("use v3 create build api", func() {
			response := []string{`{
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
				  "droplet": {
						"guid": "3cb4e243-bed4-49d5-8739-f8b45abdec1x"
					},
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
				}`}
			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			var buildPacks = []string{"ruby_buildpack"}
			buildResponse, err := resourcesData.CreateBuild("585bc3c1-3743-497d-88b0-403ad6b56d16", buildPacks)

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			Expect(buildResponse.GUID).To(Equal("585bc3c1-3743-497d-88b0-403ad6b56d16"))
			Expect(buildResponse.State).To(Equal("STAGING"))
			Expect(buildResponse.Droplet.GUID).To(Equal("3cb4e243-bed4-49d5-8739-f8b45abdec1x"))

			Expect(err).ToNot(HaveOccurred())
		})

		It("use v3 create buuld api without a buildpack", func() {
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
	Describe("CheckBuildState with curl v3 api", func() {
		It("use v3 to check build state", func() {
			response := []string{`{
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
				  "droplet": {
						"guid": "3cb4e243-bed4-49d5-8739-f8b45abdec1x"
					},
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
				}`}

			response2 := []string{`{
				  "guid": "585bc3c1-3743-497d-88b0-403ad6b56d16",
				  "created_at": "2016-03-28T23:39:34Z",
				  "updated_at": "2016-06-08T16:41:26Z",
				  "created_by": {
					"guid": "3cb4e243-bed4-49d5-8739-f8b45abdec1c",
					"name": "bill",
					"email": "bill@example.com"
				  },
				  "state": "FAILED",
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
				  "droplet": {
						"guid": "3cb4e243-bed4-49d5-8739-f8b45abdec1x"
					},
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
				}`}
			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			buildResponse, err := resourcesData.CheckBuildState("585bc3c1-3743-497d-88b0-403ad6b56d16")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			Expect(buildResponse.State).To(Equal("STAGING"))

			cliConn.CliCommandWithoutTerminalOutputReturns(response2, nil)
			buildResponse, err = resourcesData.CheckBuildState("585bc3c1-3743-497d-88b0-403ad6b56d16")
			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))

			Expect(buildResponse.State).To(Equal("FAILED"))

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("GetDropletGUID with curl v3 api", func() {
		It("use v3 to check build state", func() {
			response := []string{`{
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
				  "droplet": {
						"guid": "3cb4e243-bed4-49d5-8739-f8b45abdec1x"
					},
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
				}`}

			cliConn.CliCommandWithoutTerminalOutputReturns(response, nil)

			dropletResponse, err := resourcesData.GetDropletGUID("585bc3c1-3743-497d-88b0-403ad6b56d16")

			Expect(cliConn.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))

			Expect(dropletResponse.Droplet.GUID).To(Equal("3cb4e243-bed4-49d5-8739-f8b45abdec1x"))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
