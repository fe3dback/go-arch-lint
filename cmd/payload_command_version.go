package cmd

type payloadVersion struct {
	LinterVersion       string `json:"linter_version"`
	GoArchFileSupported string `json:"go_arch_file_supported"`
}
