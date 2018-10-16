package external

const (
	//Accuracy MetaInfo group name for accuracy component
	Accuracy = "Accuracy"
	//Impact MetaInfo group name for impact component
	Impact = "Impact"
	//Probability MetaInfo group name for probability component
	Probability = "Probability"
)

//Rule encapsulates the rule data from fortify
type Rule struct {
	MetaInfo struct {
		Group []struct {
			Name string `xml:"name,attr"`
			Text string `xml:",chardata"`
		} `xml:"Group"`
	} `xml:"MetaInfo"`
	ID string `xml:"id,attr"`
}

//RuleInfo encapsulates the rule info data from fortify
type RuleInfo struct {
	Rule []Rule `xml:"Rule"`
}

//ProgramData encapsulates the ProgramData from fortify
type ProgramData struct {
	Sources struct {
		SourceInstance []struct {
			SourceLocation struct {
				Path     string `xml:"path,attr"`
				Line     string `xml:"line,attr"`
				LineEnd  string `xml:"lineEnd,attr"`
				ColStart string `xml:"colStart,attr"`
				ColEnd   string `xml:"colEnd,attr"`
			} `xml:"SourceLocation,omitempty"`
			TaintFlags struct {
				TaintFlag []struct {
					Name string `xml:"name,attr"`
				} `xml:"TaintFlag"`
			} `xml:"TaintFlags"`
			RuleID       string `xml:"ruleID,attr"`
			FunctionCall struct {
				SourceLocation struct {
					Path     string `xml:"path,attr"`
					Line     string `xml:"line,attr"`
					LineEnd  string `xml:"lineEnd,attr"`
					ColStart string `xml:"colStart,attr"`
					ColEnd   string `xml:"colEnd,attr"`
				} `xml:"SourceLocation"`
				Function []struct {
					Name string `xml:"name,attr"`
				} `xml:"Function"`
			} `xml:"FunctionCall,omitempty"`
			FunctionEntry struct {
				SourceLocation struct {
					Path     string `xml:"path,attr"`
					Line     string `xml:"line,attr"`
					LineEnd  string `xml:"lineEnd,attr"`
					ColStart string `xml:"colStart,attr"`
					ColEnd   string `xml:"colEnd,attr"`
				} `xml:"SourceLocation"`
				Function []struct {
					Name string `xml:"name,attr"`
				} `xml:"Function"`
			} `xml:"FunctionEntry,omitempty"`
		} `xml:"SourceInstance"`
	} `xml:"Sources"`
	Sinks struct {
		SinkInstance []struct {
			FunctionCall struct {
				SourceLocation struct {
					Path     string `xml:"path,attr"`
					Line     string `xml:"line,attr"`
					LineEnd  string `xml:"lineEnd,attr"`
					ColStart string `xml:"colStart,attr"`
					ColEnd   string `xml:"colEnd,attr"`
				} `xml:"SourceLocation"`
				Function []struct {
					Name string `xml:"name,attr"`
				} `xml:"Function"`
			} `xml:"FunctionCall"`
			RuleID string `xml:"ruleID,attr"`
		} `xml:"SinkInstance"`
	} `xml:"Sinks"`
	CalledWithNoDef struct {
		Function []struct {
			Name string `xml:"name,attr"`
		} `xml:"Function"`
	} `xml:"CalledWithNoDef"`
}

//Snippets encapsulates the snippets from fortify
type Snippets struct {
	Snippet []struct {
		File      []string `xml:"File"`
		StartLine string   `xml:"StartLine"`
		EndLine   string   `xml:"EndLine"`
		Text      struct {
			Cdata string `xml:"_cdata,attr"`
		} `xml:"Text"`
		ID string `xml:"id,attr"`
	} `xml:"Snippet"`
}

//Description encapsulates the Description from fortify
type Description []struct {
	Abstract        string `xml:"Abstract"`
	Explanation     string `xml:"Explanation"`
	Recommendations string `xml:"Recommendations"`
	Tips            struct {
		Tip []string `xml:"Tip"`
	} `xml:"Tips,omitempty"`
	References struct {
		Reference []struct {
			Title     string `xml:"Title"`
			Author    string `xml:"Author,omitempty"`
			Source    string `xml:"Source,omitempty"`
			Publisher string `xml:"Publisher,omitempty"`
		} `xml:"Reference"`
	} `xml:"References"`
	ContentType string `xml:"contentType,attr"`
	ClassID     string `xml:"classID,attr"`
}

//UnifiedTracePool encapsulates the UnifiedTracePool from fortify
type UnifiedTracePool struct {
	Trace []struct {
		Primary struct {
			Entry []struct {
				Node []struct {
					SourceLocation struct {
						Path      string `xml:"path,attr"`
						Line      string `xml:"line,attr"`
						LineEnd   string `xml:"lineEnd,attr"`
						ColStart  string `xml:"colStart,attr"`
						ColEnd    string `xml:"colEnd,attr"`
						ContextID string `xml:"contextId,attr"`
						Snippet   string `xml:"snippet,attr"`
					} `xml:"SourceLocation"`
					Action struct {
						Type string `xml:"type,attr"`
						Text string `xml:"_text,attr"`
					} `xml:"Action"`
					Reason struct {
						TraceRef struct {
							ID string `xml:"id,attr"`
						} `xml:"TraceRef"`
					} `xml:"Reason"`
				} `xml:"Node"`
			} `xml:"Entry"`
		} `xml:"Primary"`
		ID string `xml:"id,attr"`
	} `xml:"Trace"`
}

//AnalysisInfo 411
type AnalysisInfo struct {
	Unified struct {
		Context []struct {
			Function []struct {
				Name string `xml:"name,attr"`
			} `xml:"Function"`
			FunctionDeclarationSourceLocation struct {
				Path     string `xml:"path,attr"`
				Line     string `xml:"line,attr"`
				LineEnd  string `xml:"lineEnd,attr"`
				ColStart string `xml:"colStart,attr"`
				ColEnd   string `xml:"colEnd,attr"`
			} `xml:"FunctionDeclarationSourceLocation"`
		} `xml:"Context"`
		ReplacementDefinitions struct {
			Def []struct {
				Key   string `xml:"key,attr"`
				Value string `xml:"value,attr"`
			} `xml:"Def"`
		} `xml:"ReplacementDefinitions"`
		Trace []struct {
			Primary struct {
				Entry []struct {
					Node []struct {
						SourceLocation struct {
							Path     string `xml:"path,attr"`
							Line     string `xml:"line,attr"`
							LineEnd  string `xml:"lineEnd,attr"`
							ColStart string `xml:"colStart,attr"`
							ColEnd   string `xml:"colEnd,attr"`
							Snippet  string `xml:"snippet,attr"`
						} `xml:"SourceLocation"`
						Action struct {
							Type string `xml:"type,attr"`
							Text string `xml:"_text,attr"`
						} `xml:"Action"`
						DetailsOnly string `xml:"detailsOnly,attr"`
						Label       string `xml:"label,attr"`
					} `xml:"Node"`
				} `xml:"Entry"`
			} `xml:"Primary"`
		} `xml:"Trace"`
	} `xml:"Unified"`
}

//UnifiedNodePool more pools
type UnifiedNodePool struct {
	Node []struct {
		SourceLocation struct {
			Path      string `xml:"path,attr"`
			Line      string `xml:"line,attr"`
			LineEnd   string `xml:"lineEnd,attr"`
			ColStart  string `xml:"colStart,attr"`
			ColEnd    string `xml:"colEnd,attr"`
			ContextID string `xml:"contextId,attr"`
			Snippet   string `xml:"snippet,attr"`
		} `xml:"SourceLocation"`
		Action struct {
			Type string `xml:"type,attr"`
			Text string `xml:"_text,attr"`
		} `xml:"Action"`
		Reason struct {
			Rule []struct {
				RuleID string `xml:"ruleID,attr"`
			} `xml:"Rule"`
		} `xml:"Reason,omitempty"`
		ID        string `xml:"id,attr"`
		Knowledge struct {
			Fact []struct {
				Primary string `xml:"primary,attr"`
				Type    string `xml:"type,attr"`
				Text    string `xml:"_text,attr"`
			} `xml:"Fact"`
		} `xml:"Knowledge,omitempty"`
		SecondaryLocation struct {
			Path     string `xml:"path,attr"`
			Line     string `xml:"line,attr"`
			LineEnd  string `xml:"lineEnd,attr"`
			ColStart string `xml:"colStart,attr"`
			ColEnd   string `xml:"colEnd,attr"`
			Snippet  string `xml:"snippet,attr"`
		} `xml:"SecondaryLocation,omitempty"`
	} `xml:"Node"`
}

//ContextPool jump in
type ContextPool struct {
	Context []struct {
		Function []struct {
			Name string `xml:"name,attr"`
		} `xml:"Function"`
		FunctionDeclarationSourceLocation struct {
			Path     string `xml:"path,attr"`
			Line     string `xml:"line,attr"`
			LineEnd  string `xml:"lineEnd,attr"`
			ColStart string `xml:"colStart,attr"`
			ColEnd   string `xml:"colEnd,attr"`
		} `xml:"FunctionDeclarationSourceLocation"`
		ID string `xml:"id,attr"`
	} `xml:"Context"`
}

//Build this is the build
type Build struct {
	BuildID     string `xml:"BuildID"`
	NumberFiles string `xml:"NumberFiles"`
	LOC         []struct {
		Type string `xml:"type,attr"`
		Text string `xml:"_text,attr"`
	} `xml:"LOC"`
	SourceBasePath string `xml:"SourceBasePath"`
	SourceFiles    struct {
		File []struct {
			Name string `xml:"Name"`
			LOC  []struct {
				Type string `xml:"type,attr"`
				Text string `xml:"_text,attr"`
			} `xml:"LOC"`
			Size      string `xml:"size,attr"`
			Timestamp string `xml:"timestamp,attr"`
			Loc       string `xml:"loc,attr,omitempty"`
			Type      string `xml:"type,attr"`
			Encoding  string `xml:"encoding,attr"`
		} `xml:"File"`
	} `xml:"SourceFiles"`
	ScanTime struct {
		Value string `xml:"value,attr"`
	} `xml:"ScanTime"`
}

//Vulnerability Fortify detects them
type Vulnerability struct {
	ClassInfo struct {
		ClassID         string `xml:"ClassID"`
		Kingdom         string `xml:"Kingdom"`
		Type            string `xml:"Type"`
		AnalyzerName    string `xml:"AnalyzerName"`
		DefaultSeverity string `xml:"DefaultSeverity"`
	} `xml:"ClassInfo"`
	InstanceInfo struct {
		InstanceID       string `xml:"InstanceID"`
		InstanceSeverity string `xml:"InstanceSeverity"`
		Confidence       string `xml:"Confidence"`
	} `xml:"InstanceInfo"`
	AnalysisInfo AnalysisInfo `xml:"AnalysisInfo"`
}

//EngineData encapsulates EngineData from fortify
type EngineData struct {
	EngineVersion   string `xml:"EngineVersion"`
	InactiveResults string `xml:"InactiveResults"`
	RulePacks       struct {
		RulePack []struct {
			RulePackID string `xml:"RulePackID"`
			SKU        string `xml:"SKU"`
			Name       string `xml:"Name"`
			Version    string `xml:"Version"`
			MAC        string `xml:"MAC"`
		} `xml:"RulePack"`
	} `xml:"RulePacks"`
	Properties []struct {
		Property []struct {
			Name  string `xml:"name"`
			Value string `xml:"value"`
		} `xml:"Property"`
		Type string `xml:"_type"`
	} `xml:"Properties"`
	CommandLine struct {
		Argument []string `xml:"Argument"`
	} `xml:"CommandLine"`
	Errors      string `xml:"Errors"`
	MachineInfo struct {
		Hostname string `xml:"Hostname"`
		Username string `xml:"Username"`
		Platform string `xml:"Platform"`
	} `xml:"MachineInfo"`
	FilterResult string   `xml:"FilterResult"`
	RuleInfo     RuleInfo `xml:"RuleInfo"`
	LicenseInfo  struct {
		Metadata []struct {
			Name  string `xml:"name"`
			Value string `xml:"value"`
		} `xml:"Metadata"`
		Capability []struct {
			Name       string `xml:"Name"`
			Expiration string `xml:"Expiration"`
			Attribute  struct {
				Name  string `xml:"name"`
				Value string `xml:"value"`
			} `xml:"Attribute,omitempty"`
		} `xml:"Capability"`
	} `xml:"LicenseInfo"`
}

//FVDL yes yes
type FVDL struct {
	UUID            string `xml:"UUID"`
	Build           Build  `xml:"Build"`
	Vulnerabilities struct {
		Vulnerability []Vulnerability `xml:"Vulnerability"`
	} `xml:"Vulnerabilities"`
	ContextPool      ContextPool      `xml:"ContextPool"`
	UnifiedNodePool  UnifiedNodePool  `xml:"UnifiedNodePool"`
	UnifiedTracePool UnifiedTracePool `xml:"UnifiedTracePool"`
	Description      Description      `xml:"Description"`
	Snippets         Snippets         `xml:"Snippets"`
	ProgramData      ProgramData      `xml:"ProgramData"`
	EngineData       EngineData       `xml:"EngineData"`
}

//Rules returns the rules used in the Fortify file
func (f *FVDL) Rules() map[string]Rule {
	rules := make(map[string]Rule)
	for _, rsf := range f.EngineData.RuleInfo.Rule {
		rules[rsf.ID] = rsf
	}
	return rules
}

//Group returns the value of the metainfo group of a rule
func (f *FVDL) Group(ruleID, groupName string) string {
	rules := f.Rules()
	rule := rules[ruleID]
	var value string
	for _, g := range rule.MetaInfo.Group {
		if g.Name == groupName {
			value = g.Text
		}
	}

	return value
}
