package digests

import (
	"encoding/json"
	"testing"

	"github.com/ion-channel/ionic/scanner"
	"github.com/ion-channel/ionic/scans"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestVulnerabilitiesDigests(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Vulnerabilities", func() {
		g.It("should produce digests", func() {
			s := &scanner.ScanStatus{}
			e := scans.NewEval()

			var r scans.VulnerabilityResults
			b := []byte(validHighCriticalVulns)
			err := json.Unmarshal(b, &r)
			Expect(err).To(BeNil())

			e.TranslatedResults = &scans.TranslatedResults{
				Type: "vulnerability",
				Data: r,
			}

			ds, err := vulnerabilityDigests(s, e)
			Expect(err).To(BeNil())
			Expect(len(ds)).To(Equal(4))

			Expect(ds[0].Title).To(Equal("total vulnerabilities"))
			Expect(string(ds[0].Data)).To(Equal(`{"count":3}`))
			Expect(ds[0].Warning).To(BeTrue())
			Expect(ds[0].Pending).To(BeFalse())
			Expect(ds[0].Errored).To(BeFalse())

			Expect(ds[1].Title).To(Equal("unique vulnerabilities"))
			Expect(string(ds[1].Data)).To(Equal(`{"count":3}`))
			Expect(ds[1].Warning).To(BeTrue())
			Expect(ds[1].Pending).To(BeFalse())
			Expect(ds[1].Errored).To(BeFalse())

			Expect(ds[2].Title).To(Equal("high vulnerabilities"))
			Expect(string(ds[2].Data)).To(Equal(`{"count":2}`))
			Expect(ds[2].Pending).To(BeFalse())
			Expect(ds[2].Errored).To(BeFalse())

			Expect(ds[3].Title).To(Equal("critical vulnerability"))
			Expect(string(ds[3].Data)).To(Equal(`{"count":1}`))
			Expect(ds[3].Pending).To(BeFalse())
			Expect(ds[3].Errored).To(BeFalse())
		})

		g.It("should differenciate value counts", func() {
			s := &scanner.ScanStatus{}
			e := scans.NewEval()

			var r scans.VulnerabilityResults
			b := []byte(validHighVulns)
			err := json.Unmarshal(b, &r)
			Expect(err).To(BeNil())

			e.TranslatedResults = &scans.TranslatedResults{
				Type: "vulnerability",
				Data: r,
			}

			ds, err := vulnerabilityDigests(s, e)
			Expect(err).To(BeNil())
			Expect(len(ds)).To(Equal(4))

			Expect(ds[0].Title).To(Equal("total vulnerabilities"))
			Expect(string(ds[0].Data)).To(Equal(`{"count":4}`))
			Expect(ds[0].Warning).To(BeTrue())
			Expect(ds[0].Pending).To(BeFalse())
			Expect(ds[0].Errored).To(BeFalse())

			Expect(ds[1].Title).To(Equal("unique vulnerabilities"))
			Expect(string(ds[1].Data)).To(Equal(`{"count":4}`))
			Expect(ds[1].Warning).To(BeTrue())
			Expect(ds[1].Pending).To(BeFalse())
			Expect(ds[1].Errored).To(BeFalse())

			Expect(ds[2].Title).To(Equal("high vulnerabilities"))
			Expect(string(ds[2].Data)).To(Equal(`{"count":3}`))
			Expect(ds[2].Passed).To(BeFalse())
			Expect(ds[2].Pending).To(BeFalse())
			Expect(ds[2].Errored).To(BeFalse())

			Expect(ds[3].Title).To(Equal("critical vulnerabilities"))
			Expect(string(ds[3].Data)).To(Equal(`{"count":0}`))
			Expect(ds[3].Passed).To(BeTrue())
			Expect(ds[3].Pending).To(BeFalse())
			Expect(ds[3].Errored).To(BeFalse())
		})
	})
}

const (
	validHighCriticalVulns = `{"vulnerabilities":[{"id":92400,"external_id":"cpe:/a:rack_project:rack:1.1.0","source_id":0,"title":"Rack_project Rack 1.1.0","name":"rack","org":"rack_project","version":"1.1.0","up":"","edition":"","aliases":null,"created_at":"2017-02-13T20:02:43.83Z","updated_at":"2018-05-25T04:10:15.382Z","references":[],"part":"/a","language":"","vulnerabilities":[{"id":267937758,"external_id":"CVE-2011-5036","source":[{"id":1,"name":"NVD","description":"National Vulnerability Database","created_at":"2017-02-09T20:18:35.385Z","updated_at":"2017-02-13T20:12:05.342Z","attribution":"Copyright © 1999–2017, The MITRE Corporation. CVE and the CVE logo are registered trademarks and CVE-Compatible is a trademark of The MITRE Corporation.","license":"Submissions: For all materials you submit to the Common Vulnerabilities and Exposures (CVE®), you hereby grant to The MITRE Corporation (MITRE) and all CVE Numbering Authorities (CNAs) a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute such materials and derivative works. Unless required by applicable law or agreed to in writing, you provide such materials on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.\n\nCVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for such purposes is authorized provided that you reproduce MITRE's copyright designation and this license in any such copy.\n","copyright_url":"http://cve.mitre.org/about/termsofuse.html"}],"title":"CVE-2011-5036","summary":"Rack before 1.1.3, 1.2.x before 1.2.5, and 1.3.x before 1.3.6 computes hash values for form parameters without restricting the ability to trigger hash collisions predictably, which allows remote attackers to cause a denial of service (CPU consumption) by sending many crafted parameters.","score":"7.5","score_version":"3.0","score_system":"CVSS","score_details":{"cvssv3":{"vectorString":"AV:N/AC:L/Au:N/C:N/I:N/A:P","accessVector":"NETWORK","accessComplexity":"LOW","authentication":"NONE","confidentialityImpact":"NONE","integrityImpact":"NONE","availabilityImpact":"PARTIAL","baseScore":7.5}},"vector":"","access_complexity":"","vulnerability_authentication":"","confidentiality_impact":"","integrity_impact":"","availability_impact":"","vulnerabilty_source":"","assessment_check":null,"scanner":null,"recommendation":"","dependencies":null,"references":[{"type":"UNKNOWN","source":"20111228 n.runs-SA-2011.004 - web programming languages and platforms - DoS through hash table","url":"http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html","text":"http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html"},{"type":"UNKNOWN","source":"DSA-2783","url":"http://www.debian.org/security/2013/dsa-2783","text":"http://www.debian.org/security/2013/dsa-2783"},{"type":"UNKNOWN","source":"VU#903934","url":"http://www.kb.cert.org/vuls/id/903934","text":"http://www.kb.cert.org/vuls/id/903934"},{"type":"UNKNOWN","source":"http://www.nruns.com/_downloads/advisory28122011.pdf","url":"http://www.nruns.com/_downloads/advisory28122011.pdf","text":"http://www.nruns.com/_downloads/advisory28122011.pdf"},{"type":"UNKNOWN","source":"http://www.ocert.org/advisories/ocert-2011-003.html","url":"http://www.ocert.org/advisories/ocert-2011-003.html","text":"http://www.ocert.org/advisories/ocert-2011-003.html"},{"type":"UNKNOWN","source":"https://gist.github.com/52bbc6b9cc19ce330829","url":"https://gist.github.com/52bbc6b9cc19ce330829","text":"https://gist.github.com/52bbc6b9cc19ce330829"}],"modified_at":"2013-10-31T03:21:00Z","published_at":"2011-12-30T01:55:00Z","created_at":"2018-03-10T22:53:07.473Z","updated_at":"2018-11-24T10:05:23.28Z"},{"id":267868924,"external_id":"CVE-2013-0184","source":[{"id":1,"name":"NVD","description":"National Vulnerability Database","created_at":"2017-02-09T20:18:35.385Z","updated_at":"2017-02-13T20:12:05.342Z","attribution":"Copyright © 1999–2017, The MITRE Corporation. CVE and the CVE logo are registered trademarks and CVE-Compatible is a trademark of The MITRE Corporation.","license":"Submissions: For all materials you submit to the Common Vulnerabilities and Exposures (CVE®), you hereby grant to The MITRE Corporation (MITRE) and all CVE Numbering Authorities (CNAs) a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute such materials and derivative works. Unless required by applicable law or agreed to in writing, you provide such materials on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.\n\nCVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for such purposes is authorized provided that you reproduce MITRE's copyright designation and this license in any such copy.\n","copyright_url":"http://cve.mitre.org/about/termsofuse.html"}],"title":"CVE-2013-0184","summary":"Unspecified vulnerability in Rack::Auth::AbstractRequest in Rack 1.1.x before 1.1.5, 1.2.x before 1.2.7, 1.3.x before 1.3.9, and 1.4.x before 1.4.4 allows remote attackers to cause a denial of service via unknown vectors related to \"symbolized arbitrary strings.\"","score":"9.5","score_version":"3.0","score_system":"CVSS","score_details":{"cvssv3":{"vectorString":"AV:N/AC:M/Au:N/C:N/I:N/A:P","accessVector":"NETWORK","accessComplexity":"MEDIUM","authentication":"NONE","confidentialityImpact":"NONE","integrityImpact":"NONE","availabilityImpact":"PARTIAL","baseScore":9.5}},"vector":"","access_complexity":"","vulnerability_authentication":"","confidentiality_impact":"","integrity_impact":"","availability_impact":"","vulnerabilty_source":"","assessment_check":null,"scanner":null,"recommendation":"","dependencies":null,"references":[{"type":"UNKNOWN","source":"openSUSE-SU-2013:0462","url":"http://lists.opensuse.org/opensuse-updates/2013-03/msg00048.html","text":"http://lists.opensuse.org/opensuse-updates/2013-03/msg00048.html"},{"type":"UNKNOWN","source":"RHSA-2013:0544","url":"http://rhn.redhat.com/errata/RHSA-2013-0544.html","text":"http://rhn.redhat.com/errata/RHSA-2013-0544.html"},{"type":"UNKNOWN","source":"RHSA-2013:0548","url":"http://rhn.redhat.com/errata/RHSA-2013-0548.html","text":"http://rhn.redhat.com/errata/RHSA-2013-0548.html"},{"type":"UNKNOWN","source":"DSA-2783","url":"http://www.debian.org/security/2013/dsa-2783","text":"http://www.debian.org/security/2013/dsa-2783"},{"type":"UNKNOWN","source":"https://bugzilla.redhat.com/show_bug.cgi?id=895384","url":"https://bugzilla.redhat.com/show_bug.cgi?id=895384","text":"https://bugzilla.redhat.com/show_bug.cgi?id=895384"}],"modified_at":"2013-10-31T03:30:00Z","published_at":"2013-03-01T05:40:00Z","created_at":"2018-03-09T00:48:08.313Z","updated_at":"2019-01-29T08:54:59.385Z"},{"id":627395877,"external_id":"CVE-2012-9186","source":[{"id":1,"name":"NVD","description":"National Vulnerability Database","created_at":"2017-02-09T20:18:35.385Z","updated_at":"2017-02-13T20:12:05.342Z","attribution":"Copyright © 1999–2017, The MITRE Corporation. CVE and the CVE logo are registered trademarks and CVE-Compatible is a trademark of The MITRE Corporation.","license":"Submissions: For all materials you submit to the Common Vulnerabilities and Exposures (CVE®), you hereby grant to The MITRE Corporation (MITRE) and all CVE Numbering Authorities (CNAs) a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute such materials and derivative works. Unless required by applicable law or agreed to in writing, you provide such materials on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.\n\nCVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for such purposes is authorized provided that you reproduce MITRE's copyright designation and this license in any such copy.\n","copyright_url":"http://cve.mitre.org/about/termsofuse.html"}],"title":"CVE-2012-9186","summary":"Rack before 1.1.3, 1.2.x before 1.2.5, and 1.3.x before 1.3.6 computes hash values for form parameters without restricting the ability to trigger hash collisions predictably, which allows remote attackers to cause a denial of service (CPU consumption) by sending many crafted parameters.","score":"7.5","score_version":"2.0","score_system":"CVSS","score_details":{"cvssv2":{"vectorString":"AV:N/AC:L/Au:N/C:N/I:N/A:P","accessVector":"NETWORK","accessComplexity":"LOW","authentication":"NONE","confidentialityImpact":"NONE","integrityImpact":"NONE","availabilityImpact":"PARTIAL","baseScore":7.5}},"vector":"","access_complexity":"","vulnerability_authentication":"","confidentiality_impact":"","integrity_impact":"","availability_impact":"","vulnerabilty_source":"","assessment_check":null,"scanner":null,"recommendation":"","dependencies":null,"references":[{"type":"UNKNOWN","source":"20111228 n.runs-SA-2011.004 - web programming languages and platforms - DoS through hash table","url":"http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html","text":"http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html"},{"type":"UNKNOWN","source":"DSA-2783","url":"http://www.debian.org/security/2013/dsa-2783","text":"http://www.debian.org/security/2013/dsa-2783"},{"type":"UNKNOWN","source":"VU#903934","url":"http://www.kb.cert.org/vuls/id/903934","text":"http://www.kb.cert.org/vuls/id/903934"},{"type":"UNKNOWN","source":"http://www.nruns.com/_downloads/advisory28122011.pdf","url":"http://www.nruns.com/_downloads/advisory28122011.pdf","text":"http://www.nruns.com/_downloads/advisory28122011.pdf"},{"type":"UNKNOWN","source":"http://www.ocert.org/advisories/ocert-2011-003.html","url":"http://www.ocert.org/advisories/ocert-2011-003.html","text":"http://www.ocert.org/advisories/ocert-2011-003.html"},{"type":"UNKNOWN","source":"https://gist.github.com/52bbc6b9cc19ce330829","url":"https://gist.github.com/52bbc6b9cc19ce330829","text":"https://gist.github.com/52bbc6b9cc19ce330829"}],"modified_at":"2013-10-31T03:21:00Z","published_at":"2011-12-30T01:55:00Z","created_at":"2018-03-10T22:53:07.473Z","updated_at":"2018-11-24T10:05:23.28Z"}]}],"meta":{"vulnerability_count":3}}`

	validHighVulns = `{
		"vulnerabilities": [
			{
				"id": 92400,
				"external_id": "cpe:/a:rack_project:rack:1.1.0",
				"source_id": 0,
				"title": "Rack_project Rack 1.1.0",
				"name": "rack",
				"org": "rack_project",
				"version": "1.1.0",
				"up": "",
				"edition": "",
				"aliases": null,
				"created_at": "2017-02-13T20:02:43.83Z",
				"updated_at": "2018-05-25T04:10:15.382Z",
				"references": [],
				"part": "/a",
				"language": "",
				"vulnerabilities": [
					{
						"id": 267937758,
						"external_id": "CVE-2011-5036",
						"source": [
							{
								"id": 1,
								"name": "NVD",
								"description": "National Vulnerability Database",
								"created_at": "2017-02-09T20:18:35.385Z",
								"updated_at": "2017-02-13T20:12:05.342Z",
								"attribution": "Copyright © 1999–2017, The MITRE Corporation. CVE and the CVE logo are registered trademarks and CVE-Compatible is a trademark of The MITRE Corporation.",
								"license": "Submissions: For all materials you submit to the Common Vulnerabilities and Exposures (CVE®), you hereby grant to The MITRE Corporation (MITRE) and all CVE Numbering Authorities (CNAs) a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute such materials and derivative works. Unless required by applicable law or agreed to in writing, you provide such materials on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.\n\nCVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for such purposes is authorized provided that you reproduce MITRE's copyright designation and this license in any such copy.\n",
								"copyright_url": "http://cve.mitre.org/about/termsofuse.html"
							}
						],
						"title": "CVE-2011-5036",
						"summary": "Rack before 1.1.3, 1.2.x before 1.2.5, and 1.3.x before 1.3.6 computes hash values for form parameters without restricting the ability to trigger hash collisions predictably, which allows remote attackers to cause a denial of service (CPU consumption) by sending many crafted parameters.",
						"score": "7.5",
						"score_version": "3.0",
						"score_system": "CVSS",
						"score_details": {
							"cvssv3": {
								"vectorString": "AV:N/AC:L/Au:N/C:N/I:N/A:P",
								"accessVector": "NETWORK",
								"accessComplexity": "LOW",
								"authentication": "NONE",
								"confidentialityImpact": "NONE",
								"integrityImpact": "NONE",
								"availabilityImpact": "PARTIAL",
								"baseScore": 7.5
							}
						},
						"vector": "",
						"access_complexity": "",
						"vulnerability_authentication": "",
						"confidentiality_impact": "",
						"integrity_impact": "",
						"availability_impact": "",
						"vulnerabilty_source": "",
						"assessment_check": null,
						"scanner": null,
						"recommendation": "",
						"dependencies": null,
						"references": [
							{
								"type": "UNKNOWN",
								"source": "20111228 n.runs-SA-2011.004 - web programming languages and platforms - DoS through hash table",
								"url": "http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html",
								"text": "http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html"
							},
							{
								"type": "UNKNOWN",
								"source": "DSA-2783",
								"url": "http://www.debian.org/security/2013/dsa-2783",
								"text": "http://www.debian.org/security/2013/dsa-2783"
							},
							{
								"type": "UNKNOWN",
								"source": "VU#903934",
								"url": "http://www.kb.cert.org/vuls/id/903934",
								"text": "http://www.kb.cert.org/vuls/id/903934"
							},
							{
								"type": "UNKNOWN",
								"source": "http://www.nruns.com/_downloads/advisory28122011.pdf",
								"url": "http://www.nruns.com/_downloads/advisory28122011.pdf",
								"text": "http://www.nruns.com/_downloads/advisory28122011.pdf"
							},
							{
								"type": "UNKNOWN",
								"source": "http://www.ocert.org/advisories/ocert-2011-003.html",
								"url": "http://www.ocert.org/advisories/ocert-2011-003.html",
								"text": "http://www.ocert.org/advisories/ocert-2011-003.html"
							},
							{
								"type": "UNKNOWN",
								"source": "https://gist.github.com/52bbc6b9cc19ce330829",
								"url": "https://gist.github.com/52bbc6b9cc19ce330829",
								"text": "https://gist.github.com/52bbc6b9cc19ce330829"
							}
						],
						"modified_at": "2013-10-31T03:21:00Z",
						"published_at": "2011-12-30T01:55:00Z",
						"created_at": "2018-03-10T22:53:07.473Z",
						"updated_at": "2018-11-24T10:05:23.28Z"
					},
					{
						"id": 268058529,
						"external_id": "NPM-118",
						"source": [
							{
								"id": 3,
								"name": "NPM",
								"description": "NPM Security advisories",
								"created_at": "2018-08-15T22:26:13.139Z",
								"updated_at": "2018-08-15T22:26:13.139Z",
								"attribution": "NPM Term and Licenses https://www.npmjs.com/policies/terms",
								"license": "https://www.npmjs.com/policies/open-source-terms",
								"copyright_url": "https://www.npmjs.com/policies/dmca"
							}
						],
						"title": "Regular Expression Denial of Service",
						"summary": "Affected+versions+of+%60minimatch%60+are+vulnerable+to+regular+expression+denial+of+service+attacks+when+user+input+is+passed+into+the+%60pattern%60+argument+of+%60minimatch%28path%2C+pattern%29%60.%0A%0A%0A%23%23+Proof+of+Concept%0A%60%60%60%0Avar+minimatch+%3D+require%28%E2%80%9Cminimatch%E2%80%9D%29%3B%0A%0A%2F%2F+utility+function+for+generating+long+strings%0Avar+genstr+%3D+function+%28len%2C+chr%29+%7B%0A++var+result+%3D+%E2%80%9C%E2%80%9D%3B%0A++for+%28i%3D0%3B+i%3C%3Dlen%3B+i%2B%2B%29+%7B%0A++++result+%3D+result+%2B+chr%3B%0A++%7D%0A++return+result%3B%0A%7D%0A%0Avar+exploit+%3D+%E2%80%9C%5B%21%E2%80%9D+%2B+genstr%281000000%2C+%E2%80%9C%5C%5C%E2%80%9D%29+%2B+%E2%80%9CA%E2%80%9D%3B%0A%0A%2F%2F+minimatch+exploit.%0Aconsole.log%28%E2%80%9Cstarting+minimatch%E2%80%9D%29%3B%0Aminimatch%28%E2%80%9Cfoo%E2%80%9D%2C+exploit%29%3B%0Aconsole.log%28%E2%80%9Cfinishing+minimatch%E2%80%9D%29%3B%0A%60%60%60",
						"score": "7.0",
						"score_system": "NPM",
						"score_details": {},
						"vector": "",
						"access_complexity": "",
						"vulnerability_authentication": "",
						"confidentiality_impact": "",
						"integrity_impact": "",
						"availability_impact": "",
						"vulnerabilty_source": "",
						"assessment_check": null,
						"scanner": null,
						"recommendation": "Update to version 3.0.2 or later.",
						"dependencies": null,
						"references": [
							{
								"type": "url",
								"source": "NPM",
								"url": "https://npmjs.com/advisories/118",
								"text": ""
							},
							{
								"type": "reporter",
								"source": "NPM",
								"url": "",
								"text": "Nick Starke"
							},
							{
								"type": "CVE",
								"source": "NVD",
								"url": "",
								"text": "CVE-2016-10540"
							}
						],
						"modified_at": "2018-03-01T21:58:01.072Z",
						"published_at": "2016-05-25T16:37:20Z",
						"created_at": "2019-09-24T00:04:04.659Z",
						"updated_at": "2020-03-05T20:01:55.557Z"
					},
					{
						"id": 267868924,
						"external_id": "CVE-2013-0184",
						"source": [
							{
								"id": 1,
								"name": "NVD",
								"description": "National Vulnerability Database",
								"created_at": "2017-02-09T20:18:35.385Z",
								"updated_at": "2017-02-13T20:12:05.342Z",
								"attribution": "Copyright © 1999–2017, The MITRE Corporation. CVE and the CVE logo are registered trademarks and CVE-Compatible is a trademark of The MITRE Corporation.",
								"license": "Submissions: For all materials you submit to the Common Vulnerabilities and Exposures (CVE®), you hereby grant to The MITRE Corporation (MITRE) and all CVE Numbering Authorities (CNAs) a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute such materials and derivative works. Unless required by applicable law or agreed to in writing, you provide such materials on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.\n\nCVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for such purposes is authorized provided that you reproduce MITRE's copyright designation and this license in any such copy.\n",
								"copyright_url": "http://cve.mitre.org/about/termsofuse.html"
							}
						],
						"title": "CVE-2013-0184",
						"summary": "Unspecified vulnerability in Rack::Auth::AbstractRequest in Rack 1.1.x before 1.1.5, 1.2.x before 1.2.7, 1.3.x before 1.3.9, and 1.4.x before 1.4.4 allows remote attackers to cause a denial of service via unknown vectors related to \"symbolized arbitrary strings.\"",
						"score": "3.2",
						"score_version": "3.0",
						"score_system": "CVSS",
						"score_details": {
							"cvssv3": {
								"vectorString": "AV:N/AC:M/Au:N/C:N/I:N/A:P",
								"accessVector": "NETWORK",
								"accessComplexity": "MEDIUM",
								"authentication": "NONE",
								"confidentialityImpact": "NONE",
								"integrityImpact": "NONE",
								"availabilityImpact": "PARTIAL",
								"baseScore": 3.2
							}
						},
						"vector": "",
						"access_complexity": "",
						"vulnerability_authentication": "",
						"confidentiality_impact": "",
						"integrity_impact": "",
						"availability_impact": "",
						"vulnerabilty_source": "",
						"assessment_check": null,
						"scanner": null,
						"recommendation": "",
						"dependencies": null,
						"references": [
							{
								"type": "UNKNOWN",
								"source": "openSUSE-SU-2013:0462",
								"url": "http://lists.opensuse.org/opensuse-updates/2013-03/msg00048.html",
								"text": "http://lists.opensuse.org/opensuse-updates/2013-03/msg00048.html"
							},
							{
								"type": "UNKNOWN",
								"source": "RHSA-2013:0544",
								"url": "http://rhn.redhat.com/errata/RHSA-2013-0544.html",
								"text": "http://rhn.redhat.com/errata/RHSA-2013-0544.html"
							},
							{
								"type": "UNKNOWN",
								"source": "RHSA-2013:0548",
								"url": "http://rhn.redhat.com/errata/RHSA-2013-0548.html",
								"text": "http://rhn.redhat.com/errata/RHSA-2013-0548.html"
							},
							{
								"type": "UNKNOWN",
								"source": "DSA-2783",
								"url": "http://www.debian.org/security/2013/dsa-2783",
								"text": "http://www.debian.org/security/2013/dsa-2783"
							},
							{
								"type": "UNKNOWN",
								"source": "https://bugzilla.redhat.com/show_bug.cgi?id=895384",
								"url": "https://bugzilla.redhat.com/show_bug.cgi?id=895384",
								"text": "https://bugzilla.redhat.com/show_bug.cgi?id=895384"
							}
						],
						"modified_at": "2013-10-31T03:30:00Z",
						"published_at": "2013-03-01T05:40:00Z",
						"created_at": "2018-03-09T00:48:08.313Z",
						"updated_at": "2019-01-29T08:54:59.385Z"
					},
					{
						"id": 627395877,
						"external_id": "CVE-2012-9186",
						"source": [
							{
								"id": 1,
								"name": "NVD",
								"description": "National Vulnerability Database",
								"created_at": "2017-02-09T20:18:35.385Z",
								"updated_at": "2017-02-13T20:12:05.342Z",
								"attribution": "Copyright © 1999–2017, The MITRE Corporation. CVE and the CVE logo are registered trademarks and CVE-Compatible is a trademark of The MITRE Corporation.",
								"license": "Submissions: For all materials you submit to the Common Vulnerabilities and Exposures (CVE®), you hereby grant to The MITRE Corporation (MITRE) and all CVE Numbering Authorities (CNAs) a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute such materials and derivative works. Unless required by applicable law or agreed to in writing, you provide such materials on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied, including, without limitation, any warranties or conditions of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A PARTICULAR PURPOSE.\n\nCVE Usage: MITRE hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for such purposes is authorized provided that you reproduce MITRE's copyright designation and this license in any such copy.\n",
								"copyright_url": "http://cve.mitre.org/about/termsofuse.html"
							}
						],
						"title": "CVE-2012-9186",
						"summary": "Rack before 1.1.3, 1.2.x before 1.2.5, and 1.3.x before 1.3.6 computes hash values for form parameters without restricting the ability to trigger hash collisions predictably, which allows remote attackers to cause a denial of service (CPU consumption) by sending many crafted parameters.",
						"score": "7.5",
						"score_version": "2.0",
						"score_system": "CVSS",
						"score_details": {
							"cvssv2": {
								"vectorString": "AV:N/AC:L/Au:N/C:N/I:N/A:P",
								"accessVector": "NETWORK",
								"accessComplexity": "LOW",
								"authentication": "NONE",
								"confidentialityImpact": "NONE",
								"integrityImpact": "NONE",
								"availabilityImpact": "PARTIAL",
								"baseScore": 7.5
							}
						},
						"vector": "",
						"access_complexity": "",
						"vulnerability_authentication": "",
						"confidentiality_impact": "",
						"integrity_impact": "",
						"availability_impact": "",
						"vulnerabilty_source": "",
						"assessment_check": null,
						"scanner": null,
						"recommendation": "",
						"dependencies": null,
						"references": [
							{
								"type": "UNKNOWN",
								"source": "20111228 n.runs-SA-2011.004 - web programming languages and platforms - DoS through hash table",
								"url": "http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html",
								"text": "http://archives.neohapsis.com/archives/bugtraq/2011-12/0181.html"
							},
							{
								"type": "UNKNOWN",
								"source": "DSA-2783",
								"url": "http://www.debian.org/security/2013/dsa-2783",
								"text": "http://www.debian.org/security/2013/dsa-2783"
							},
							{
								"type": "UNKNOWN",
								"source": "VU#903934",
								"url": "http://www.kb.cert.org/vuls/id/903934",
								"text": "http://www.kb.cert.org/vuls/id/903934"
							},
							{
								"type": "UNKNOWN",
								"source": "http://www.nruns.com/_downloads/advisory28122011.pdf",
								"url": "http://www.nruns.com/_downloads/advisory28122011.pdf",
								"text": "http://www.nruns.com/_downloads/advisory28122011.pdf"
							},
							{
								"type": "UNKNOWN",
								"source": "http://www.ocert.org/advisories/ocert-2011-003.html",
								"url": "http://www.ocert.org/advisories/ocert-2011-003.html",
								"text": "http://www.ocert.org/advisories/ocert-2011-003.html"
							},
							{
								"type": "UNKNOWN",
								"source": "https://gist.github.com/52bbc6b9cc19ce330829",
								"url": "https://gist.github.com/52bbc6b9cc19ce330829",
								"text": "https://gist.github.com/52bbc6b9cc19ce330829"
							}
						],
						"modified_at": "2013-10-31T03:21:00Z",
						"published_at": "2011-12-30T01:55:00Z",
						"created_at": "2018-03-10T22:53:07.473Z",
						"updated_at": "2018-11-24T10:05:23.28Z"
					}
				]
			}
		],
		"meta": {
			"vulnerability_count": 4
		}
	}`
)
