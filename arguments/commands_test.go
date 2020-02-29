package arguments

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands TestSuite")
}

var _ = Describe("List Commands", func() {
	var ()

	BeforeEach(func() {

	})

	It("load all commands", func() {
		commands := commands()
		Expect(len(commands)).Should(Equal(18))
		Expect(commands["env"]).Should(Equal("Variable key value pair for adding dynamic environment variables; can specify multiple times"))
		Expect(commands["f"]).Should(Equal("path to an application manifest"))
	})

	It("load all UsageDetail", func() {
		commands := UsageDetailsOptionCommands()
		Expect(len(commands)).Should(Equal(18))
		Expect(commands["-env"]).Should(Equal("Variable key value pair for adding dynamic environment variables; can specify multiple times"))
		Expect(commands["f"]).Should(Equal("path to an application manifest"))
	})
})
