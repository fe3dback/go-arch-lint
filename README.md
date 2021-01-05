![Logo image](https://user-images.githubusercontent.com/2073883/94179282-f82cd200-fea4-11ea-85c5-bf685293220e.png)

Check all project imports and compare to arch rules defined in yml file

[![Go Report Card](https://goreportcard.com/badge/github.com/fe3dback/go-arch-lint)](https://goreportcard.com/report/github.com/fe3dback/go-arch-lint)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/fe3dback/go-arch-lint)

## Quick start

### Install globally and run

```bash
go get -u github.com/fe3dback/go-arch-lint
```

go will download and install binary to bin folder, usually
is ~/go/bin

Run binary with flag "check --project-path" and absolutely path
to your project, for example:

```bash
go-arch-lint check --project-path /home/user/go/src/github.com/fe3dback/galaxy
```

### alternative - Run with docker

```bash
docker run --rm \
    -v /home/user/go/src/github.com/fe3dback/galaxy:/app \
    fe3dback/go-arch-lint:latest-stable-release check --project-path /app
```

[docker hub](https://hub.docker.com/r/fe3dback/go-arch-lint/tags)

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

Make archfile called '.go-arch-lint.yml' in root directory
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
#
# Used for split real modules and 
# packages to abstract namespaces
# ----------------------------------
components:
  main:
    in: .
  engine:
    in: engine
  engineVendorEvents:
    in: engine/vendor/*/event
  game:
    in: game
  gameComponent:
    in: game/components/**
  utils:
    in: utils

# ----------------------------------
# All components can import any 
# packages from "common" list
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
    canUse:
      - loaders 

  engineVendorEvents:
    mayDependOn:
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
| version           | +     | int   | schema version (__latest: 2__)  |
| workdir           | -     | str   | relative directory for analyse  |
| allow             | -     | map   | global rules |
| . depOnAnyVendor  | -     | bool  | allow import any vendor code to any project file |
| exclude           | -     | []str  | list of directories (relative path) for exclude from analyse |
| excludeFiles      | -     | []str  | regular expression rules for file names, will exclude this files and it's packages from analyse |
| components        | +     | map   | project components used for split real modules and packages to abstract thing |
| . %name%          | +     | str   | name of component |
| . . in            | +     | str   | relative directory name, support glob masking (src/\*/engine/\*\*) |
| vendors           | -     | map   | vendor libs |
| . %name%          | +     | str   | name of vendor component |
| . . in            | +     | str, []str   | one or more import path of vendor libs, support glob masking (src/\*/engine/\*\*) |
| commonComponents  | -     | []str  | list of components, allow import them into any code |
| commonVendors     | -     | []str  | list of vendors, allow import them into any code |
| deps              | +     | map   | dependency rules |
| . %name%          | +     | str   | name of component, exactly as defined in "components" section |
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

Read more in [docs](docs):
