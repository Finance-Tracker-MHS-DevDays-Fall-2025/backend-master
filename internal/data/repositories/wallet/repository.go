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

	CreateTransaction(
		ctx context.Context,
		tx *Transaction,
	) (*Transaction, error)

	UpdateAccountBalance(
		ctx context.Context,
		accountID uuid.UUID,
		amount int64,
	) error
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

		WHERE 1=1
			AND user_id = $1

		ORDER BY created_at DESC
	`

	var accounts []Account
	err := repo.db.GetDB().SelectContext(ctx, &accounts, query, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get accounts for uid: %s %w",
			userID.String(),
			err,
		)
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

		WHERE 1=1
			AND account_id = $1

		ORDER BY created_at DESC
		LIMIT $2
	`

	var transactions []Transaction
	err := repo.db.GetDB().SelectContext(ctx, &transactions, query, accountID, limit)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get transactions for aid %s: %w",
			accountID.String(),
			err,
		)
	}

	return transactions, nil
}

func (repo *walletRepositoryImpl) CreateTransaction(
	ctx context.Context,
	tx *Transaction,
) (*Transaction, error) {
	query := `
		INSERT INTO transactions (
			id,
			account_id,
			to_account_id,
			type,
			amount,
			currency,
			mcc,
			description,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, account_id, to_account_id, type, amount, currency, mcc, description, created_at
	`

	tx.ID = uuid.New()

	err := repo.db.GetDB().GetContext(
		ctx,
		tx,
		query,
		tx.ID,
		tx.AccountID,
		tx.ToAccountID,
		tx.Type,
		tx.Amount,
		tx.Currency,
		tx.MCC,
		tx.Description,
		tx.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create transaction for aid %s: %w",
			tx.AccountID.String(),
			err,
		)
	}

	return tx, nil
}

func (repo *walletRepositoryImpl) UpdateAccountBalance(
	ctx context.Context,
	accountID uuid.UUID,
	amount int64,
) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`

	_, err := repo.db.GetDB().ExecContext(ctx, query, amount, accountID)
	if err != nil {
		return fmt.Errorf(
			"failed to update account balance for aid %s and amount %d: %w",
			accountID.String(),
			amount,
			err,
		)
	}

	return nil
}
