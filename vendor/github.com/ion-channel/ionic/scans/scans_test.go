package scans

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestScan(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Scan", func() {
		g.Describe("Translating", func() {
			g.It("should translate an untranslated scan", func() {
				s := &Scan{
					UntranslatedResults: &UntranslatedResults{
						License: &LicenseResults{},
					},
				}
				Expect(s.UntranslatedResults).NotTo(BeNil())
				Expect(s.TranslatedResults).To(BeNil())

				err := s.Translate()
				Expect(err).To(BeNil())
				Expect(s.UntranslatedResults).To(BeNil())
				Expect(s.TranslatedResults).NotTo(BeNil())
				Expect(s.TranslatedResults.Type).To(Equal("license"))
				Expect(s.Results).NotTo(BeNil())
				Expect(len(s.Results)).NotTo(Equal(0))
			})

			g.It("should not translate an already translated scan", func() {
				s := &Scan{
					UntranslatedResults: &UntranslatedResults{
						License: &LicenseResults{},
					},
				}
				Expect(s.UntranslatedResults).NotTo(BeNil())
				Expect(s.TranslatedResults).To(BeNil())

				err := s.Translate()
				Expect(err).To(BeNil())
				Expect(s.UntranslatedResults).To(BeNil())
				Expect(s.TranslatedResults).NotTo(BeNil())
				Expect(s.TranslatedResults.Type).To(Equal("license"))
				Expect(s.Results).NotTo(BeNil())
				Expect(len(s.Results)).NotTo(Equal(0))

				err = s.Translate()
				Expect(err).To(BeNil())
				Expect(s.UntranslatedResults).To(BeNil())
				Expect(s.TranslatedResults).NotTo(BeNil())
				Expect(s.TranslatedResults.Type).To(Equal("license"))
				Expect(s.Results).NotTo(BeNil())
				Expect(len(s.Results)).NotTo(Equal(0))
			})

			g.It("should return string in JSON", func() {
				createdAt := time.Date(2018, 07, 07, 13, 42, 47, 651387237, time.UTC)
				updatedAt := time.Date(2018, 07, 07, 13, 42, 47, 651387237, time.UTC)
				r := json.RawMessage(`{"result": "someresult"}`)

				s, _ := NewScan("someid", "someteamid", "someprojectid", "someanalysisid", "somesummary", "somename", "somedesc", r, createdAt, updatedAt, 19.22)
				Expect(fmt.Sprintf("%v", s)).To(Equal(`{"id":"someid","team_id":"someteamid","project_id":"someprojectid","analysis_id":"someanalysisid","summary":"somesummary","results":{"result":"someresult"},"created_at":"2018-07-07T13:42:47.651387237Z","updated_at":"2018-07-07T13:42:47.651387237Z","duration":19.22,"name":"somename","description":"somedesc"}`))
			})
		})

		g.Describe("Unmarshalling", func() {
			g.It("should populate results with untranslated result", func() {
				var ss Scan
				err := json.Unmarshal([]byte(sampleUntranslatedResults), &ss)

				Expect(err).To(BeNil())
				Expect(ss.TeamID).To(Equal("cuketest"))
				Expect(ss.TranslatedResults).To(BeNil())
				Expect(ss.UntranslatedResults).NotTo(BeNil())
				Expect(ss.UntranslatedResults.License).NotTo(BeNil())
				Expect(ss.UntranslatedResults.License.Name).To(Equal("some license"))
				Expect(ss.UntranslatedResults.License.Type[0].Confidence).To(Equal(float32(1.0)))
			})

			g.It("should populate results with translated result", func() {
				var ss Scan
				err := json.Unmarshal([]byte(sampleTranslatedResults), &ss)

				Expect(err).To(BeNil())
				Expect(ss.TeamID).To(Equal("cuketest"))
				Expect(ss.UntranslatedResults).To(BeNil())
				Expect(ss.TranslatedResults).NotTo(BeNil())
				Expect(ss.TranslatedResults.Type).To(Equal("community"))
			})

			g.It("should unmarshal a scan with a bad, list-ified community results member", func() {
				var ss Scan
				err := json.Unmarshal([]byte(badCommunityResults), &ss)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		g.Describe("Marshalling", func() {
			g.It("should output results with untranslated result", func() {
				s := &Scan{
					scan: &scan{
						ID:          "41e6905a-a16d-45a7-9d2c-2794840ca03e",
						TeamID:      "cuketest",
						AnalysisID:  "c2a79402-e3bf-4069-89c8-7a4ecb10d33f",
						CreatedAt:   time.Now(),
						Description: "This scan data has not been evaluated against a rule.",
						Duration:    1000.1,
						Name:        "license",
						ProjectID:   "35b06118-da91-4ac8-a3d0-db25a3e554c5",
						Summary:     "oh hi",
						UpdatedAt:   time.Now(),
					},
					UntranslatedResults: &UntranslatedResults{
						License: &LicenseResults{
							&License{
								Name: "some license",
								Type: []LicenseType{
									LicenseType{
										Name:       "a license",
										Confidence: 1.0,
									},
								},
							},
						},
					},
				}

				b, err := json.MarshalIndent(s, "", "  ")
				Expect(err).To(BeNil())

				body := string(b)
				Expect(body).To(ContainSubstring("team_id\": \"cuketest\""))
				Expect(body).NotTo(ContainSubstring("about_yml"))
				Expect(body).NotTo(ContainSubstring("community"))
				Expect(body).NotTo(ContainSubstring("coverage"))
				Expect(body).NotTo(ContainSubstring("dependency"))
				Expect(body).NotTo(ContainSubstring("difference"))
				Expect(body).NotTo(ContainSubstring("ecosystem"))
				Expect(body).NotTo(ContainSubstring("external_vulnerability"))
				Expect(body).NotTo(ContainSubstring("virus"))
				Expect(body).NotTo(ContainSubstring("vulnerability"))
			})

			g.It("should output results with translated result", func() {
				s := &Scan{
					scan: &scan{
						ID:          "41e6905a-a16d-45a7-9d2c-2794840ca03e",
						TeamID:      "cuketest",
						AnalysisID:  "c2a79402-e3bf-4069-89c8-7a4ecb10d33f",
						CreatedAt:   time.Now(),
						Description: "This scan data has not been evaluated against a rule.",
						Duration:    1000.1,
						Name:        "community",
						ProjectID:   "35b06118-da91-4ac8-a3d0-db25a3e554c5",
						Summary:     "oh hi",
						UpdatedAt:   time.Now(),
					},
					TranslatedResults: &TranslatedResults{
						Type: "community",
						Data: &CommunityResults{
							Committers: 5,
							Name:       "reponame",
							URL:        "http://github.com/reponame",
						},
					},
				}

				b, err := json.MarshalIndent(s, "", "  ")
				Expect(err).To(BeNil())

				body := string(b)
				Expect(body).To(ContainSubstring("team_id\": \"cuketest\""))
				Expect(body).To(ContainSubstring("committers\": 5"))
				Expect(body).NotTo(ContainSubstring("about_yml"))
				Expect(body).NotTo(ContainSubstring("coverage"))
				Expect(body).NotTo(ContainSubstring("dependency"))
				Expect(body).NotTo(ContainSubstring("difference"))
				Expect(body).NotTo(ContainSubstring("ecosystem"))
				Expect(body).NotTo(ContainSubstring("external_vulnerability"))
				Expect(body).NotTo(ContainSubstring("license"))
				Expect(body).NotTo(ContainSubstring("virus"))
				Expect(body).NotTo(ContainSubstring("vulnerability"))
			})
		})
	})
}
