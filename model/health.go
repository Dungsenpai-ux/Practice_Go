package model

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}