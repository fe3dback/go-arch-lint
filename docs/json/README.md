# Json output mode

## models.Check

*struct:*
```go
type Check struct {
    DocumentNotices        []CheckNotice                `json:"ExecutionWarnings"`
    ArchHasWarnings        bool                         `json:"ArchHasWarnings"`
    ArchWarningsDependency []CheckArchWarningDependency `json:"ArchWarningsDeps"`
    ArchWarningsMatch      []CheckArchWarningMatch      `json:"ArchWarningsNotMatched"`
    ModuleName             string                       `json:"ModuleName"`
}
```

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
        "Text": "version '333' is not supported, supported: [1]",
        "File": "/home/neo/go/src/github.com/fe3dback/galaxy/.go-arch-lint.yml",
        "Line": 1,
        "Offset": 10
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

*struct:*
```go
type Version struct {
    LinterVersion       string `json:"LinterVersion"`
    GoArchFileSupported string `json:"GoArchFileSupported"`
    BuildTime           string `json:"BuildTime"`
    CommitHash          string `json:"CommitHash"`
}
```

*example:*
```json
{
  "Type": "models.Version",
  "Payload": {
    "LinterVersion": "1.4.3",
    "GoArchFileSupported": "1",
    "BuildTime": "unknown",
    "CommitHash": "unknown"
  }
}
```