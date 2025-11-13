package wallet

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/wallet"
	"backend-master/internal/data/repositories/wallet"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletController interface {
	GetUserAccounts(
		ctx context.Context,
		userID string,
	) (*pb.GetAccountsResponse, error)

	GetAccountTransactions(
		ctx context.Context,
		accountID string,
		limit int32,
	) (*pb.GetTransactionsResponse, error)

	CreateTransaction(
		ctx context.Context,
		accountID string,
		toAccountID string,
		txType common.TransactionType,
		amount int64,
		currency string,
		mcc string,
		description string,
		date time.Time,
	) (*pb.Transaction, error)
}

type walletControllerImpl struct {
	repo   wallet.WalletRepository
	client *wallet.WalletClient
	logger *zap.Logger
}

func NewController(
	repo wallet.WalletRepository,
	client *wallet.WalletClient,
	logger *zap.Logger,
) WalletController {
	return &walletControllerImpl{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

func (cont *walletControllerImpl) GetUserAccounts(
	ctx context.Context,
	userID string,
) (*pb.GetAccountsResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	accounts, err := cont.repo.GetAccountsByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts from repository: %w", err)
	}

	pbAccounts := make([]*pb.Account, 0, len(accounts))
	for _, acc := range accounts {
		pbAccounts = append(
			pbAccounts,
			acc.ToProto(),
		)
	}

	return &pb.GetAccountsResponse{
		Accounts: pbAccounts,
	}, nil
}

func (cont *walletControllerImpl) GetAccountTransactions(
	ctx context.Context,
	accountID string,
	limit int32,
) (*pb.GetTransactionsResponse, error) {
	aid, err := uuid.Parse(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	transactions, err := cont.repo.GetTransactionsByAccountID(ctx, aid, int(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions from repository: %w", err)
	}

	pbTransactions := make([]*pb.Transaction, 0, len(transactions))
	for _, tx := range transactions {
		pbTx := tx.ToProto()

		if tx.ToAccountID.Valid {
			pbTx.ToAccountId = tx.ToAccountID.String
		}

		if tx.Description.Valid {
			pbTx.Description = tx.Description.String
		}

		pbTransactions = append(pbTransactions, pbTx)
	}

	return &pb.GetTransactionsResponse{
		Transactions: pbTransactions,
	}, nil
}

func (cont *walletControllerImpl) CreateTransaction(
	ctx context.Context,
	accountID string,
	toAccountID string,
	txType common.TransactionType,
	amount int64,
	currency string,
	mcc string,
	description string,
	date time.Time,
) (*pb.Transaction, error) {
	aid, err := uuid.Parse(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	txTypeStr := wallet.TransactionPbTypeToDbType(txType)

	tx := &wallet.Transaction{
		AccountID: aid,
		Type:      txTypeStr,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: date,
	}

	if toAccountID != "" {
		tx.ToAccountID = sql.NullString{String: toAccountID, Valid: true}
	}

	if mcc != "" {
		mccInt, err := strconv.Atoi(mcc)
		if err == nil {
			tx.MCC = sql.NullInt32{Int32: int32(mccInt), Valid: true}
		}
	}

	if description != "" {
		tx.Description = sql.NullString{String: description, Valid: true}
	}

	createdTx, err := cont.repo.CreateTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction in repository: %w", err)
	}

	balanceChange := amount

	switch txType {
	case common.TransactionType_TRANSACTION_TYPE_EXPENSE:
		balanceChange = -amount
	case common.TransactionType_TRANSACTION_TYPE_TRANSFER:
		balanceChange = -amount

		if toAccountID != "" {
			toAid, parseErr := uuid.Parse(toAccountID)
			if parseErr == nil {
				err = cont.repo.UpdateAccountBalance(ctx, toAid, amount)

				if err != nil {
					cont.logger.Error(
						"failed to update account balance",
						zap.Error(err),
					)
				}
			}
		}
	}

	err = cont.repo.UpdateAccountBalance(ctx, aid, balanceChange)
	if err != nil {
		cont.logger.Error(
			"failed to update account balance",
			zap.Error(err),
		)
	}

	return createdTx.ToProto(), nil
}
