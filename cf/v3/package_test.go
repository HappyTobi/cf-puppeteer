package v3_test

import (
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"fmt"
	"github.com/happytobi/cf-puppeteer/cf/cli"
	"github.com/happytobi/cf-puppeteer/cf/v3"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfPackageActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Package")
}

var _ = Describe("cf-package test", func() {
	var (
		cliConn       *pluginfakes.FakeCliConnection
		resourcesData *v3.ResourcesData
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		resourcesData = &v3.ResourcesData{Connection: cliConn, Cli: cli.NewCli(cliConn, true)}

	})
	Describe("Test temp file generation", func() {
		It("app-name without prefix", func() {
			appName := "myApplication"
			zipFile := resourcesData.GenerateTempFile(appName, "zip")
			fmt.Printf("zipFilePath %s", zipFile)
			Expect(strings.HasSuffix(zipFile, "/myApplication.zip")).To(Equal(true))
			Expect(strings.HasSuffix(zipFile, "//myApplication.zip")).To(Equal(false))
		})

		It("app-name with prefix", func() {
			appName := "/myApplication"
			zipFile := resourcesData.GenerateTempFile(appName, "zip")
			fmt.Printf("zipFilePath %s", zipFile)
			Expect(strings.HasSuffix(zipFile, "/myApplication.zip")).To(Equal(true))
			Expect(strings.HasSuffix(zipFile, "//myApplication.zip")).To(Equal(false))

		})
	})
})
