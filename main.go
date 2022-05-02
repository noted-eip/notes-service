package main

import (
    "notes-service/grpc/notespb"
    "net"
    
    "fmt"
)

func main() {
    srv := grpc.NewServer()
    notesSrv := notesService{}
    notespb.RegisterAccountsServiceServer(srv, &notesSrv)

    lis, err := net.Listen("tcp", ":3000")
    
    ftm.Println("server listen on port 3000")

	if err != nil {
		panic(err)
	}
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}