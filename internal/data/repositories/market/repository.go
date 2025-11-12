package market

import (
	"backend-master/internal/data/database"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MarketRepository interface {
	GetInvestmentPositionsByAccountID(
		ctx context.Context,
		accountID uuid.UUID,
	) ([]InvestmentPosition, error)

	GetSecurityByFIGI(
		ctx context.Context,
		figi string,
	) (*Security, error)

	GetSecuritiesByFIGIs(
		ctx context.Context,
		figis []string,
	) ([]Security, error)

	GetSecurityPaymentsByFIGI(
		ctx context.Context,
		figi string,
	) ([]SecurityPayment, error)
}

type marketRepositoryImpl struct {
	db     database.DBManager
	logger *zap.Logger
}

func NewRepository(
	db database.DBManager,
	logger *zap.Logger,
) MarketRepository {
	return &marketRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (repo *marketRepositoryImpl) GetInvestmentPositionsByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
) ([]InvestmentPosition, error) {
	query := `
		SELECT 
			id,
			account_id,
			security_id,
			quantity,
			created_at
		FROM investment_positions
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	var positions []InvestmentPosition
	err := repo.db.GetDB().SelectContext(ctx, &positions, query, accountID)
	if err != nil {
		repo.logger.Error(
			"failed to get investment positions",
			zap.Error(err),
			zap.String("account_id", accountID.String()),
		)
		return nil, fmt.Errorf("failed to get investment positions: %w", err)
	}

	return positions, nil
}

func (repo *marketRepositoryImpl) GetSecurityByFIGI(
	ctx context.Context,
	figi string,
) (*Security, error) {
	query := `
		SELECT 
			figi,
			name,
			current_price,
			type,
			price_updated_at,
			created_at
		FROM securities
		WHERE figi = $1
	`

	var security Security
	err := repo.db.GetDB().GetContext(ctx, &security, query, figi)
	if err != nil {
		repo.logger.Error(
			"failed to get security",
			zap.Error(err),
			zap.String("figi", figi),
		)
		return nil, fmt.Errorf("failed to get security: %w", err)
	}

	return &security, nil
}

func (repo *marketRepositoryImpl) GetSecuritiesByFIGIs(
	ctx context.Context,
	figis []string,
) ([]Security, error) {
	query := `
		SELECT 
			figi,
			name,
			current_price,
			type,
			price_updated_at,
			created_at
		FROM securities
		WHERE figi = ANY($1)
	`

	var securities []Security
	err := repo.db.GetDB().SelectContext(ctx, &securities, query, figis)
	if err != nil {
		repo.logger.Error(
			"failed to get securities",
			zap.Error(err),
			zap.Int("count", len(figis)),
		)
		return nil, fmt.Errorf("failed to get securities: %w", err)
	}

	return securities, nil
}

func (repo *marketRepositoryImpl) GetSecurityPaymentsByFIGI(
	ctx context.Context,
	figi string,
) ([]SecurityPayment, error) {
	query := `
		SELECT 
			id,
			security_id,
			amount_per_share,
			payment_date,
			created_at
		FROM securities_payments
		WHERE security_id = $1
		ORDER BY payment_date DESC
	`

	var payments []SecurityPayment
	err := repo.db.GetDB().SelectContext(ctx, &payments, query, figi)
	if err != nil {
		repo.logger.Error(
			"failed to get security payments",
			zap.Error(err),
			zap.String("figi", figi),
		)
		return nil, fmt.Errorf("failed to get security payments: %w", err)
	}

	return payments, nil
}

