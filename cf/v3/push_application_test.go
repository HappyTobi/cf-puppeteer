package v3_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	"github.com/happytobi/cf-puppeteer/cf/v3"
	manifest "github.com/happytobi/cf-puppeteer/manifest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfPushApplicationActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Push Application")
}

var _ = Describe("cf-push application test", func() {
	var (
		cliConn       *pluginfakes.FakeCliConnection
		resourcesData *v3.ResourcesData
		manifestFile  manifest.Manifest
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, true)}
		manifestFile, _ = manifest.ParseApplicationManifest("../../fixtures/manifest.yml", "")
	})
	Describe("Test temp file generation without routes", func() {
		It("app-name without prefix", func() {
			noRouteYmlPath, err := resourcesData.GenerateNoRouteYml("my-test-application", manifestFile)
			noRouteYml, errNoRouteYml := manifest.ParseApplicationManifest(noRouteYmlPath, "")

			Expect(err).ToNot(HaveOccurred())
			Expect(errNoRouteYml).ToNot(HaveOccurred())

			Expect(len(manifestFile.ApplicationManifests[0].Routes)).To(Equal(2))
			Expect(len(noRouteYml.ApplicationManifests[0].Routes)).To(Equal(0))
			Expect(manifestFile.ApplicationManifests[0].DiskQuota).To(Equal(noRouteYml.ApplicationManifests[0].DiskQuota))
			Expect(manifestFile.ApplicationManifests[0].Instances).To(Equal(noRouteYml.ApplicationManifests[0].Instances))
			Expect(manifestFile.ApplicationManifests[0].Memory).To(Equal(noRouteYml.ApplicationManifests[0].Memory))
		})
	})
})
