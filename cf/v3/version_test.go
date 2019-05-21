package v3

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer Version Check Test")
}

var _ = Describe("check version", func() {

	BeforeEach(func() {

	})
	Describe("check min v3 version", func() {
		It("check version string", func() {
			Expect(MinControllerVersion).To(Equal("3.27.0"))
		})

		It("check sem version", func() {
			v3SemVer, err := GetMinSemVersion()

			Expect(v3SemVer.String()).To(Equal("3.27.0"))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
