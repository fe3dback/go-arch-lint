package cmd

type outputPayloadType = string

const (
	outputPayloadTypeHalt           outputPayloadType = "halt"
	outputPayloadTypeCommandVersion outputPayloadType = "command.version"
	outputPayloadTypeCommandCheck   outputPayloadType = "command.check"
)
