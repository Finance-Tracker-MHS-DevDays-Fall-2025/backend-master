package presentation

import (
	"context"

	pb "backend-master/internal/api/generated/proto/master"

	"go.uber.org/zap"
)

type MasterService struct {
	pb.UnimplementedMasterServiceServer
	logger *zap.Logger
}

func NewMasterService(logger *zap.Logger) *MasterService {
	return &MasterService{
		logger: logger,
	}
}

func (s *MasterService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	s.logger.Info("CreateTransaction", zap.String("user_id", req.UserId))
	return &pb.CreateTransactionResponse{}, nil
}

func (s *MasterService) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	s.logger.Info("GetBalance", zap.String("user_id", req.UserId))
	return &pb.GetBalanceResponse{}, nil
}

func (s *MasterService) GetAnalytics(ctx context.Context, req *pb.GetAnalyticsRequest) (*pb.GetAnalyticsResponse, error) {
	s.logger.Info("GetAnalytics", zap.String("user_id", req.UserId))
	return &pb.GetAnalyticsResponse{}, nil
}

func (s *MasterService) GetForecast(ctx context.Context, req *pb.GetForecastRequest) (*pb.GetForecastResponse, error) {
	s.logger.Info("GetForecast", zap.String("user_id", req.UserId))
	return &pb.GetForecastResponse{}, nil
}
