# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-alpine

ENV GO111MODULE=on

WORKDIR /example-grpc-auth


# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY ./docker/config.json ./config/

RUN go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=$(hostname)
RUN go build -v -o ./exampleapp ./cmd/main.go


CMD [ "./exampleapp" ]
