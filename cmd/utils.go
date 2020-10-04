package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fe3dback/go-arch-lint/models"
)

type (
	combinedPayload struct {
		Type    string
		Payload interface{}
	}
)

func mustFetchFlags(ctx context.Context) *rootInput {
	inp := ctx.Value(ctxRootInputFlags)
	switch flags := inp.(type) {
	case *rootInput:
		return flags
	default:
		panic(models.NewUserSpaceError("Failed to extract root command input flags from context"))
	}
}

func output(flags *rootInput, payloadType outputPayloadType, payload interface{}, rawOutputFn func()) {
	payload = combinedPayload{
		Type:    payloadType,
		Payload: payload,
	}

	switch flags.outputType {
	case outputTypeJSON:
		var jsonBuffer []byte
		var err error

		if flags.outputJsonOneLine {
			jsonBuffer, err = json.Marshal(payload)
		} else {
			jsonBuffer, err = json.MarshalIndent(payload, "", "  ")
		}

		if err != nil {
			panic(models.NewUserSpaceError(fmt.Sprintf("failed to marshal payload '%v' to json: %s", payload, err)))
		}

		fmt.Println(string(jsonBuffer))

	case outputTypeASCII:
		rawOutputFn()

	default:
		panic(models.NewUserSpaceError(fmt.Sprintf("unknown output type: %s", flags.outputType)))
	}
}

func halt(flags *rootInput, err error) {
	if flags.outputType == outputTypeDefault {
		// cobra error before parsing args
		// for example "Error: unknown flag --example"

		// output in default ascii mode
		flags.outputType = outputTypeASCII
	}

	payload := payloadTypeHalt{Error: err.Error()}
	output(flags, outputPayloadTypeHalt, payload, func() {
		if flags.useColors && flags.au != nil {
			fmt.Printf("%s\n", flags.au.Yellow(err.Error()))
		} else {
			fmt.Println(err.Error())
		}
	})
}
