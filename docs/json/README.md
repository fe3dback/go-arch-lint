# Json output mode

## command.check

*struct:*
```go
type payloadTypeCommandCheck struct {
	ExecutionWarnings      []spec.YamlAnnotatedWarning `json:"execution_warnings"`
	ExecutionError         string                      `json:"execution_error"`
	ArchHasWarnings        bool                        `json:"arch_has_warnings"`
	ArchWarningsDeps       []checker.WarningDep        `json:"arch_warnings_deps"`
	ArchWarningsNotMatched []checker.WarningNotMatched `json:"arch_warnings_not_matched"`
}
```

*example (checker warnings):*
```json
{
  "Type": "command.check",
  "Payload": {
    "execution_warnings": [],
    "execution_error": "",
    "arch_has_warnings": true,
    "arch_warnings_deps": [
      {
        "ComponentName": "game_loader",
        "FileRelativePath": "/game/loader/weaponloader/parser.go",
        "FileAbsolutePath": "/home/neo/go/src/github.com/fe3dback/galaxy/game/loader/weaponloader/parser.go",
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

*example (invalid archfile syntax):*
```json
{
  "Type": "command.check",
  "Payload": {
    "execution_warnings": [
      {
        "Text": "path '$.version': version 2 is not supported, supported: [1]",
        "Path": "$.version",
        "Line": 1,
        "Offset": 10
      }
    ],
    "execution_error": "failed to parse archfile: spec '/home/neo/go/src/github.com/fe3dback/galaxy/.go-arch-lint.yml' has warnings",
    "arch_has_warnings": false,
    "arch_warnings_deps": [],
    "arch_warnings_not_matched": []
  }
}
```

## command.version

*struct:*
```go
type payloadVersion struct {
	LinterVersion       string `json:"linter_version"`
	GoArchFileSupported string `json:"go_arch_file_supported"`
}
```

*example:*
```json
{
  "Type": "command.version",
  "Payload": {
    "linter_version": "1.1.0",
    "go_arch_file_supported": "1"
  }
}
```

## runtime cli panic/error

*struct:*
```go
type payloadTypeHalt struct {
	Error string `json:"error"`
}
```

*example:*
```json
{
  "Type": "halt",
  "Payload": {
    "error": "panic: flag 'max-warnings' should by in range [1 .. 32768]"
  }
}
```
