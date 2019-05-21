package cf

import (
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer Version Check Test")
}

var _ = Describe("cf push version check", func() {
	var (
		cliConn *pluginfakes.FakeCliConnection
		push    PuppeteerPush
	)

	BeforeEach(func() {
		cliConn = &pluginfakes.FakeCliConnection{}
		push = NewApplicationPush(cliConn, false)
	})

	Describe("check controller call", func() {
		It("should return version", func() {
			push.pushApplication()
			/*v2Ver, v3Ver, err := getCloudControllerAPIVersion()
			Expect(v3Ver).To(Equal("3.27.0"))
			Expect(v2Ver).To(Equal("3.27.0"))
			Expect(err).ToNot(HaveOccurred())*/

		})

		/*It("check sem version", func() {
			v3SemVer, err := GetMinSemVersion()

			Expect(v3SemVer.String()).To(Equal("3.27.0"))
			Expect(err).ToNot(HaveOccurred())
		})*/
	})
})
