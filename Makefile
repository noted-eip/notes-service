########### Variables ##########

EXE	=	notes-service

########### !Variables #########

########### Project rules ##########

all:	build

build:
	go build .

init:
	rm go.mod
	go mod init ${EXE}
	go mod tidy

re:	init all

clean:
	rm ${EXE}

########### !Project rules #########

########### Submodules rules ##########

codegen: update-submodules
	docker run --rm -v `pwd`/protorepo:/app/protorepo -v `pwd`/misc:/app/misc -w /app noted-go-protoc /bin/sh -c misc/gen_proto.sh

update-submodules:
	git submodule update --remote

init-submodules:
	git submodule init .gitmodules
	docker build -t noted-go-protoc -f misc/Dockerfile .

lint:
	docker run -w /app -v `pwd`:/app:ro golangci/golangci-lint golangci-lint run

########### DataBase rules ##########

run-db:
	docker run --name accounts-mongo --detach --publish 27017:27017 mongo

stop-db:
	docker kill accounts-mongo

########### !DataBase rules #########
