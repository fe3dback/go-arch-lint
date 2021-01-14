## Development

1. Use DI containers for all struct constructors (see cmd/container)
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

### Docker images

Docker will build image for some refs:

| git tag | docker image | description | example |
| ------- | ------------ | ----------- | ------- |
| master  | :latest | latest dev |
| vX.Y.Z  | :release-vX.Y.Z | release version | v1.3.0
| vX.Y.Z-rcX  | :dev-vX.Y.Z-rcX | dev version | v1.3.0-rc1, v1.3.0-rc16

Tags must expect the semantic version.

### Currently, version maintained by hand

Before release, update version in code places

1. go - internal/version/version.go (const `VERSION`)
2. git - tag commit with same version from const
4. fix functional tests (test/version)