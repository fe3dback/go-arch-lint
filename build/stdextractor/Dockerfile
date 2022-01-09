FROM golang:1-alpine

COPY ./* /src/

WORKDIR /src

RUN go build /src && \
    mkdir /.cache && \
    chown 1000:1000 -R /.cache

USER 1000:1000

ENTRYPOINT ["/src/stdextractor"]
