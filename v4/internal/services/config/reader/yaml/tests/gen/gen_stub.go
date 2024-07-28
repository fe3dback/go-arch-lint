package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/reader/yaml"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/reader/yaml/tests"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("could not determine filename")
	}

	currentDir := filepath.Dir(filename)
	stubDir := filepath.Clean(fmt.Sprintf("%s/../data", currentDir))

	fmt.Printf("working at %s\n", stubDir)

	entries, err := os.ReadDir(stubDir)
	if err != nil {
		panic(fmt.Sprintf("failed read dir %s: %v", stubDir, err))
	}

	reader := yaml.NewReader()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".yml" {
			continue
		}

		sourceCode, err := os.ReadFile(fmt.Sprintf("%s/%s", stubDir, entry.Name()))
		if err != nil {
			panic(fmt.Sprintf("failed read file %s: %v", entry.Name(), err))
		}

		conf, err := reader.Parse("/conf.yml", sourceCode)
		if err != nil {
			panic(fmt.Sprintf("failed read YML config %s: %v", entry.Name(), err))
		}

		nameParts := strings.Split(entry.Name(), ".")
		stubFilePath := fmt.Sprintf("%s/%s_parsed.gold", stubDir, nameParts[0])
		err = os.WriteFile(stubFilePath, []byte(tests.Dump(conf)), 0600)
		if err != nil {
			panic(fmt.Sprintf("failed write JSON config %s: %v", entry.Name(), err))
		}
	}
}
