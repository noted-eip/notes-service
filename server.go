package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"notes-service/auth"
	"notes-service/language"

	"notes-service/communication"

	background "github.com/noted-eip/noted/background-service"
	mailing "github.com/noted-eip/noted/mailing-service"

	"context"
	"errors"
	"fmt"
	"net"
	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"strings"
	"time"

	"notes-service/models/mongo"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	logger  *zap.Logger
	slogger *zap.SugaredLogger

	authService       auth.Service
	backgroundService background.Service
	mailingService    mailing.Service
	languageService   language.Service // NOTE: Could put directly service typed as NaturalAPIService, remove Init() from interface and just put it in NaturalAPIService
	accountsClient    *communication.AccountsServiceClient

	mongoDB *mongo.Database

	notesRepository      models.NotesRepository
	groupsRepository     models.GroupsRepository
	activitiesRepository models.ActivitiesRepository

	notesAPI           notesv1.NotesAPIServer
	groupsAPI          notesv1.GroupsAPIServer
	recommendationsAPI notesv1.RecommendationsAPIServer

	grpcServer *grpc.Server
}

func (s *server) Init(opt ...grpc.ServerOption) {
	s.initLogger()
	s.initAuthService()
	s.initRepositories()
	s.initLanguageService()
	s.initBackgroundService()
	s.initMailingService()
	s.initAccountsClient()
	s.initGroupsAPI()
	s.initNotesAPI()
	s.initRecommendationsAPI()
	s.initgrpcServer(opt...)

	s.validateOldBackgroundService()
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
	s.languageService = &language.NotedLanguageService{}
	err := s.languageService.Init(s.logger)
	must(err, "unable to instantiate language service")
}

func (s *server) initAuthService() {
	rawKey, err := base64.StdEncoding.DecodeString(*jwtPrivateKey)
	must(err, "could not decode jwt private key")
	var key ed25519.PrivateKey = rawKey
	pubKey := key.Public().(ed25519.PublicKey)
	s.authService = auth.NewService(pubKey)
}

func (s *server) initBackgroundService() {
	s.backgroundService = background.NewService(s.logger)
}

func (s *server) initMailingService() {
	s.mailingService = mailing.NewService(s.logger, *gmailSuperSecret)
}

func (s *server) initAccountsClient() {
	accountsClient, err := communication.NewAccountsServiceClient(*accountsServiceUrl)
	if *environment == envIsDev && err != nil {
		s.logger.Warn(fmt.Sprintf("could not instantiate accounts service connection: %v", err))
		accountsClient = nil
	} else {
		must(err, "could not instantiate accounts service connection")
	}
	s.accountsClient = accountsClient
}

func (s *server) initGroupsAPI() {
	s.groupsAPI = &groupsAPI{
		auth:           s.authService,
		logger:         s.logger,
		notes:          s.notesRepository,
		groups:         s.groupsRepository,
		activities:     s.activitiesRepository,
		background:     s.backgroundService,
		mailing:        s.mailingService,
		accountsClient: s.accountsClient,
	}
}

func (s *server) initNotesAPI() {
	s.notesAPI = &notesAPI{
		auth:       s.authService,
		logger:     s.logger,
		notes:      s.notesRepository,
		groups:     s.groupsRepository,
		activities: s.activitiesRepository,
		language:   s.languageService,
		background: s.backgroundService,
	}
}

func (s *server) initRecommendationsAPI() {
	s.recommendationsAPI = &recommendationsAPI{
		auth:     s.authService,
		logger:   s.logger,
		notes:    s.notesRepository,
		language: s.languageService,
	}
}

func (s *server) initgrpcServer(opt ...grpc.ServerOption) {
	s.grpcServer = grpc.NewServer(opt...)
	notesv1.RegisterNotesAPIServer(s.grpcServer, s.notesAPI)
	notesv1.RegisterGroupsAPIServer(s.grpcServer, s.groupsAPI)
	notesv1.RegisterRecommendationsAPIServer(s.grpcServer, s.recommendationsAPI)
}

func (s *server) initRepositories() {
	var err error
	s.mongoDB, err = mongo.NewDatabase(context.Background(), *mongoUri, *mongoDbName, s.logger)
	must(err, "could not instantiate mongo database")
	s.notesRepository = mongo.NewNotesRepository(s.mongoDB.DB, s.logger)
	s.groupsRepository = mongo.NewGroupsRepository(s.mongoDB.DB, s.logger)
	s.activitiesRepository = mongo.NewActivitiesRepository(s.mongoDB.DB, s.logger)
}

func (s *server) validateOldBackgroundService() {
	s.validateQuizsExpiration()
}

func (s *server) validateQuizsExpiration() {
	res, err := s.notesRepository.ListQuizsCreatedDateInternal(context.Background())

	if err != nil {
		s, ok := status.FromError(err)
		if !ok {
			must(errors.New("wrong error code"), "not a grpc status")
		}
		must(err, "could not check quiz listing error")
		if s.Code() != codes.NotFound {
			must(err, "could not list quizs for expiration validation")
		}
	}

	for _, quiz := range *res {
		err := s.backgroundService.AddProcess(&background.Process{
			Identifier: models.NoteIdentifier{
				ActionType: models.NoteDeleteQuiz,
				Metadata: models.QuizIdentifier{
					QuizID: quiz.ID,
				},
			},
			CallBackFct: func() error {
				return s.notesRepository.DeleteQuizFromIDInternal(context.Background(), quiz.ID)
			},
			CancelProcessOnSameIdentifier: true,
			RepeatProcess:                 false,
			SecondsToDebounce:             uint32(time.Until(quiz.CreatedAt.Add(time.Hour * 24 * 7)).Seconds()),
		})
		must(err, "couldn't validate quiz "+quiz.ID+"'s expiration date")
	}
}

func must(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %v", msg, err))
	}
}
