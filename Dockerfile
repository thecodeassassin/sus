### BUILDER ###

FROM golang:alpine AS builder

# Install build dependencies
RUN apk update && \
    apk add \
        git \
        ca-certificates \
        build-base \
        librdkafka-dev
WORKDIR /src
COPY . .
RUN  go get -d -v .
RUN GOOS=linux GOARCH=amd64 go build -v -o bin/appname ./cmd/appname

### FINAL IMAGE ###

# Install runtime dependencies
FROM alpine:3.8
RUN apk update && \
    apk add \
        ca-certificates

WORKDIR /src

COPY --from=0 /src/bin/appname /usr/local/bin/appname
ENTRYPOINT ["appname"]
