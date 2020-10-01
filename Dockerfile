FROM golang:1.15-alpine as builder

ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOOS=linux
ENV GOARCH=386

COPY . /tmp/builder
WORKDIR /tmp/builder

RUN go mod vendor
RUN go mod tidy
RUN go build -o /bin/go-arch-lint /tmp/builder

###

FROM scratch

COPY --from=builder /bin/go-arch-lint /bin/go-arch-lint

ENTRYPOINT ["/bin/go-arch-lint"]