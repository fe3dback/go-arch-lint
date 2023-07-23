## Project structure

All code related to linter/binary specified in `internal` directory

Internal directory contain some subdirectories:
- app (run logic and DI)
- models (shared DTO for usage in all other layers)
- operations (entrypoint for every binary command (check, graph, version, etc..))
- services (additional code that can be reused between `operations`)
- - checker - code related to `check` cmd and go code linter main process itself.
- - project - code related to working with project files (scan, group, mapping to components, etc..)
- - render - code related to terminal output and ascii rendering
- - schema - json schemas for every supported version of the DSL config
- - spec - code related to work with yaml DSL config (parse/validate/transform)
- view - stdout templates for every command

## Development

1. Use DI containers for all struct constructors (see app/internal/container)
2. All other struct/processors should be injected only wia constructor as interface
3. All injected interfaces should be described in file called "types.go" 

## Quality control

Please use external linters and quality control tools, such as IDE

- Before push, check unit and functional tests

    ```bash
    make tests
    ```

- if command output changed, functional tests will fall. You can
update test files wia command

    ```bash
    make tests-functional-update-ct
    ```

## Releases

Use semver for git tags 

### Docker images

CI will build docker image for some refs:

| git tag    | docker image    | description     | example                 |
|------------|-----------------|-----------------|-------------------------|
| master     | :latest         | latest dev      |                         |
| vX.Y.Z     | :release-vX.Y.Z | release version | v1.3.0                  |
| vX.Y.Z-rcX | :dev-vX.Y.Z-rcX | dev version     | v1.3.0-rc1, v1.3.0-rc16 |
