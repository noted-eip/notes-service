## Project rules
# default rule
all:	build
# Make if go.mod is already init
build:
	go build .
#init project
init:
	rm go.mod
	go mod init notes-service
	go mod tidy
# clean go.mod, init project, build it and create an executable 
re:	init all


## Submodules rules
# Run the protoc compiler to generate the Golang server code
codegen: update-submodules
	docker run --rm -v `pwd`/protorepo:/app/protorepo -v `pwd`/misc:/app/misc -w /app noted-go-protoc /bin/sh -c misc/gen_proto.sh
# Fetch the latest version of the protos submodule.                             
update-submodules:
#	git submodule update --recursive protorepo
	git submodule update --remote
#	git pull --recurse-submodules
# After cloning the repo, run init
init-submodules:
	git submodule init .gitmodules
	docker build -t noted-go-protoc -f misc/Dockerfile .
# Run the golangci-lint linter.
lint:
	docker run -w /app -v `pwd`:/app:ro golangci/golangci-lint golangci-lint run


## DataBase rules
# Run MongoDB database as a docker container.
run-db:
	docker run --name accounts-mongo --detach --publish 27017:27017 mongo
# Stop MongoDB database.
stop-db:
	docker kill accounts-mongo