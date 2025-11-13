package internal

import (
	"backend-master/configs"
	pb "backend-master/internal/api-gen/proto/master"
	"backend-master/internal/data/database"
	analRepo "backend-master/internal/data/repositories/analyzer"
	walletRepo "backend-master/internal/data/repositories/wallet"
	analyzerController "backend-master/internal/domain/controllers/analyzer"
	walletController "backend-master/internal/domain/controllers/wallet"
	"backend-master/internal/presentation"
	"backend-master/internal/presentation/docs"
	"context"
	_ "embed"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed api-gen/openapi/master/master.swagger.json
var swaggerJSON []byte

type Service interface {
	Start() error
	Shutdown() error
}

type serviceImpl struct {
	Service

	cfg        *configs.ServiceConfig
	grpcServer *grpc.Server
	ginEngine  *gin.Engine
	logger     *zap.Logger
}

func NewService(
	cfg *configs.ServiceConfig,
	logger *zap.Logger,
) Service {
	gin.SetMode(gin.ReleaseMode)

	dbManager, err := database.NewManager(cfg.DatabaseCfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}

	walletRepository := walletRepo.NewRepository(dbManager, logger)
	analyzerRepository := analRepo.NewRepository(dbManager, logger)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(presentation.UnaryClientInterceptor(logger)),
	}

	walletClient, err := walletRepo.NewClient(
		cfg.SlavesCfg.WalletUrl,
		logger,
		opts...,
	)
	if err != nil {
		logger.Fatal("failed to initialize wallet client", zap.Error(err))
	}

	analyzerClient, err := analRepo.NewClient(
		cfg.SlavesCfg.AnalyzerUrl,
		logger,
		opts...,
	)
	if err != nil {
		logger.Fatal("failed to initialize analyzer client", zap.Error(err))
	}

	walletCtrl := walletController.NewController(walletRepository, walletClient, logger)
	analyzerCtrl := analyzerController.NewController(analyzerRepository, analyzerClient, logger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(presentation.UnaryServerInterceptor(logger)),
	)
	masterService := presentation.NewMasterService(logger, walletCtrl, analyzerCtrl)
	pb.RegisterMasterServiceServer(grpcServer, masterService)

	s := &serviceImpl{
		cfg:        cfg,
		grpcServer: grpcServer,
		ginEngine:  gin.New(),
		logger:     logger,
	}

	return s
}

func (s *serviceImpl) Start() error {
	ctx := context.Background()

	grpcAddr := fmt.Sprintf(
		"%s:%d",
		s.cfg.ServerCfg.ServerHost,
		s.cfg.ServerCfg.GrpcPort,
	)
	httpAddr := fmt.Sprintf(
		"%s:%d",
		s.cfg.ServerCfg.ServerHost,
		s.cfg.ServerCfg.HttpPort,
	)

	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			s.logger.Fatal("failed to listen gRPC", zap.Error(err))
		}

		s.logger.Info("starting gRPC server", zap.String("addr", grpcAddr))

		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(presentation.UnaryClientInterceptor(s.logger)),
	}

	err := pb.RegisterMasterServiceHandlerFromEndpoint(
		ctx,
		grpcMux,
		grpcAddr,
		opts,
	)
	if err != nil {
		s.logger.Fatal("failed to register gateway", zap.Error(err))
	}

	s.ginEngine.Use(gin.Recovery())

	apiRouter := s.ginEngine.Group("/api")
	apiRouter.GET("/docs", docs.NewSwaggerHandler(swaggerJSON))

	apiV1Router := apiRouter.Group("/v1")
	apiV1Router.Any("/*path", gin.WrapH(grpcMux))

	s.logger.Info("starting HTTP server", zap.String("addr", httpAddr))

	return s.ginEngine.Run(httpAddr)
}

func (s *serviceImpl) Shutdown() error {
	s.logger.Info("shutting down servers")
	s.grpcServer.GracefulStop()
	return nil
}
