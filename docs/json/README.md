# Json output mode

## models.Check

*example (checker warnings):*
```json
{
  "Type": "models.Check",
  "Payload": {
    "ExecutionWarnings": [],
    "ArchHasWarnings": true,
    "ArchWarningsDeps": [
      {
        "ComponentName": "game_loader",
        "FileRelativePath": "/game/loader/weaponloader/parser.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/game/loader/weaponloader/parser.go",
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

*example (invalid archfile syntax):*
```json
{
  "Type": "models.Check",
  "Payload": {
    "ExecutionWarnings": [
      {
        "Text": "vendor path 'golang.org/x/mod/madfile' not valid, or no packages found by glob (project use gomod? try 'go mod vendor'), err: not found directories for 'vendor/golang.org/x/mod/madfile' in '/home/neo/go/src/github.com/fe3dback/go-arch-lint/vendor/golang.org/x/mod/madfile'",
        "File": "/home/neo/go/src/github.com/fe3dback/go-arch-lint/.go-arch-lint.yml",
        "Line": 11,
        "Offset": 9
      }
    ],
    "ArchHasWarnings": false,
    "ArchWarningsDeps": [],
    "ArchWarningsNotMatched": [],
    "ModuleName": "github.com/fe3dback/galaxy"
  }
}
```

## models.Version

*example:*
```json
{
  "Type": "models.Version",
  "Payload": {
    "LinterVersion": "2.0.0-rc7",
    "GoArchFileSupported": "1, 2",
    "BuildTime": "2020-12-30T15:47:12Z",
    "CommitHash": "f2f5624a070e0babf82598fb756b316423e47ba7"
  }
}
```

## models.Mapping

*example:*
```json
{
  "Type": "models.Mapping",
  "Payload": {
    "ProjectDirectory": "/home/neo/go/src/github.com/fe3dback/go-arch-lint",
    "ModuleName": "github.com/fe3dback/go-arch-lint",
    "MappingGrouped": [
      {
        "ComponentName": "app",
        "FileNames": [
          "/home/neo/go/src/github.com/fe3dback/go-arch-lint/internal/app/cli.go"
        ]
      },
      {
        "ComponentName": "commands",
        "FileNames": [
          "/home/neo/go/src/github.com/fe3dback/go-arch-lint/internal/commands/check/command.go",
          "/home/neo/go/src/github.com/fe3dback/go-arch-lint/internal/commands/check/flags.go",
          "/home/neo/go/src/github.com/fe3dback/go-arch-lint/internal/commands/mapping/command.go",
          "/home/neo/go/src/github.com/fe3dback/go-arch-lint/internal/commands/mapping/flags.go"
        ]
      }
    ],
    "MappingList": [
      {
        "FileName": "/home/neo/go/src/github.com/fe3dback/go-arch-lint/internal/app/cli.go",
        "ComponentName": "app"
      },
      {
        "FileName": "/home/neo/go/src/github.com/fe3dback/go-arch-lint/internal/app/internal/container/cmd_check.go",
        "ComponentName": "container"
      }
    ]
  }
}
```