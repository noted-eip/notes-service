package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"notes-service/auth"
	"notes-service/language"

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

	authService     auth.Service
	languageService language.Service // NOTE: Could put directly service typed as NaturalAPIService, remove Init() from interface and just put it in NaturalAPIService

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
	s.initLanguageService()
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

func (s *server) initLanguageService() {
	s.languageService = &language.NaturalAPIService{}
	err := s.languageService.Init()

	must(err, "unable to instantiate language service")
}

func (s *server) initAuthService() {
	rawKey, err := base64.StdEncoding.DecodeString(*jwtPrivateKey)
	must(err, "could not decode jwt private key")
	var key ed25519.PrivateKey = rawKey
	pubKey := key.Public().(ed25519.PublicKey)
	s.authService = auth.NewService(pubKey)
}

func (s *server) initNotesService() {
	s.notesService = &notesService{
		language:  s.languageService,
		auth:      s.authService,
		logger:    s.logger,
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
	s.mongoDB, err = mongoServices.NewDatabase(context.Background(), *mongoUri, *mongoDbName, s.logger)
	must(err, "could not instantiate mongo database")
	s.notesRepository = mongoServices.NewNotesRepository(s.mongoDB.DB, s.logger)
	s.blocksRepository = mongoServices.NewBlocksRepository(s.mongoDB.DB, s.logger)
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %v", msg, err))
	}
}
