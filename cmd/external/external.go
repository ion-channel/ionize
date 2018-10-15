package external

//AnalysisID contains data fields that will identify a given analysis
type AnalysisID struct {
	ID        string
	TeamID    string
	ProjectID string
	APIKey    string
}

//NewAnalysisID Creates and returns a new AnalysisID struct
func NewAnalysisID(id, teamID, projectID, apiKey string) *AnalysisID {
	return &AnalysisID{
		ID:        id,
		TeamID:    teamID,
		ProjectID: projectID,
		APIKey:    apiKey,
	}
}
