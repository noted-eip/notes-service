## Project rules
# default rule
all:	build
# Make if go.mod is already init
build:
	go build .
# Create an executable for service note
re:
	rm go.mod
	go mod init notes-service
	go mod tidy
# clean go.mod, init project, build it and create an executable 
re:	init all

## Submodules rules
# Run the protoc compiler to generate the Golang server code
codegen: update-submodules
	docker run --rm -v `pwd`/grpc:/app/grpc -v `pwd`/misc:/app/misc -w /app noted-go-protoc /bin/sh -c misc/gen_proto.sh
# Run the golangci-lint linter.
lint:
	docker run -w /app -v `pwd`:/app:ro golangci/golangci-lint golangci-lint run

# Fetch the latest version of the protos submodule.
update-submodules:
	git submodule update --recursive grpc/protos
	git submodule update --remote
#	git pull --recurse-submodules
# After cloning the repo, run init
init-submodules:
	git submodule init .gitmodules
	docker build -t noted-go-protoc -f misc/Dockerfile .

## DataBase rules
# no rules we have a oncline cluster
