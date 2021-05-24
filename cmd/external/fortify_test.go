package external

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/franela/goblin"
	"github.com/ion-channel/ionize/dropbox"
	. "github.com/onsi/gomega"
)

type mockS3Client struct {
	s3iface.S3API
}

func (c mockS3Client) PutObjectRequest(input *s3.PutObjectInput) (req *request.Request, output *s3.PutObjectOutput) {
	out := s3.PutObjectOutput{}
	req = new(request.Request)
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.HTTPRequest = r
	return req, &out
}

func TestFortify(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("FPR file handling", func() {
		g.BeforeEach(func() {
			mc := mockS3Client{}
			dropbox.Uploader = s3manager.NewUploaderWithClient(mc)
		})

		g.It("should unzip and fpr file", func() {
			dir, _ := filepath.Abs(filepath.Join(os.Getenv("PWD"), "..", ".."))

			_, err := unzip(strings.Join([]string{dir, "fortify.zip"}, "/"))
			Expect(err).To(BeNil())
		})

		g.It("should parse an fpr file", func() {
			dir, _ := filepath.Abs(filepath.Join(os.Getenv("PWD"), "..", ".."))

			path := strings.Join([]string{dir, "fortify.zip"}, "/")

			fort, err := ParseFortify(path)
			Expect(err).To(BeNil())
			Expect(fort.Value).NotTo(BeNil())
			// matches the pdf for the fpr input
			Expect(fort.Value.Vulnerability.Critcal).To(Equal(43))
			Expect(fort.Value.Vulnerability.High).To(Equal(262))
			Expect(fort.Value.Vulnerability.Medium).To(Equal(0))
			Expect(fort.Value.Vulnerability.Low).To(Equal(79))
		})

		g.It("should gather all of the rules", func() {
			dir, _ := filepath.Abs(filepath.Join(os.Getenv("PWD"), "..", ".."))

			path := strings.Join([]string{dir, "fortify.zip"}, "/")

			fort, err := ParseFortify(path)
			Expect(err).To(BeNil())

			rules := fort.FVDL.Rules()
			Expect(rules).NotTo(BeNil())
			Expect(len(rules)).To(Equal(42))
		})

		g.It("should get group value for rule and group name", func() {
			dir, _ := filepath.Abs(filepath.Join(os.Getenv("PWD"), "..", ".."))

			path := strings.Join([]string{dir, "fortify.zip"}, "/")

			fort, _ := ParseFortify(path)

			value := fort.FVDL.Group("10683D0C-25FA-4984-41CC-651C955D640A", Accuracy)
			Expect(value).NotTo(Equal(""))
			Expect(value).To(Equal("4.0"))
		})
	})
}
