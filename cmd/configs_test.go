package cmd

import (
	"testing"

	. "github.com/franela/goblin"
	"github.com/gomicro/penname"
	. "github.com/onsi/gomega"
)

func TestFunctionality(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Configs Command", func() {
		g.It("should output the loaded configs", func() {
			mw := penname.New()

			output = mw
			cfgFile = "somefile.yaml"

			runConfigsCmd(nil, nil)
			Expect(string(mw.Written())).To(Equal("Config File: \nSecret Key: \nAPI: \n"))
		})
	})
}
