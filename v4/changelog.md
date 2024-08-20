### SDK

All logic code separated from CLI code into own repository `github.com/fe3dback/go-arch-lint-sdk`

- linter can be used without CLI program
- linter can be used directly from GO unit tests
- The SDK still has the option to load the config from a file as usual. But now it is also possible to manually build it in GO code. 

### V4

Linter rewritten from scratch for better extensibility and support for several new features.

#### Config changes

- new config version: **v4**
  - major changes:
    - **deepscan** is **ON** by default
    - added new **strictMode** in `settings.imports.strictMode`
    - added new struct tags checker in `settings.structTags.allowed`
      - you can set it to 4 different modes:
        - `not set` = by default, if not set - will allow **ALL** tags
        - `true` = all struct tags is allowed (linter will not check tags)
        - `false` = no tags is allowed. Components can allow specific tags
        - `[db, bd, etc]` = list of allowed tags
      - added component override in `dependencies.{name}.canContainTags`
        - allow use specified tags in this component
  - minor changes:
    - `workdir` renamed to `workingDirectory`
    - `allow.depOnAnyVendor` moved into `settings.imports.allowAnyVendorImports`
    - `exclude` moved into `exclude.directories`
    - `excludeFiles` moved into `exclude.files`
    - `deps` renamed to `dependencies`
- added miss-use config validation
  - added flag `--skip-missuse` for skip this validation

#### Command changes

##### Mapping

> [!WARNING]
> Breaking changes

- when using command with json flag (ex `go-arch-lint mapping --json`), output DTO now contains list of go packages paths instead of individual go files paths.