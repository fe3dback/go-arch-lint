## config

### validation
- if "depOnAnyVendor=false" - print warning about all defined vendor imports in config (will not be checked)

### tools
- fmt?

## documentation

### glob

- used https://github.com/gobwas/glob
- examples: path "operations/mapping" matched by:
  - `operations`:
    - operations
  - `operations/*`
    - operations/mapping
  - `operations/**`
    - operations
    - operations/mapping