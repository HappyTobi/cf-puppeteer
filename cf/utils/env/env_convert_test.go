package env_test

import (
	"github.com/happytobi/cf-puppeteer/cf/utils/env"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCfUtilEnvConverter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Puppeteer CF Util Resources")
}

var _ = Describe("cf util environment converter", func() {
	Describe("convert simple environment variables to map", func() {
		It("convert simple", func() {
			var simpleVars = []string{"foo=bar", "santa=claus"}
			convertedEnvs := env.Convert(simpleVars)
			Expect(convertedEnvs["foo"]).To(Equal("bar"))
			Expect(convertedEnvs["santa"]).To(Equal("claus"))
			Expect(len(convertedEnvs)).To(Equal(2))
		})
	})

	Describe("convert simple environment variable with equals value to map", func() {
		It("convert simple equals value", func() {
			var simpleVars = []string{"foo=bar", "santa=claus=cool"}
			convertedEnvs := env.Convert(simpleVars)
			Expect(convertedEnvs["foo"]).To(Equal("bar"))
			Expect(convertedEnvs["santa"]).To(Equal("claus=cool"))
			Expect(len(convertedEnvs)).To(Equal(2))
		})
	})

	Describe("convert complex environment variables to map", func() {
		It("convert complex", func() {
			var simpleVars = []string{"jdbc=jdbc:oracle:thin:username/password@amrood:1521:EMP", "fstab=\"uid=1000,gid=100,umask=0,allow_other\""}
			convertedEnvs := env.Convert(simpleVars)
			Expect(convertedEnvs["jdbc"]).To(Equal("jdbc:oracle:thin:username/password@amrood:1521:EMP"))
			Expect(convertedEnvs["fstab"]).To(Equal("uid=1000,gid=100,umask=0,allow_other"))
			Expect(len(convertedEnvs)).To(Equal(2))
		})

		It("convert complex with single quotes", func() {
			var simpleVars = []string{"jdbc=jdbc:oracle:thin:username/password@amrood:1521:EMP", "fstab='uid=1000,gid=100,umask=0,allow_other'"}
			convertedEnvs := env.Convert(simpleVars)
			Expect(convertedEnvs["jdbc"]).To(Equal("jdbc:oracle:thin:username/password@amrood:1521:EMP"))
			Expect(convertedEnvs["fstab"]).To(Equal("uid=1000,gid=100,umask=0,allow_other"))
			Expect(len(convertedEnvs)).To(Equal(2))
		})
	})

	Describe("dont convert var's didn't match the pattern", func() {
		It("convert non patter matching vars", func() {
			var simpleVars = []string{"foo_bar", "santa-claus"}
			convertedEnvs := env.Convert(simpleVars)
			Expect(len(convertedEnvs)).To(Equal(0))
		})
	})
})
