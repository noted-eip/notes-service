rm -rf protorepo/noted/notes/v1/*pb/
protoc --go_out=. --go-grpc_out=. protorepo/noted/notes/v1/*.proto