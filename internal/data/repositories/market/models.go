package market

import (
	"database/sql"
	"time"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/market"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvestmentPosition struct {
	ID         uuid.UUID `db:"id"`
	AccountID  uuid.UUID `db:"account_id"`
	SecurityID uuid.UUID `db:"security_id"`
	Quantity   int32     `db:"quantity"`
	CreatedAt  time.Time `db:"created_at"`
}

type Security struct {
	Figi           string        `db:"figi"`
	Name           string        `db:"name"`
	CurrentPrice   sql.NullInt64 `db:"current_price"`
	Type           string        `db:"type"`
	PriceUpdatedAt sql.NullTime  `db:"price_updated_at"`
	CreatedAt      time.Time     `db:"created_at"`
}

type SecurityPayment struct {
	ID             uuid.UUID `db:"id"`
	SecurityID     string    `db:"security_id"`
	AmountPerShare int64     `db:"amount_per_share"`
	PaymentDate    time.Time `db:"payment_date"`
	CreatedAt      time.Time `db:"created_at"`
}

func (pos *InvestmentPosition) ToProto() *pb.InvestmentPosition {
	return &pb.InvestmentPosition{
		Figi:     pos.SecurityID.String(),
		Quantity: pos.Quantity,
	}
}

func (sec *Security) ToProto() *pb.Security {
	pbSec := &pb.Security{
		Figi: sec.Figi,
		Name: sec.Name,
	}

	if sec.CurrentPrice.Valid {
		pbSec.CurrentPrice = &common.Money{
			Amount:   sec.CurrentPrice.Int64,
			Currency: "RUB",
		}
	}

	if sec.PriceUpdatedAt.Valid {
		pbSec.PriceUpdatedAt = timestamppb.New(sec.PriceUpdatedAt.Time)
	}

	return pbSec
}

func (pay *SecurityPayment) ToProto() *pb.SecurityPayment {
	return &pb.SecurityPayment{
		Figi: pay.SecurityID,
		Payment: &common.Money{
			Amount:   pay.AmountPerShare,
			Currency: "RUB",
		},
		PaymentDate: timestamppb.New(pay.PaymentDate),
	}
}
