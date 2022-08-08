FROM golang:1.18.1-alpine as build

WORKDIR /app

#useless ?
ENV PATH="/home/.local/bin:${PATH}"
ENV GO111MODULE=on

COPY . .

RUN apk update && apk add --no-cache make protobuf-dev=3.18.1-r1
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

#WORKDIR /app/grpc
#COPY grpc/* ./
#WORKDIR /app

# le faire préalablement ?
#RUN go mod init notes-service
#RUN go mod tidy

RUN ./misc/gen_proto.sh

#COPY . .
#RUN go mod download
#RUN go mod vendor
#RUN go mod verify

#COPY *.go ./
#RUN ./misc/gen_proto.sh

#RUN go build -buildvcs=false -o /app .
RUN go build -buildvcs=false .

FROM alpine:latest

COPY --from=build /app/notes-service .

EXPOSE 3000

CMD [ "./notes-service" ]