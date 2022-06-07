FROM golang:1.18.1-alpine as build

WORKDIR /app

ENV PATH="/home/.local/bin:${PATH}"
ENV GO111MODULE=on

COPY . .

RUN apk update && apk add --no-cache make protobuf-dev=3.18.1-r1
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN ./misc/gen_proto.sh

RUN go build -buildvcs=false .

FROM alpine:latest

COPY --from=build /app/notes-service .

EXPOSE 3000

CMD [ "./notes-service" ]