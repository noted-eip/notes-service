EXE	=	notes-service

all:	build

build:
	go build .

re: clean all

clean:
	rm ${EXE}

update-submodules:
	git submodule update --init --remote

lint:
	docker run -w /app -v `pwd`:/app:ro golangci/golangci-lint golangci-lint run

run-db:
	docker run --name notes-mongo --detach --publish 27017:27017 mongo

stop-db:
	docker kill notes-mongo
	docker rm notes-mongo
