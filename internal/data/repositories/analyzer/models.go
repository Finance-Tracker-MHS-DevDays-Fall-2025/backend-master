package analyzer

import (
	"time"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/analyzer"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PeriodBalance struct {
	PeriodStart time.Time `db:"period_start"`
	PeriodEnd   time.Time `db:"period_end"`
	Income      int64     `db:"income"`
	Expense     int64     `db:"expense"`
	Balance     int64     `db:"balance"`
	Currency    string    `db:"currency"`
}

type CategorySpending struct {
	CategoryID  string `db:"category_id"`
	TotalAmount int64  `db:"total_amount"`
	Currency    string `db:"currency"`
}

type Forecast struct {
	PeriodStart     time.Time `db:"period_start"`
	PeriodEnd       time.Time `db:"period_end"`
	ExpectedIncome  int64     `db:"expected_income"`
	ExpectedExpense int64     `db:"expected_expense"`
	ExpectedBalance int64     `db:"expected_balance"`
	Currency        string    `db:"currency"`
}

func (p *PeriodBalance) ToProto(categoryBreakdown []*pb.CategorySpending) *pb.PeriodBalance {
	return &pb.PeriodBalance{
		PeriodStart: timestamppb.New(p.PeriodStart),
		PeriodEnd:   timestamppb.New(p.PeriodEnd),
		Income: &common.Money{
			Amount:   p.Income,
			Currency: p.Currency,
		},
		Expense: &common.Money{
			Amount:   p.Expense,
			Currency: p.Currency,
		},
		Balance: &common.Money{
			Amount:   p.Balance,
			Currency: p.Currency,
		},
		CategoryBreakdown: categoryBreakdown,
	}
}

func (cs *CategorySpending) ToProto() *pb.CategorySpending {
	return &pb.CategorySpending{
		CategoryId: cs.CategoryID,
		TotalAmount: &common.Money{
			Amount:   cs.TotalAmount,
			Currency: cs.Currency,
		},
	}
}

func (f *Forecast) ToProto(categoryBreakdown []*pb.CategorySpending) *pb.Forecast {
	return &pb.Forecast{
		PeriodStart: timestamppb.New(f.PeriodStart),
		PeriodEnd:   timestamppb.New(f.PeriodEnd),
		ExpectedIncome: &common.Money{
			Amount:   f.ExpectedIncome,
			Currency: f.Currency,
		},
		ExpectedExpense: &common.Money{
			Amount:   f.ExpectedExpense,
			Currency: f.Currency,
		},
		ExpectedBalance: &common.Money{
			Amount:   f.ExpectedBalance,
			Currency: f.Currency,
		},
		CategoryBreakdown: categoryBreakdown,
	}
}

