package models

type Application struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone,omitempty"`
	Education      string `json:"education,omitempty"`
	Experience     string `json:"experience,omitempty"`
	Location       string `json:"location,omitempty"`
	CoverLetter    string `json:"coverLetter,omitempty"`
	JobID          string `json:"jobId,omitempty"`
	CVOriginalName string `json:"cvOriginalName"`
	CVStoredName   string `json:"cvStoredName"`
	CVMimeType     string `json:"cvMimeType"`
	CVSize         int64  `json:"cvSize"`
	AppliedAt      string `json:"appliedAt"`
	IP             string `json:"ip"`
}
