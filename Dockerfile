FROM golang:1.14.2 AS builder
WORKDIR /go/src/app
COPY . .
RUN \
	GO111MODULE=on GOOS=linux GOARCH=amd64 \
	go build -tags netgo -ldflags '-w -extldflags "-static"' -o hello-go

FROM alpine:3.12@sha256:769fddc7cc2f0a1c35abb2f91432e8beecf83916c421420e6a6da9f8975464b6
COPY --from=builder /go/src/app/hello-go /bin/hello-go
WORKDIR /app

ENTRYPOINT ["/bin/hello-go"]
