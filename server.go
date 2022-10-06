package main

import (
	"encoding/base64"
	"notes-service/auth"

	"context"
	"errors"
	"fmt"
	"net"
	"notes-service/models"
	notespb "notes-service/protorepo/noted/notes/v1"
	"strings"
	"time"

	mongoServices "notes-service/models/mongo"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	logger  *zap.Logger
	slogger *zap.SugaredLogger

	authService auth.Service

	mongoDB *mongoServices.Database

	notesRepository  models.NotesRepository
	blocksRepository models.BlocksRepository

	notesService notespb.NotesAPIServer

	grpcServer *grpc.Server
}

func (s *server) Init(opt ...grpc.ServerOption) {
	s.initLogger()
	s.initAuthService()
	s.initRepositories()
	s.initNotesService()
	s.initgrpcServer(opt...)
}

func (s *server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprint(":", *port))
	must(err, "failed to create tcp listener")
	reflection.Register(s.grpcServer)
	s.slogger.Infof("service running on : %d", *port)
	err = s.grpcServer.Serve(lis)
	must(err, "failed to run grpc server")
}

func (s *server) Close() {
	s.logger.Info("shutdown")
	s.mongoDB.Disconnect(context.Background())
	s.logger.Sync()
}

func (s *server) LoggerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	res, err := handler(ctx, req)
	end := time.Now()

	method := info.FullMethod[strings.LastIndexByte(info.FullMethod, '/')+1:]

	if err != nil {
		var displayErr = err
		st, ok := status.FromError(err)
		if ok {
			displayErr = errors.New(st.Message())
		}
		s.logger.Warn("failed rpc",
			zap.String("code", status.Code(err).String()),
			zap.String("method", method),
			zap.Duration("duration", end.Sub(start)),
			zap.Error(displayErr),
		)
		return res, err
	}

	s.logger.Info("rpc",
		zap.String("code", status.Code(err).String()),
		zap.String("method", method),
		zap.Duration("duration", end.Sub(start)),
	)

	return res, nil
}

func (s *server) initLogger() {
	var err error
	if *environment == envIsProd {
		s.logger, err = zap.NewProduction(zap.AddStacktrace(zapcore.FatalLevel), zap.WithCaller(false))
	} else {
		s.logger, err = zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel), zap.WithCaller(false))
	}
	must(err, "unable to instantiate zap.Logger")
	s.slogger = s.logger.Sugar()
}

func (s *server) initAuthService() {
	rawKey, err := base64.StdEncoding.DecodeString(*jwtPrivateKey)
	must(err, "could not decode jwt private key")
	s.authService = auth.NewService(rawKey)
}

func (s *server) initNotesService() {
	s.notesService = &notesService{
		auth:      s.authService,
		logger:    s.slogger,
		repoNote:  s.notesRepository,
		repoBlock: s.blocksRepository,
	}
}

func (s *server) initgrpcServer(opt ...grpc.ServerOption) {
	s.grpcServer = grpc.NewServer(opt...)
	notespb.RegisterNotesAPIServer(s.grpcServer, s.notesService)
}

func (s *server) initRepositories() {
	var err error

	/*
		client, err := mongo.NewClient(option.Client().ApplyURI(*mongoUri))
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			defer client.Disconnect(ctx)
			log.Fatal(err)
		}
		//s.mongoDB.DB = client.Database("Notes-database")
		if err != nil {
			log.Fatal(err)
		}*/
	s.mongoDB, err = mongoServices.NewDatabase(context.Background(), *mongoUri, *mongoDbName, s.logger)
	must(err, "could not instantiate mongo database")
	s.notesRepository = mongoServices.NewNotesRepository(s.mongoDB.DB, s.logger)
	s.blocksRepository = mongoServices.NewBlocksRepository(s.mongoDB.DB, s.logger, s.notesRepository)
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %v", msg, err))
	}
}
