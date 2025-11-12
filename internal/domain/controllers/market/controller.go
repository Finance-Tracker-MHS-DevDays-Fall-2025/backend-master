package market

import (
	"context"
	"fmt"

	pb "backend-master/internal/api-gen/proto/market"
	"backend-master/internal/data/repositories/market"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MarketController interface {
	GetInvestmentPositions(
		ctx context.Context,
		accountID string,
	) (*pb.GetInvestmentPositionsResponse, error)

	GetSecurity(
		ctx context.Context,
		figi string,
	) (*pb.GetSecurityResponse, error)

	GetSecuritiesPrices(
		ctx context.Context,
		figis []string,
	) (*pb.GetSecuritiesPricesResponse, error)

	GetSecurityPayments(
		ctx context.Context,
		figi string,
	) (*pb.GetSecuritiesPaymentsResponse, error)
}

type marketControllerImpl struct {
	repo   market.MarketRepository
	client *market.MarketClient
	logger *zap.Logger
}

func NewController(
	repo market.MarketRepository,
	client *market.MarketClient,
	logger *zap.Logger,
) MarketController {
	return &marketControllerImpl{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

func (cont *marketControllerImpl) GetInvestmentPositions(
	ctx context.Context,
	accountID string,
) (*pb.GetInvestmentPositionsResponse, error) {
	aid, err := uuid.Parse(accountID)
	if err != nil {
		cont.logger.Error(
			"invalid account ID",
			zap.Error(err),
			zap.String("account_id", accountID),
		)
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	positions, err := cont.repo.GetInvestmentPositionsByAccountID(ctx, aid)
	if err != nil {
		cont.logger.Error(
			"failed to get investment positions from repository",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get investment positions: %w", err)
	}

	pbPositions := make([]*pb.InvestmentPosition, 0, len(positions))
	for _, pos := range positions {
		pbPositions = append(pbPositions, pos.ToProto())
	}

	return &pb.GetInvestmentPositionsResponse{
		Positions: pbPositions,
	}, nil
}

func (cont *marketControllerImpl) GetSecurity(
	ctx context.Context,
	figi string,
) (*pb.GetSecurityResponse, error) {
	security, err := cont.repo.GetSecurityByFIGI(ctx, figi)
	if err != nil {
		cont.logger.Error(
			"failed to get security from repository",
			zap.Error(err),
			zap.String("figi", figi),
		)
		return nil, fmt.Errorf("failed to get security: %w", err)
	}

	return &pb.GetSecurityResponse{
		Security: security.ToProto(),
	}, nil
}

func (cont *marketControllerImpl) GetSecuritiesPrices(
	ctx context.Context,
	figis []string,
) (*pb.GetSecuritiesPricesResponse, error) {
	securities, err := cont.repo.GetSecuritiesByFIGIs(ctx, figis)
	if err != nil {
		cont.logger.Error(
			"failed to get securities from repository",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get securities: %w", err)
	}

	pbSecurities := make([]*pb.Security, 0, len(securities))
	for _, sec := range securities {
		pbSecurities = append(pbSecurities, sec.ToProto())
	}

	return &pb.GetSecuritiesPricesResponse{
		Securities: pbSecurities,
	}, nil
}

func (cont *marketControllerImpl) GetSecurityPayments(
	ctx context.Context,
	figi string,
) (*pb.GetSecuritiesPaymentsResponse, error) {
	payments, err := cont.repo.GetSecurityPaymentsByFIGI(ctx, figi)
	if err != nil {
		cont.logger.Error(
			"failed to get security payments from repository",
			zap.Error(err),
			zap.String("figi", figi),
		)
		return nil, fmt.Errorf("failed to get security payments: %w", err)
	}

	pbPayments := make([]*pb.SecurityPayment, 0, len(payments))
	for _, pay := range payments {
		pbPayments = append(pbPayments, pay.ToProto())
	}

	return &pb.GetSecuritiesPaymentsResponse{
		Payments: pbPayments,
	}, nil
}
