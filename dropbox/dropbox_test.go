package dropbox

import (
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestDropbox(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("generate UUID", func() {
		g.It("should generate a uuid v3", func() {
			rando, err := Randomizer()
			Expect(err).To(BeNil())
			Expect(len(rando)).To(Equal(32))
		})
	})
}
