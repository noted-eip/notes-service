# This image is meant for building the notes-service and generating
# protobuf code in a consistent way.
FROM golang:1.18.1-alpine

RUN apk update && apk add --no-cache make protobuf-dev=3.18.1-r1
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
