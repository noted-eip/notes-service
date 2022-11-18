package main

import (
	"os"

	"notes-service/auth"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"

	"gopkg.in/alecthomas/kingpin.v2"
)

var password = "Gyy628\\nAWS"

var (
	app = kingpin.New("notes-service", "Notes service for the Noted backend").DefaultEnvars()

	environment = app.Flag("env", "either development or production").Default(envIsProd).Enum(envIsProd, envIsDev)
	port        = app.Flag("port", "grpc server port").Default("3001").Int()
	//mongoUri    = app.Flag("mongo-uri", "mongo uri with password to connect client").Default("mongodb://localhost:27017").String()
	mongoUri      = app.Flag("mongo-uri", "mongo uri with password to connect client").Default("mongodb+srv://gabriel:" + password + "@cluster0.2ckb3.mongodb.net/myFirstDatabase?retryWrites=true&w=majority").String()
	mongoDbName   = app.Flag("mongo-db-name", "name of the mongo database").Default("notes-service").String()
	jwtPrivateKey = app.Flag("jwt-private-key", "base64 encoded ed25519 private key").Default("SGfCQAb05CtmhEesWxcrfXSQR6JjmEMeyjR7Mo21S60ZDW9VVTUuCvEMlGjlqiw4I/z8T11KqAXexvGIPiuffA==").String()
	//recommandationUrl = app.Flag("recommandationUrl", "recommandation api Url").Default(":3000").String()
	recommandationUrl = app.Flag("recommandationUrl", "recommandation api Url").Default("https://localhost:3000").String()

	NotesDatabase   *mongo.Database   = nil
	NotesCollection *mongo.Collection = nil
)

var (
	envIsProd = "production"
	envIsDev  = "development"
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	s := &server{}
	s.Init(grpc.ChainUnaryInterceptor(s.LoggerUnaryInterceptor, auth.ForwardAuthMetadatathUnaryInterceptor))
	s.Run()
	defer s.Close()
}
