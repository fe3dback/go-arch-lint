![Logo image](https://user-images.githubusercontent.com/2073883/94179282-f82cd200-fea4-11ea-85c5-bf685293220e.png)

Check all project imports and compare to arch rules defined in yml file

[![Go Report Card](https://goreportcard.com/badge/github.com/fe3dback/go-arch-lint)](https://goreportcard.com/report/github.com/fe3dback/go-arch-lint)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/fe3dback/go-arch-lint)

## Quick start

### Install globally and run

```bash
go get -u github.com/fe3dback/go-arch-lint
```

go will download and install `go-arch-lint` binary to bin folder, usually
is `~/go/bin`

```bash
go-arch-lint check --project-path /home/user/go/src/github.com/fe3dback/galaxy

// or inside project directory:
cd project_dir
go-arch-lint check
```

### alternative - run with docker

```bash
docker run --rm \
    -v ${PWD}:/app \
    fe3dback/go-arch-lint:latest-stable-release check --project-path /app
```

[all docker tags](https://hub.docker.com/r/fe3dback/go-arch-lint/tags)

### precompiled binaries and other options

[releases page](https://github.com/fe3dback/go-arch-lint/releases)

### IDE plugin for autocompletion and other help

<img src="https://user-images.githubusercontent.com/2073883/104641610-0f453900-56bb-11eb-8419-6d94fbcb4d2f.png" alt="jetbrains goland logo" align="left" width="80px" height="80px">

https://plugins.jetbrains.com/plugin/15423-goarchlint-file-support

currently IDE plugin in alpha

## Usage

```
Usage:
  go-arch-lint check [flags]

Flags:
      --arch-file string      arch file path (default ".go-arch-lint.yml")
  -h, --help                  help for check
      --max-warnings int      max number of warnings to output (default 512)
      --project-path string   absolute path to project directory (where '.go-arch-lint.yml' is located) (default "./")

Global Flags:
      --json                   (alias for --output-type=json)
      --output-color           use ANSI colors in terminal output (default true)
      --output-json-one-line   format JSON as single line payload (without line breaks), only for json output type
      --output-type string     type of command output, variants: [ascii, json] (default "default")
```

## Archfile example

Make archfile called `.go-arch-lint.yml` in root directory
of your project, and put some arch rules to it

```yaml
version: 2
workdir: ./
allow:
  depOnAnyVendor: false

# ----------------------------------
# Exclude from analyse
# ----------------------------------
exclude:
  - .idea
  - vendor

excludeFiles:
  - "^.*_test\\.go$"
  - "^.*test/mock/.*\\.go$"

# ----------------------------------
# Vendor libs
# ----------------------------------
vendors:
  vectors:
    in: github.com/fe3dback/go-vec
  company-libs:
    in: example.com/*/libs/**
  loaders:
    in:
      - gopkg.in/yaml.v2
      - github.com/mailru/easyjson

# ----------------------------------
# Project components
# ----------------------------------
components:
  main:
    in: .
  engine:
    in: engine
  engineVendorEvents:
    in: 
      - engine/vendor/*/event
      - engine/company/*/event
  game:
    in: game
  gameComponent:
    in: game/components/**
  utils:
    in: utils

# ----------------------------------
# All components can import any 
# other components from "common" list
# ----------------------------------
commonComponents:
  - utils

# ----------------------------------
# All components can import any 
# vendors from "common" list
# ----------------------------------
commonVendors:
  - vectors
  - company-libs

# ----------------------------------
# Dependency rules
# ----------------------------------
deps:
  engine:
    canUse: # = can import vendor lib
      - loaders 

  engineVendorEvents:
    mayDependOn: # = can import another project package
      - engine

  game:
    mayDependOn:
      - engine
      - gameComponent

  main:
    mayDependOn:
      - game
      - engine
```

This project also uses arch lint, see example in [.go-arch-lint.yml](.go-arch-lint.yml)

## Archfile Syntax

| Path              | Req?  | Type  | Description         |
| -------------     | ----- | ----- | ------------------- |
| version           | `+`   | int   | schema version (__latest: 2__)  |
| workdir           | -     | str   | relative directory for analyse  |
| allow             | -     | map   | global rules |
| . depOnAnyVendor  | -     | bool  | allow import any vendor code to any project file |
| exclude           | -     | []str  | list of directories (relative path) for exclude from analyse |
| excludeFiles      | -     | []str  | regular expression rules for file names, will exclude this files and it's packages from analyse |
| components        | `+`   | map   | project components used for split real modules and packages to abstract thing |
| . %name%          | `+`   | str   | name of component |
| . . in            | `+`   | str, []str   | one or more relative directory name, support glob masking (src/\*/engine/\*\*) |
| vendors           | -     | map   | vendor libs |
| . %name%          | `+`   | str   | name of vendor component |
| . . in            | `+`   | str, []str   | one or more import path of vendor libs, support glob masking (src/\*/engine/\*\*) |
| commonComponents  | -     | []str  | list of components, allow import them into any code |
| commonVendors     | -     | []str  | list of vendors, allow import them into any code |
| deps              | `+`   | map   | dependency rules |
| . %name%          | `+`   | str   | name of component, exactly as defined in "components" section |
| . . anyVendorDeps | -     | bool  | all component code can import any vendor code |
| . . anyProjectDeps| -     | bool  | all component code can import any other project code, useful for DI/main component |
| . . mayDependOn   | -     | []str  | list of components that can by imported in %name% |
| . . canUse        | -     | []str  | list of vendors that can by imported in %name% |

## Example of usage

This linter will return:

| Status Code | Description |
| ----------- | ----------- |
| 0           | Project corresponds for defined rules |
| 1           | Found warnings |

```text
$ go-arch-lint check --project-path ~/go/src/github.com/fe3dback/galaxy
Module: github.com/fe3dback/galaxy
[WARN] Component 'game_entities_factory': file '/game/entities/factory/bullet.go' shouldn't depend on 'github.com/fe3dback/galaxy/game/entities/components/game'
[WARN] Component 'game_loader': file '/game/loader/weaponloader/loader.go' shouldn't depend on 'github.com/fe3dback/galaxy/engine'
[WARN] File '/shared/ui/layer_shared_fps.go' not attached to any component in archfile

warnings found: 3
```

## Advanced examples

### json

Same warnings in json format

```
$ go-arch-lint check --project-path ~/go/src/github.com/fe3dback/galaxy --json
```

```json
{
  "Type": "models.Check",
  "Payload": {
    "ExecutionWarnings": [],
    "ArchHasWarnings": true,
    "ArchWarningsDeps": [
      {
        "ComponentName": "game_entities_factory",
        "FileRelativePath": "/game/entities/factory/bullet.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/game/entities/factory/bullet.go",
        "ResolvedImportName": "github.com/fe3dback/galaxy/game/entities/components/game"
      },
      {
        "ComponentName": "game_loader",
        "FileRelativePath": "/game/loader/weaponloader/loader.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/game/loader/weaponloader/loader.go",
        "ResolvedImportName": "github.com/fe3dback/galaxy/engine"
      }
    ],
    "ArchWarningsNotMatched": [
      {
        "FileRelativePath": "/shared/ui/layer_shared_fps.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/shared/ui/layer_shared_fps.go"
      }
    ],
    "ModuleName": "github.com/fe3dback/galaxy"
  }
}
```

More json examples in [docs](docs):

### json schema for archfile

linter can export self schema wia `schema --version X` command

```bash
go-arch-lint schema --version 2
{"$schema":"http://json-schema.org/draft-07/schema#","additionalProperties":false,"definitions":{"commonComponents":{"description":"All project packages ... }
```

this will be useful for auto-complete and validation in another editors

### mapping

you can see archfile mapping to source files wia `mapping` command

two modes available:
- list (default)
- grouped by component

```bash
go-arch-lint mapping

module: github.com/fe3dback/go-arch-lint
Project Packages:
   app                 /internal/app
   container           /internal/app/internal/container
   commands            /internal/commands/check
   commands            /internal/commands/mapping
   ...
```

```bash
go-arch-lint mapping --scheme grouped

module: github.com/fe3dback/go-arch-lint
Project Packages:
   app:
     /internal/app
   commands:
     /internal/commands/check
     /internal/commands/mapping
   ...
```

same data available in json format, with `--json` option
