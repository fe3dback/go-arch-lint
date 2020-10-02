package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

type (
	combinedPayload struct {
		Type    string
		Payload interface{}
	}
)

func output(payloadType outputPayloadType, payload interface{}, asciiFn func()) {
	payload = combinedPayload{
		Type:    payloadType,
		Payload: payload,
	}

	switch flagOutputType {
	case outputTypeJSON:
		var jsonBuffer []byte
		var err error

		if flagOutputJsonOneLine {
			jsonBuffer, err = json.Marshal(payload)
		} else {
			jsonBuffer, err = json.MarshalIndent(payload, "", "  ")
		}

		if err != nil {
			panic(fmt.Sprintf("failed to marshal payload '%v' to json: %s", payload, err))
		}

		fmt.Println(string(jsonBuffer))

	case outputTypeASCII:
		asciiFn()

	default:
		panic(fmt.Sprintf("unknown output type: %s", flagOutputType))
	}
}

func halt(err error) {
	if flagOutputType == outputTypeDefault {
		// cobra error before parsing args
		// for example "Error: unknown flag --example"

		// output in default ascii mode
		flagOutputType = outputTypeASCII
	}

	payload := payloadTypeHalt{Error: err.Error()}
	output(outputPayloadTypeHalt, payload, func() {
		if flagUseColors && au != nil {
			fmt.Printf("%s\n", au.Yellow(err.Error()))
		} else {
			fmt.Println(err.Error())
		}
	})
}

func getModuleNameFromGoModFile(goModPath string) (string, error) {
	gomod, err := parseGoMod(goModPath)
	if err != nil {
		return "", fmt.Errorf("can`t parse gomod: %v", err)
	}

	moduleName := gomod.Module.Mod.Path
	if moduleName == "" {
		return "", fmt.Errorf("%s should contain module name in 'module'", goModFileName)
	}

	return moduleName, nil
}

func parseGoMod(path string) (*modfile.File, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %v", path, err)
	}

	mod, err := modfile.ParseLax(path, file, nil)
	if err != nil {
		return nil, fmt.Errorf("modfile parseLax failed '%s': %v", path, err)
	}

	return mod, nil
}
