package external

type AnalysisID struct {
	ID        string
	TeamID    string
	ProjectID string
	APIKey    string
}

func NewAnalysisID(id, teamID, projectID, apiKey string) *AnalysisID {
	return &AnalysisID{
		ID:        id,
		TeamID:    teamID,
		ProjectID: projectID,
		APIKey:    apiKey,
	}
}
