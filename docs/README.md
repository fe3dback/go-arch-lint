# Docs

- [syntax](syntax/README.md)

## Advanced usages

### json schema for archfile

linter can export self schema wia `schema --version X` command

```bash
go-arch-lint schema --version 3
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
