package wallet

import (
	"backend-master/internal/data/database"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletRepository interface {
	GetAccountsByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]Account, error)

	GetTransactionsByAccountID(
		ctx context.Context,
		accountID uuid.UUID,
		limit int,
	) ([]Transaction, error)
}

type walletRepositoryImpl struct {
	db     database.DBManager
	logger *zap.Logger
}

func NewRepository(
	db database.DBManager,
	logger *zap.Logger,
) WalletRepository {
	return &walletRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (repo *walletRepositoryImpl) GetAccountsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]Account, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			name, 
			type, 
			balance, 
			currency, 
			created_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var accounts []Account
	err := repo.db.GetDB().SelectContext(ctx, &accounts, query, userID)
	if err != nil {
		repo.logger.Error(
			"failed to get accounts",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (repo *walletRepositoryImpl) GetTransactionsByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
) ([]Transaction, error) {
	query := `
		SELECT 
			id,
			account_id,
			to_account_id,
			type,
			amount,
			currency,
			mcc,
			description,
			created_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	var transactions []Transaction
	err := repo.db.GetDB().SelectContext(ctx, &transactions, query, accountID, limit)
	if err != nil {
		repo.logger.Error(
			"failed to get transactions",
			zap.Error(err),
			zap.String("account_id", accountID.String()),
			zap.Int("limit", limit),
		)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
