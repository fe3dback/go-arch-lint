package models

type (
	SelfInspect struct {
		LinterVersion string       `json:"LinterVersion"`
		Notices       []Annotation `json:"Notices"`
		Suggestions   []Annotation `json:"Suggestions"`
	}

	Annotation struct {
		Text      string    `json:"Text"`
		Reference Reference `json:"Reference"`
	}
)
