# golang architecture linter

Check all project imports and compare to arch rules defined in yml file

![Logo image](https://user-images.githubusercontent.com/2073883/93022069-84124480-f5ef-11ea-93b6-614015a3d6d6.png)

## Install

```bash
go get -u github.com/fe3dback/go-arch-lint
```

go will download and install binary to bin folder, usually
is ~/go/bin

## Run

Run binary with flag "--project-path" and absolutely path
to your project, for example:
```
go-arch-lint --project-path ~/go/src/github.com/fe3dback/galaxy
```

flags:

| Flag              | Req?  | Default  | Example             |
| ----------------- | ----- | -------- | ------------------- |
| project-path      | +     | n/a      | --project-path ~/go/src/github.com/fe3dback/galaxy |
| max-warnings      | -     | 512      | --max-warnings=32 |
| output-color      | -     | true     | --output-color=false |
| output-type       | -     | ascii    | --output-color=json |

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

## Example of usage

This linter will return:

| Status Code | Description |
| ----------- | ----------- |
| 0           | Project corresponds for defined rules |
| 1           | Found warnings |

```
λ ~/ go-arch-lint check --project-path ~/go/src/github.com/fe3dback/galaxy
used arch file: /home/neo/go/src/github.com/fe3dback/galaxy/.go-arch-lint.yml
        module: github.com/fe3dback/galaxy
[WARNING] File '/home/neo/go/src/github.com/fe3dback/galaxy/engine/lib/sound/manager.go' not attached to any component in archfile
[WARNING] Component 'engine_loader': file '/engine/loader/assets_loader.go' shouldn't depend on 'github.com/fe3dback/galaxy/engine'
```

### json

Same warnings in json format

```
λ ~/ go-arch-lint check --project-path ~/go/src/github.com/fe3dback/galaxy --output-type=json
```

```json
{
  "Type": "cmd.checkPayload",
  "Payload": {
    "HasWarnings": true,
    "WarningsDeps": [
      {
        "ComponentName": "engine_loader",
        "FileRelativePath": "/engine/loader/assets_loader.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/engine/loader/assets_loader.go",
        "ResolvedImportName": "github.com/fe3dback/galaxy/engine"
      }
    ],
    "WarningsNotMatched": [
      {
        "FileRelativePath": "/engine/lib/sound/manager.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/engine/lib/sound/manager.go"
      }
    ]
  }
}
```