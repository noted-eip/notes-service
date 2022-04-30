
# Create an executable for service note
re:
	rm go.mod
	go mod init service
	go build

# Run the protoc compiler to generate the Golang server code.
codegen: update-submodules
	protoc --go_out=. --go-grpc_out=. grpc/protos/notes/*.proto

# Fetch the latest version of the protos submodule.
update-submodules:
	git submodule init
	git submodule update --remote
