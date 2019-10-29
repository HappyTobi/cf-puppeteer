package manifest_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/happytobi/cf-puppeteer/manifest"
)

func TestManifestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manifest Testsuite")
}

var _ = Describe("Parse Manifest", func() {
	It("parses complete manifest", func() {
		manifest, err := Parse("../fixtures/manifest.yml")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.ApplicationManifests[0].Name).Should(Equal("myApp"))
		Expect(manifest.ApplicationManifests[0].Buildpacks[0]).Should(Equal("java_buildpack"))
		Expect(manifest.ApplicationManifests[0].Buildpacks[1]).Should(Equal("go_buildpack"))
		Expect(manifest.ApplicationManifests[0].Timeout).Should(Equal(2))
	})

	It("parses complete manifest with services", func() {
		manifest, err := Parse("../fixtures/manifest.yml")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.ApplicationManifests[0].Name).Should(Equal("myApp"))
		Expect(manifest.ApplicationManifests[0].Services[0]).Should(Equal("service1"))
		Expect(manifest.ApplicationManifests[0].Services[1]).Should(Equal("service2"))
		Expect(manifest.ApplicationManifests[0].AppPath).Should(Equal("myJavaApp.jar"))
	})
	It("parses complete manifest with buildpack url", func() {
		manifest, err := Parse("../fixtures/phpManifest.yml")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.ApplicationManifests[0].Name).Should(Equal("appname"))
		Expect(manifest.ApplicationManifests[0].Services[0]).Should(Equal("ma-db"))
		Expect(manifest.ApplicationManifests[0].Services[1]).Should(Equal("app-db"))
		Expect(manifest.ApplicationManifests[0].Services[2]).Should(Equal("credentials"))
		Expect(manifest.ApplicationManifests[0].Stack).Should(Equal("cflinuxfs3"))
		Expect(manifest.ApplicationManifests[0].Buildpacks[0]).Should(Equal("https://github.com/cloudfoundry/php-buildpack.git"))
		Expect(manifest.ApplicationManifests[0].Buildpacks[1]).Should(Equal("https://github.com/cloudfoundry/php-buildpack.git"))
	})
})

var _ = Describe("Parse multi Application Manifest", func() {
	It("parses complete manifest", func() {
		manifest, err := Parse("../fixtures/multiManifest.yml")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.ApplicationManifests[0].Name).Should(Equal("myApp"))
		Expect(manifest.ApplicationManifests[1].Name).Should(Equal("myApp2"))
	})
})

var _ = Describe("Parse invalid Application Manifest", func() {
	It("parses invalid manifest", func() {
		manifest, err := Parse("../fixtures/invalidManifest.yml")
		Expect(err).ShouldNot(BeNil())
		Expect(manifest.ApplicationManifests).Should(BeNil())
	})
})

var _ = Describe("Parse comp Manifest", func() {
	It("parses complicated manifest", func() {
		manifest, err := Parse("../fixtures/defaultMultiManifest.yml")
		Expect(err).Should(BeNil())
		Expect(manifest.ApplicationManifests).ShouldNot(BeNil())
		Expect(manifest.ApplicationManifests[0].Name).Should(Equal("app"))
		Expect(manifest.ApplicationManifests[0].DiskQuota).Should(Equal("1G"))
	})
})

var _ = Describe("Write new manifest", func() {
	It("write manifest file to specified path", func() {
		manifest, err := Parse("../fixtures/manifest.yml")
		Expect(err).ShouldNot(HaveOccurred())
		tempFile := fmt.Sprintf("%s%s", os.TempDir(), "testManifest.yml")
		parsedTempManifest, err := WriteYmlFile(tempFile, manifest)
		Expect(err).ShouldNot(HaveOccurred())
		fmt.Printf("%s", tempFile)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.ApplicationManifests[0].Name).Should(Equal(parsedTempManifest.ApplicationManifests[0].Name))
	})
})
