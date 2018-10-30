package dropbox

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spf13/viper"
)

//Uploader s3 uploader global for mocking
var Uploader *s3manager.Uploader

//Randomizer generate a random UUID for dropbox entries
func Randomizer() (string, error) {
	var uuid [16]byte
	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return hex.EncodeToString(uuid[:]), nil
}

//ParseURL takes a potential url.  Based on scheme will either upload to a
//bucket and return a http url or just return the url
func ParseURL(input, randomizer string) (string, error) {
	url, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %v", err.Error())
	}

	if randomizer != "" {
		randomizer = strings.Join([]string{randomizer, "/"}, "")
	}

	// It is a local file.  upload to a bucket
	// and create a timed url for downloading after analysis
	if url.Scheme == "" || url.Scheme == "file" {

		f, err := os.Open(url.Hostname() + url.EscapedPath())
		if err != nil {
			return "", fmt.Errorf("failed to read file for url (%s): %v", url.String(), err.Error())
		}

		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})

		if Uploader == nil {
			Uploader = s3manager.NewUploader(sess)
		}

		_, err = Uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(viper.GetString("bucket")),
			Key:    aws.String("ionize/" + randomizer + f.Name()),
			Body:   f,
		})
		if err != nil {
			return "", fmt.Errorf("failed to write file to Ion Channel: %v", err.Error())
		}

		svc := s3.New(sess)
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(viper.GetString("bucket")),
			Key:    aws.String("ionize/" + randomizer + f.Name()),
		})
		erl, err := req.Presign(15 * time.Minute)
		if err != nil {
			return "", fmt.Errorf("Failed to sign request: %v", err)
		}

		return erl, nil
	}

	return url.String(), nil
}
