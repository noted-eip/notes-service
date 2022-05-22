package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"notes-service/grpc/notespb"

	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	log "github.com/sirupsen/logrus"
)

var NotesDatabase *mongo.Database = nil
var NotesCollection *mongo.Collection = nil

var password = "Gyy628\\nAWS"
var mongoUri string = "mongodb+srv://gabriel:" + password + "@cluster0.2ckb3.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

func InitDatabase() {
	client, err := mongo.NewClient(option.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	fmt.Println("service is connected to MongoDB")
	//defer client.Disconect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	NotesDatabase = client.Database("Notes-database")
	NotesCollection = NotesDatabase.Collection("Notes-collection")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	InitDatabase()

	srv := grpc.NewServer()
	notesSrv := notesService{}
	notespb.RegisterNotesServiceServer(srv, &notesSrv)

	lis, err := net.Listen("tcp", ":3000")

	fmt.Println("service listen on port 3000")

	if err != nil {
		panic(err)
	}
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
