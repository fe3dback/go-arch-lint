FROM scratch
COPY go-arch-lint /
ENTRYPOINT ["/go-arch-lint"]