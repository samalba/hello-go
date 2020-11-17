FROM golang:1.15 AS builder
WORKDIR /go/src/app
COPY . .
RUN \
	GO111MODULE=on GOOS=linux GOARCH=amd64 \
	go build -o hello-go

FROM golang:1.15
COPY --from=builder /go/src/app/hello-go /bin/hello-go
WORKDIR /app

ENV PORT 80

ENTRYPOINT ["/bin/hello-go"]
