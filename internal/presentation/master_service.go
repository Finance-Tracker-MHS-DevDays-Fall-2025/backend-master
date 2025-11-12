package presentation

import (
	"context"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/master"
	anal "backend-master/internal/domain/controllers/analyzer"
	"backend-master/internal/domain/controllers/wallet"

	"go.uber.org/zap"
)

type masterServiceImpl struct {
	pb.UnimplementedMasterServiceServer
	logger             *zap.Logger
	walletController   wallet.WalletController
	analyzerController anal.AnalyzerController
}

func NewMasterService(
	logger *zap.Logger,
	walletCtrl wallet.WalletController,
	analyzerCtrl anal.AnalyzerController,
) pb.MasterServiceServer {
	return &masterServiceImpl{
		logger:             logger,
		walletController:   walletCtrl,
		analyzerController: analyzerCtrl,
	}
}

func (s *masterServiceImpl) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	s.logger.Info("CreateTransaction", zap.String("user_id", req.UserId))
	// TODO: implement transaction creation
	return &pb.CreateTransactionResponse{}, nil
}

func (s *masterServiceImpl) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	s.logger.Info("GetBalance", zap.String("user_id", req.UserId))

	accountsResp, err := s.walletController.GetUserAccounts(ctx, req.UserId)
	if err != nil {
		s.logger.Error("failed to get user accounts", zap.Error(err))
		return nil, err
	}

	var totalBalance int64
	for _, account := range accountsResp.Accounts {
		if account.Balance != nil {
			totalBalance += account.Balance.Amount
		}
	}

	return &pb.GetBalanceResponse{
		TotalBalance: &common.Money{
			Amount:   totalBalance,
			Currency: "RUB",
		},
		Accounts: accountsResp.Accounts,
	}, nil
}

func (s *masterServiceImpl) GetAnalytics(ctx context.Context, req *pb.GetAnalyticsRequest) (*pb.GetAnalyticsResponse, error) {
	s.logger.Info("GetAnalytics", zap.String("user_id", req.UserId))

	stats, err := s.analyzerController.GetStatistics(
		ctx,
		req.UserId,
		req.StartDate.AsTime(),
		req.EndDate.AsTime(),
		common.TimePeriod_TIME_PERIOD_MONTH,
	)
	if err != nil {
		s.logger.Error("failed to get statistics", zap.Error(err))
		return nil, err
	}

	return &pb.GetAnalyticsResponse{
		Statistics: stats,
	}, nil
}

func (s *masterServiceImpl) GetForecast(ctx context.Context, req *pb.GetForecastRequest) (*pb.GetForecastResponse, error) {
	s.logger.Info("GetForecast", zap.String("user_id", req.UserId))

	forecast, err := s.analyzerController.GetForecast(
		ctx,
		req.UserId,
		req.Period,
		req.PeriodsAhead,
	)
	if err != nil {
		s.logger.Error("failed to get forecast", zap.Error(err))
		return nil, err
	}

	return &pb.GetForecastResponse{
		Forecasts: forecast.Forecasts,
	}, nil
}
