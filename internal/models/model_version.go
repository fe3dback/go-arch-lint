package models

type Version struct {
	LinterVersion       string `json:"LinterVersion"`
	GoArchFileSupported string `json:"GoArchFileSupported"`
	BuildTime           string `json:"BuildTime"`
	CommitHash          string `json:"CommitHash"`
}
