package scans

import (
	"encoding/json"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestScanResults(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Untranslated Scan Results", func(){
		g.It("should translate untranslated scan results", func(){
			var untranslatedResult UntranslatedResults
			err := json.Unmarshal([]byte(SampleValidUntranslatedScanResultsLicense), &untranslatedResult)

			// validate the json parsing
			Expect(err).NotTo(HaveOccurred())
			Expect(untranslatedResult.AboutYML).To(BeNil())
			Expect(untranslatedResult.Community).To(BeNil())
			Expect(untranslatedResult.Coverage).To(BeNil())
			Expect(untranslatedResult.Dependency).To(BeNil())
			Expect(untranslatedResult.Difference).To(BeNil())
			Expect(untranslatedResult.Ecosystem).To(BeNil())
			Expect(untranslatedResult.ExternalVulnerabilities).To(BeNil())
			Expect(untranslatedResult.Vulnerability).To(BeNil())
			Expect(untranslatedResult.License).NotTo(BeNil())
			license := untranslatedResult.License
			Expect(license.Name).To(Equal("some license"))
			Expect(license.Type).To(HaveLen(1))
			Expect(license.Type[0].Name).To(Equal("a license"))

			// translate it
			translatedResult := untranslatedResult.Translate()

			// validate translated object
			Expect(translatedResult).NotTo(BeNil())
			Expect(translatedResult.Type).To(Equal("license"))
			Expect(translatedResult.Data).NotTo(BeNil())
			wasLicenseResults := false
			switch translatedResult.Data.(type) {
			case *LicenseResults:
				wasLicenseResults = true
			}
			Expect(wasLicenseResults).To(BeTrue())
			licenseResults := translatedResult.Data.(*LicenseResults)
			Expect(licenseResults.Type).To(HaveLen(1))
			Expect(licenseResults.Type[0].Name).To(Equal("a license"))
			Expect(licenseResults.Name).To(Equal("some license"))
			Expect(licenseResults.License.Type).To(HaveLen(1))
			Expect(licenseResults.License.Type[0].Name).To(Equal("a license"))
		})
	})
	g.Describe("Translated Scan Results", func() {
		g.It("should unmarshal a scan results with about yml data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsAboutYML), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("about_yml"))

			a, ok := r.Data.(AboutYMLResults)
			Expect(ok).To(Equal(true))
			Expect(a.Content).To(Equal("some content"))
		})

		g.It("should unmarshal a scan results with community data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsCommunity), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("community"))

			a, ok := r.Data.(CommunityResults)
			Expect(ok).To(Equal(true))
			Expect(a.Committers).To(Equal(7))
			Expect(a.Name).To(Equal("ion-channel/ion-connect"))
			Expect(a.URL).To(Equal("https://github.com/ion-channel/ion-connect"))
		})

		g.It("should unmarshal a scan results with coverage data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsCoverage), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("external_coverage"))

			c, ok := r.Data.(CoverageResults)
			Expect(ok).To(Equal(true))
			Expect(c.Value).To(Equal(42.0))
		})

		g.It("should unmarshal a scan results with dependency data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsDependency), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("dependency"))

			d, ok := r.Data.(DependencyResults)
			Expect(ok).To(Equal(true))
			Expect(len(d.Dependencies)).To(Equal(7))
			Expect(d.Meta.FirstDegreeCount).To(Equal(3))
			Expect(d.Meta.NoVersionCount).To(Equal(0))
			Expect(d.Meta.TotalUniqueCount).To(Equal(7))
			Expect(d.Meta.UpdateAvailableCount).To(Equal(2))
		})

		g.It("should unmarshal a scan results with difference data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsDifference), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("difference"))

			d, ok := r.Data.(DifferenceResults)
			Expect(ok).To(Equal(true))
			Expect(d.Checksum).To(Equal("checksumishere"))
			Expect(d.Difference).To(BeTrue())
		})

		g.It("should unmarshal a scan results with ecosystem data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsEcosystems), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("ecosystems"))

			e, ok := r.Data.(EcosystemResults)
			Expect(ok).To(Equal(true))
			Expect(len(e.Ecosystems)).To(Equal(3))
		})

		g.It("should unmarshal a scan results with external vulnerabilities scan data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidExternalVulnerabilities), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("external_vulnerability"))

			e, ok := r.Data.(ExternalVulnerabilitiesResults)
			Expect(ok).To(Equal(true))
			Expect(e.Critical).To(Equal(1))
			Expect(e.High).To(Equal(0))
			Expect(e.Medium).To(Equal(1))
			Expect(e.Low).To(Equal(0))
		})

		g.It("should unmarshal a scan results with license data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsLicense), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("license"))

			l, ok := r.Data.(LicenseResults)
			Expect(ok).To(Equal(true))
			Expect(l.License.Name).To(Equal("Not found"))
		})

		g.It("should unmarshal a scan results with virus data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsVirus), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("virus"))

			v, ok := r.Data.(VirusResults)
			Expect(ok).To(Equal(true))
			Expect(v.KnownViruses).To(Equal(10))
		})

		g.It("should unmarshal a scan results with vulnerability data", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleValidScanResultsVulnerability), &r)

			Expect(err).To(BeNil())
			Expect(r.Type).To(Equal("vulnerability"))

			v, ok := r.Data.(VulnerabilityResults)
			Expect(ok).To(Equal(true))
			Expect(v.Meta.VulnerabilityCount).To(Equal(1))
		})

		g.It("should return an error for an invalid results type", func() {
			var r TranslatedResults
			err := json.Unmarshal([]byte(SampleInvalidResults), &r)

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("unsupported results type found:"))
		})
	})
}

const (
	SampleValidScanResultsAboutYML      = `{"type":"about_yml", "data":{"message": "foo message", "valid": true, "content": "some content"}}`
	SampleValidScanResultsCommunity     = `{"type":"community", "data":{"committers":7,"name":"ion-channel/ion-connect","url":"https://github.com/ion-channel/ion-connect"}}`
	SampleValidScanResultsCoverage      = `{"type":"external_coverage", "data":{"value":42.0}}`
	SampleValidScanResultsDependency    = `{"type":"dependency","data":{"dependencies":[{"latest_version":"2.0","org":"net.sourceforge.javacsv","name":"javacsv","type":"maven","package":"jar","version":"2.0","scope":"compile"},{"latest_version":"4.12","org":"junit","name":"junit","type":"maven","package":"jar","version":"4.11","scope":"test"},{"latest_version":"1.4-atlassian-1","org":"org.hamcrest","name":"hamcrest-core","type":"maven","package":"jar","version":"1.3","scope":"test"},{"latest_version":"4.5.2","org":"org.apache.httpcomponents","name":"httpclient","type":"maven","package":"jar","version":"4.3.4","scope":"compile"},{"latest_version":"4.4.5","org":"org.apache.httpcomponents","name":"httpcore","type":"maven","package":"jar","version":"4.3.2","scope":"compile"},{"latest_version":"99.0-does-not-exist","org":"commons-logging","name":"commons-logging","type":"maven","package":"jar","version":"1.1.3","scope":"compile"},{"latest_version":"20041127.091804","org":"commons-codec","name":"commons-codec","type":"maven","package":"jar","version":"1.6","scope":"compile"}],"meta":{"first_degree_count":3,"no_version_count":0,"total_unique_count":7,"update_available_count":2}}}`
	SampleValidScanResultsEcosystems    = `{"type":"ecosystems","data":{"ecosystems":[{"ecosystem":"Java","lines":2430},{"ecosystem":"Makefile","lines":210},{"ecosystem":"Ruby","lines":666}]}}`
	SampleValidScanResultsLicense       = `{"type":"license","data":{"license":{"name":"Not found","type":[]}}}`
	SampleValidScanResultsVirus         = `{"type":"virus","data":{"known_viruses":10,"engine_version":"","scanned_directories":1,"scanned_files":2,"infected_files":1,"data_scanned":"some cool data was scanned","data_read":"we read some data","time":"10PM","file_notes": {"empty_file":["file1","file2","file3"]},"clam_av_details":{"clamav_version":"1.0.0","clamav_db_version":"1.1.0"}}}`
	SampleValidScanResultsVulnerability = `{"type":"vulnerability","data":{"vulnerabilities":[{"id":316274974,"name":"hadoop","org":"apache","version":"2.8.0","up":null,"edition":null,"aliases":null,"created_at":"2017-02-13T20:02:32.785Z","updated_at":"2017-02-13T20:02:32.785Z","title":null,"references":null,"part":null,"language":null,"source_id":1,"external_id":"cpe:/a:apache:hadoop:2.8.0","vulnerabilities":[{"id":92596,"external_id":"CVE-2017-7669","title":"CVE-2017-7669","summary":"In Apache Hadoop 2.8.0, 3.0.0-alpha1, and 3.0.0-alpha2, the LinuxContainerExecutor runs docker commands as root with insufficient input validation. When the docker feature is enabled, authenticated users can run commands as root.","score":"8.5","score_version":"2.0","score_system":"CVSS","score_details":{"cvssv2":{"vectorString":"(AV:N/AC:M/Au:S/C:C/I:C/A:C)","accessVector":"NETWORK","accessComplexity":"MEDIUM","authentication":"SINGLE","confidentialityImpact":"COMPLETE","integrityImpact":"COMPLETE","availabilityImpact":"COMPLETE","baseScore":8.5},"cvssv3":{"vectorString":"AV:N/AC:H/PR:L/UI:N/S:U/C:H/I:H/A:H","attackVector":"NETWORK","attackComplexity":"HIGH","privilegesRequired":"LOW","userInteraction":"NONE","scope":"UNCHANGED","confidentialityImpact":"HIGH","integrityImpact":"HIGH","availabilityImpact":"HIGH","baseScore":7.5,"baseSeverity":"HIGH"}},"vector":"NETWORK","access_complexity":"MEDIUM","vulnerability_authentication":"SINGLE","confidentiality_impact":"COMPLETE","integrity_impact":"COMPLETE","availability_impact":"COMPLETE","vulnerability_source":null,"assessment_check":null,"scanner":null,"recommendation":"","references":[{"type":"UNKNOWN","source":"","url":"http://www.securityfocus.com/bid/98795","text":"http://www.securityfocus.com/bid/98795"},{"type":"UNKNOWN","source":"","url":"https://mail-archives.apache.org/mod_mbox/hadoop-user/201706.mbox/%3C4A2FDA56-491B-4C2A-915F-C9D4A4BDB92A%40apache.org%3E","text":"https://mail-archives.apache.org/mod_mbox/hadoop-user/201706.mbox/%3C4A2FDA56-491B-4C2A-915F-C9D4A4BDB92A%40apache.org%3E"}],"modified_at":"2017-06-09T16:21:00.000Z","published_at":"2017-06-05T01:29:00.000Z","created_at":"2017-07-12T23:07:35.491Z","updated_at":"2017-07-12T23:07:35.491Z","source_id":1}]}],"meta":{"vulnerability_count":1}}}`
	SampleInvalidResults                = `{"type":"fooresult", "data":"I pitty the foo"}`
	SampleValidScanResultsDifference    = `{"data": {"checksum": "checksumishere","difference": true},"type": "difference"}`
	SampleValidExternalVulnerabilities  = `{"type":"external_vulnerability","data":{"critical":1,"high":0,"medium":1,"low": 0}}`

	SampleValidUntranslatedScanResultsLicense = `{"license": {"license": {"type": [{"name": "a license"}], "name": "some license"}}}`
)
