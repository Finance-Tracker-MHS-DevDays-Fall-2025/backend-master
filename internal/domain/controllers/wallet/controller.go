package wallet

import (
	"context"
	"fmt"

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
		cont.logger.Error(
			"invalid user ID",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	accounts, err := cont.repo.GetAccountsByUserID(ctx, uid)
	if err != nil {
		cont.logger.Error(
			"failed to get accounts from repository",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get accounts: %w", err)
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
		cont.logger.Error(
			"invalid account ID",
			zap.Error(err),
			zap.String("account_id", accountID),
		)
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	transactions, err := cont.repo.GetTransactionsByAccountID(ctx, aid, int(limit))
	if err != nil {
		cont.logger.Error(
			"failed to get transactions from repository",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
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
