package presentation

import (
	"context"
	"fmt"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/master"
	anal "backend-master/internal/domain/controllers/analyzer"
	"backend-master/internal/domain/controllers/market"
	"backend-master/internal/domain/controllers/wallet"

	"go.uber.org/zap"
)

type masterServiceImpl struct {
	pb.UnimplementedMasterServiceServer
	logger       *zap.Logger
	walletCtrl   wallet.WalletController
	marketCtrl   market.MarketController
	analyzerCtrl anal.AnalyzerController
}

func NewMasterService(
	logger *zap.Logger,
	walletCtrl wallet.WalletController,
	marketCtrl market.MarketController,
	analyzerCtrl anal.AnalyzerController,
) pb.MasterServiceServer {
	return &masterServiceImpl{
		logger:       logger,
		walletCtrl:   walletCtrl,
		marketCtrl:   marketCtrl,
		analyzerCtrl: analyzerCtrl,
	}
}

func (s *masterServiceImpl) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	s.logger.Info("CreateTransaction", zap.String("body", fmt.Sprintf("%v", req)))

	tx, err := s.walletCtrl.CreateTransaction(
		ctx,
		req.FromAccountId,
		req.ToAccountId,
		req.Type,
		req.Amount.Amount,
		req.Amount.Currency,
		req.CategoryId, // mcc
		req.Description,
		req.Date.AsTime(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return &pb.CreateTransactionResponse{
		Transaction: tx,
	}, nil
}

func (s *masterServiceImpl) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.GetTransactionsResponse, error) {
	s.logger.Info("GetBalance", zap.String("body", fmt.Sprintf("%v", req)))

	resp, err := s.walletCtrl.GetUserTransactions(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	return &pb.GetTransactionsResponse{
		Transactions: resp.Transactions,
	}, nil
}

func (s *masterServiceImpl) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	s.logger.Info("GetBalance", zap.String("body", fmt.Sprintf("%v", req)))

	accountsResp, err := s.walletCtrl.GetUserAccounts(ctx, req.UserId)
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
	s.logger.Info("GetAnalytics", zap.String("body", fmt.Sprintf("%v", req)))

	stats, err := s.analyzerCtrl.GetStatistics(
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
	s.logger.Info("GetForecast", zap.String("body", fmt.Sprintf("%v", req)))

	forecast, err := s.analyzerCtrl.GetForecast(
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
