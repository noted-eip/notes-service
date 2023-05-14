FROM golang:1.20.4-alpine as build
WORKDIR /app
COPY . .
# Required for VCS stamping.
RUN apk add --no-cache git
RUN go build .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/notes-service .
ENTRYPOINT [ "./notes-service" ]
