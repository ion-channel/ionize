package external

func ParseVulnerabilities(contents string) (*Fortify, error) {
	return nil, nil
}

type Vulnerabilities struct {
	Value int
}

func (c *Vulnerabilities) Save() error {
	return nil
}
