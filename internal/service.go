package internal

import (
	"backend-master/configs"
	pb "backend-master/internal/api-gen/proto/master"
	"backend-master/internal/presentation"
	"backend-master/internal/presentation/docs"
	"context"
	_ "embed"
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
	Start(httpAddr, grpcAddr string) error
	Shutdown() error
}

type serviceImpl struct {
	Service

	grpcServer *grpc.Server
	ginEngine  *gin.Engine
	logger     *zap.Logger
}

func NewService(
	cfg *configs.ServiceConfig,
	logger *zap.Logger,
) Service {
	gin.SetMode(gin.ReleaseMode)

	grpcServer := grpc.NewServer()
	masterService := presentation.NewMasterService(logger)
	pb.RegisterMasterServiceServer(grpcServer, masterService)

	s := &serviceImpl{
		grpcServer: grpcServer,
		ginEngine:  gin.New(),
		logger:     logger,
	}

	return s
}

func (s *serviceImpl) Start(httpAddr, grpcAddr string) error {
	ctx := context.Background()

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
	}

	err := pb.RegisterMasterServiceHandlerFromEndpoint(ctx, grpcMux, grpcAddr, opts)
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
