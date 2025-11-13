package market

import (
	"context"
	"fmt"
	"time"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/market"
	"backend-master/internal/data/repositories/market"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	client *market.MarketClient
	logger *zap.Logger
}

func NewController(
	client *market.MarketClient,
	logger *zap.Logger,
) MarketController {
	return &marketControllerImpl{
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
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	positions, err := cont.client.GetInvestmentPositions(
		ctx,
		&pb.GetInvestmentPositionsRequest{
			AccountId: aid.String(),
			UserId:    "uid",
			Backend: &common.AccountBackend{
				Type:      "TInvest",
				AccountId: "aid",
				Token:     "tok",
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get investment positions from client: %w", err)
	}

	return positions, nil
}

func (cont *marketControllerImpl) GetSecurity(
	ctx context.Context,
	figi string,
) (*pb.GetSecurityResponse, error) {
	security, err := cont.client.GetSecurity(
		ctx,
		&pb.GetSecurityRequest{Figi: figi},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get security from client: %w", err)
	}
	return security, nil
}

func (cont *marketControllerImpl) GetSecuritiesPrices(
	ctx context.Context,
	figis []string,
) (*pb.GetSecuritiesPricesResponse, error) {
	securities, err := cont.client.GetSecuritiesPrices(
		ctx,
		&pb.GetSecuritiesPricesRequest{Figis: figis},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get securities from client: %w", err)
	}
	return securities, nil
}

func (cont *marketControllerImpl) GetSecurityPayments(
	ctx context.Context,
	figi string,
) (*pb.GetSecuritiesPaymentsResponse, error) {
	payments, err := cont.client.GetSecurityPayments(
		ctx,
		&pb.GetSecuritiesPaymentsRequest{
			Figis:     []string{figi},
			StartDate: timestamppb.New(time.Now().AddDate(0, -6, 0)),
			EndDate:   timestamppb.New(time.Now().AddDate(0, 6, 0)),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get securities from client: %w", err)
	}
	return payments, nil
}
