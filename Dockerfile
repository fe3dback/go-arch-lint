FROM golang:1.18-buster

# currently we use `packages.Load(..)` feature of go AST parser
# this require go binary and some random files from go setup
# to work. And this is reason why current build not from scratch
#
# todo: remove all not needed files from buster image
# todo: replace packages.Load to something else?

COPY go-arch-lint /
ENTRYPOINT ["/go-arch-lint"]
