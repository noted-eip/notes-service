FROM golang:1.18.1-alpine

WORKDIR /notes-service

ENV PATH="/home/.local/bin:${PATH}"
ENV GO111MODULE=on

RUN apk update && apk add --no-cache make protobuf-dev=3.18.1-r1
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


WORKDIR /notes-service/grpc/
COPY grpc/* ./

WORKDIR /notes-service
RUN go mod init notes-service
RUN go mod tidy

COPY . .
RUN go mod download
#RUN go mod vendor
#RUN go mod verify

COPY *.go ./

#RUN ./misc/gen_proto.sh

RUN go build -buildvcs=false -o /notes-service .
EXPOSE 3000
CMD [ "./notes-service" ]