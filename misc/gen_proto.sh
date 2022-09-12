rm -rf grpc/*pb/
protoc --go_out=. --go-grpc_out=. grpc/protos/notes/*.proto