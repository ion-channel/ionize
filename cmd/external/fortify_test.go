package external

import (
	"path/filepath"
	"strings"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestFortify(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("FPR file handling", func() {
		g.It("should unzip and fpr file", func() {
			dir, _ := filepath.Abs("/tmp")

			path, err := unzip(strings.Join([]string{dir, "fortifyResults-nginx-33.fpr"}, "/"))
			Expect(err).To(BeNil())
			Expect(path).To(Equal("/tmp"))
		})

		g.It("should parse an fpr file", func() {
			dir, _ := filepath.Abs("/tmp")

			path := strings.Join([]string{dir, "fortifyResults-nginx-33.fpr"}, "/")

			fort, err := ParseFortify(path)
			Expect(err).To(BeNil())
			Expect(fort.Value).NotTo(BeNil())
			// matches the pdf for the fpr input
			Expect(fort.Value.Vulnerabilities.Critcal).To(Equal(43))
			Expect(fort.Value.Vulnerabilities.High).To(Equal(262))
			Expect(fort.Value.Vulnerabilities.Medium).To(Equal(0))
			Expect(fort.Value.Vulnerabilities.Low).To(Equal(79))
		})

		g.It("should gather all of the rules", func() {
			dir, _ := filepath.Abs("/tmp")

			path := strings.Join([]string{dir, "fortifyResults-nginx-33.fpr"}, "/")

			fort, _ := ParseFortify(path)

			rules := fort.FVDL.Rules()
			Expect(rules).NotTo(BeNil())
			Expect(len(rules)).To(Equal(42))
		})

		g.It("should get group value for rule and group name", func() {
			dir, _ := filepath.Abs("/tmp")

			path := strings.Join([]string{dir, "fortifyResults-nginx-33.fpr"}, "/")

			fort, _ := ParseFortify(path)

			value := fort.FVDL.Group("10683D0C-25FA-4984-41CC-651C955D640A", Accuracy)
			Expect(value).NotTo(Equal(""))
			Expect(value).To(Equal("4.0"))
		})
	})
}
