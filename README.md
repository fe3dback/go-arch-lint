# golang architecture linter

[![Go Report Card](https://goreportcard.com/badge/github.com/fe3dback/go-arch-lint)](https://goreportcard.com/report/github.com/fe3dback/go-arch-lint)

Check all project imports and compare to arch rules defined in yml file

![Logo image](https://user-images.githubusercontent.com/2073883/94179282-f82cd200-fea4-11ea-85c5-bf685293220e.png)

## Install

```bash
go get -u github.com/fe3dback/go-arch-lint
```

go will download and install binary to bin folder, usually
is ~/go/bin

## Run

Run binary with flag "check --project-path" and absolutely path
to your project, for example:
```
go-arch-lint check --project-path ~/go/src/github.com/fe3dback/galaxy
```

## Usage

```
Usage:
  go-arch-lint [flags]
  go-arch-lint [command]

Available Commands:
  check       check project architecture by yaml file
  help        Help about any command
  version     Print go arch linter version

Flags:
  -h, --help                   help for go-arch-lint
      --json                   (alias for --output-type=json)
      --max-warnings int       max number of warnings to output (default 512)
      --output-color           use ANSI colors in terminal output (default true)
      --output-json-one-line   format JSON as single line payload (without line breaks), only for json output type
      --output-type string     type of command output, variants: [ascii, json] (default "default")
      --project-path string    absolute path to project directory (where '.go-arch-lint.yml' is located)
```

## Archfile example

Make archfile called '.go-arch-lint.yml' in root directory
of your project, and put some arch rules to it

See full examples in /examples

```yaml
version: 1
allow:
  depOnAnyVendor: false

# ----------------------------------
# Excluded folders from analyse
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
  loader-yaml:
    in: gopkg.in/yaml.v2
  vectors:
    in: github.com/fe3dback/go-vec

# ----------------------------------
# Project components
#
# Used for split real modules and 
# packages to abstract thing
# ----------------------------------
components:
  main:
    in: .
  engine:
    in: engine
  engine_vendor_events:
    in: engine/vendor/*/event
  game:
    in: game
  game_component:
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

# ----------------------------------
# Dependency rules
# ----------------------------------
deps:
  engine:
    canUse:
      - loader-yaml  

  engine_vendor_events:
    mayDependOn:
      - engine

  game:
    mayDependOn:
      - engine
      - game_component

  main:
    mayDependOn:
      - game
      - engine
```

## Archfile Syntax

| Path              | Req?  | Type  | Description         |
| -------------     | ----- | ----- | ------------------- |
| version           | +     | int   | schema version, currently support "1"  |
| allow             | -     | map   | global rules |
| . depOnAnyVendor  | -     | bool  | allow import any vendor code to any project file |
| exclude           | -     | list  | list of directories (relative path) for exclude from analyse |
| excludeFiles      | -     | list  | regExp rules for file names, for exclude from analyse |
| components        | +     | map   | project components used for split real modules and packages to abstract thing |
| . %name%          | +     | str   | name of component |
| . . in            | +     | str   | relative directory name, support glob masking (src/\*/engine/\*\*) |
| vendors           | -     | map   | vendor libs |
| . %name%          | +     | str   | name of vendor component |
| . . in            | +     | str   | full import path |
| commonComponents  | -     | list  | list of components, allow import them into any code |
| commonVendors     | -     | list  | list of vendors, allow import them into any code |
| deps              | +     | map   | dependency rules |
| . %name%          | +     | str   | name of component, exactly as defined in "components" section |
| . . anyVendorDeps | -     | bool  | all component code can import any vendor code |
| . . anyProjectDeps| -     | bool  | all component code can import any other project code, useful for DI/main component |
| . . mayDependOn   | -     | list  | list of components that can by imported in %name% |
| . . canUse        | -     | list  | list of vendors that can by imported in %name% |


## Example of usage

This linter will return:

| Status Code | Description |
| ----------- | ----------- |
| 0           | Project corresponds for defined rules |
| 1           | Found warnings |

```text
$ go-arch-lint check --project-path ~/go/src/github.com/fe3dback/galaxy
used arch file: /home/neo/go/src/github.com/fe3dback/galaxy/.go-arch-lint.yml
        module: github.com/fe3dback/galaxy
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
  "Type": "command.check",
  "Payload": {
    "execution_warnings": [],
    "execution_error": "",
    "arch_has_warnings": true,
    "arch_warnings_deps": [
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
    "arch_warnings_not_matched": [
      {
        "FileRelativePath": "/shared/ui/layer_shared_fps.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/shared/ui/layer_shared_fps.go"
      }
    ]
  }
}
```

Read more in [docs](docs):