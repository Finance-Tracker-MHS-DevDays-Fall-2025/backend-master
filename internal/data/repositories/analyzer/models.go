package analyzer

import (
	"time"

	"backend-master/internal/api-gen/proto/common"
	pb "backend-master/internal/api-gen/proto/analyzer"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PeriodBalance struct {
	PeriodStart time.Time
	PeriodEnd   time.Time
	Income      int64
	Expense     int64
	Balance     int64
	Currency    string
}

type CategorySpending struct {
	CategoryID  string
	TotalAmount int64
	Currency    string
}

type Forecast struct {
	PeriodStart     time.Time
	PeriodEnd       time.Time
	ExpectedIncome  int64
	ExpectedExpense int64
	ExpectedBalance int64
	Currency        string
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

