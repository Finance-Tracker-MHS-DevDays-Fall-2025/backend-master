package wallet

import (
	"database/sql"
	"time"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/wallet"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Account struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Balance   int64     `db:"balance"` // копейки
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}

type Transaction struct {
	ID          uuid.UUID      `db:"id"`
	AccountID   uuid.UUID      `db:"account_id"`
	ToAccountID sql.NullString `db:"to_account_id"`
	Type        string         `db:"type"`
	Amount      int64          `db:"amount"` // копейки
	Currency    string         `db:"currency"`
	MCC         sql.NullInt32  `db:"mcc"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
}

func (acc *Account) ToProto() *pb.Account {
	money := &common.Money{
		Amount:   acc.Balance,
		Currency: acc.Currency,
	}

	return &pb.Account{
		AccountId: acc.ID.String(),
		UserId:    acc.UserID.String(),
		Name:      acc.Name,
		Type:      common.AccountType(common.AccountType_value[acc.Type]),
		Balance:   money,
		CreatedAt: timestamppb.New(acc.CreatedAt),
	}
}

func (tx *Transaction) ToProto() *pb.Transaction {
	money := &common.Money{
		Amount:   tx.Amount,
		Currency: tx.Currency,
	}

	return &pb.Transaction{
		AccountId: tx.AccountID.String(),
		Type:      common.TransactionType(common.TransactionType_value[tx.Type]),
		Amount:    money,
		Date:      timestamppb.New(tx.CreatedAt),
	}
}
