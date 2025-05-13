# Archfile Syntax

| Path                       | Req? | Type       | Description                                                                                     |
|----------------------------|------|------------|-------------------------------------------------------------------------------------------------|
| version                    | `+`  | int        | schema version (__latest: 3__)                                                                  |
| workdir                    |      | str        | relative directory for analyse                                                                  |
| allow                      |      | map        | global rules                                                                                    |
| . depOnAnyVendor           |      | bool       | allow import any vendor code to any project file                                                |
| . deepScan                 |      | bool       | use advanced AST code analyse (default `true`, since v3+).                                      |
| . ignoreNotFoundComponents |      | bool       | ignore not found components (default `false`)                                                   |
| exclude                    |      | []str      | list of directories (relative path) for exclude from analyse                                    |
| excludeFiles               |      | []str      | regular expression rules for file names, will exclude this files and it's packages from analyse |
| components                 | `+`  | map        | component is abstraction on go packages. One component = one or more go packages                |
| . %name%                   | `+`  | str        | name of component                                                                               |
| . . in                     | `+`  | str, []str | one or more relative directory name, support glob masking (src/\*/engine/\*\*)                  |
| vendors                    |      | map        | vendor libs (go.mod)                                                                            |
| . %name%                   | `+`  | str        | name of vendor component                                                                        |
| . . in                     | `+`  | str, []str | one or more import path of vendor libs, support glob masking (github.com/abc/\*/engine/\*\*)    |
| commonComponents           |      | []str      | list of components, allow import them into any code                                             |
| commonVendors              |      | []str      | list of vendors, allow import them into any code                                                |
| deps                       | `+`  | map        | dependency rules                                                                                |
| . %name%                   | `+`  | str        | name of component, exactly as defined in "components" section                                   |
| . . anyVendorDeps          |      | bool       | all component code can import any vendor code                                                   |
| . . anyProjectDeps         |      | bool       | all component code can import any other project code, useful for DI/main component              |
| . . mayDependOn            |      | []str      | list of components that can by imported in %name%                                               |
| . . canUse                 |      | []str      | list of vendors that can by imported in %name%                                                  |
| . . deepScan               |      | bool       | override of allow.deepScan for this component. Default `nil` = use global settings              |

Examples:
- [.go-arch-lint.yml](../../.go-arch-lint.yml)
