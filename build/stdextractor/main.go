package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"golang.org/x/tools/go/packages"
)

type (
	args struct {
		outputPath    string
		goVersionInfo string
	}
)

func main() {
	args := args{}

	flag.StringVar(&args.outputPath, "output", "", "specify path to save generated file")
	flag.StringVar(&args.goVersionInfo, "goVersion", "", "go env version info")
	flag.Parse()

	err := process(args)
	if err != nil {
		fmt.Println(fmt.Errorf("err: %w", err))
		os.Exit(1)
	}

	os.Exit(0)
}

func process(args args) error {
	if args.outputPath == "" {
		return fmt.Errorf("specify 'output' flag for writing go code")
	}

	pkgs, err := extractStdPackages()
	if err != nil {
		return fmt.Errorf("failed extract std packages: %w", err)
	}

	outputDto := assembleOutputDto(args, pkgs)
	bytes, err := renderOutput(outputDto)
	if err != nil {
		return fmt.Errorf("failed render output: %w", err)
	}

	err = ioutil.WriteFile(args.outputPath, bytes, 0755)
	if err != nil {
		return fmt.Errorf("failed write output to '%s': %w", args.outputPath, err)
	}

	return nil
}

func assembleOutputDto(args args, stdPackages []string) outputDto {
	packageName := path.Base(path.Dir(path.Clean(args.outputPath)))

	return outputDto{
		GoVersionInfo: args.goVersionInfo,
		GeneratedAt:   time.Now().Format(time.RFC3339),
		PackageName:   packageName,
		StdPackages:   assembleStdPackagesDTO(stdPackages),
	}
}

func assembleStdPackagesDTO(stdPackages []string) []outputStdPackage {
	result := make([]outputStdPackage, 0, len(stdPackages))

	for _, packageName := range stdPackages {
		result = append(result, outputStdPackage{
			Name: packageName,
		})
	}

	return result
}

func extractStdPackages() ([]string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName,
	}

	stdList, err := packages.Load(cfg, "std")
	if err != nil {
		return nil, fmt.Errorf("failed load std packages info: %w", err)
	}

	result := make([]string, 0, len(stdList))
	for _, stdPkg := range stdList {
		result = append(result, stdPkg.PkgPath)
	}

	return result, nil
}
