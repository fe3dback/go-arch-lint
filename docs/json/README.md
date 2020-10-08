# Json output mode

## command.check

*struct:*
```go
type payloadTypeCommandCheck struct {
	ExecutionWarnings      []annotated_validator.YamlAnnotatedWarning
	ExecutionError         string
	ArchHasWarnings        bool
	ArchWarningsDeps       []checker.WarningDep
	ArchWarningsNotMatched []checker.WarningNotMatched
}
```

*example (checker warnings):*
```json
{
  "Type": "command.check",
  "Payload": {
    "ExecutionWarnings": [],
    "ExecutionError": "",
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
    ]
  }
}
```

*example (invalid archfile syntax):*
```json
{
  "Type": "command.check",
  "Payload": {
    "ExecutionWarnings": [
      {
        "Text": "path '$.version': version 2 is not supported, supported: [1]",
        "Path": "$.version",
        "Line": 1,
        "Offset": 10
      }
    ],
    "ExecutionError": "failed to parse archfile: spec '/home/neo/go/src/github.com/fe3dback/galaxy/.go-arch-lint.yml' has warnings",
    "ArchHasWarnings": false,
    "ArchWarningsDeps": [],
    "ArchWarningsNotMatched": []
  }
}
```

## command.version

*struct:*
```go
type payloadVersion struct {
	LinterVersion       string
	GoArchFileSupported string
}
```

*example:*
```json
{
  "Type": "command.version",
  "Payload": {
    "LinterVersion": "1.1.0",
    "GoArchFileSupported": "1"
  }
}
```

## runtime cli panic/error

*struct:*
```go
type payloadTypeHalt struct {
	Error string
}
```

*example:*
```json
{
  "Type": "halt",
  "Payload": {
    "Error": "panic: cmd: flag 'max-warnings' should by in range [1 .. 32768]"
  }
}
```
