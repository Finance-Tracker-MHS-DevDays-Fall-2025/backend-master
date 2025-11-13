package analyzer

import (
	"backend-master/internal/data/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AnalyzerRepository interface {
	GetStatistics(
		ctx context.Context,
		userID uuid.UUID,
		startDate time.Time,
		endDate time.Time,
	) (totalIncome int64, totalExpense int64, err error)

	GetPeriodBalances(
		ctx context.Context,
		userID uuid.UUID,
		startDate time.Time,
		endDate time.Time,
		groupBy string,
	) ([]PeriodBalance, error)

	GetCategorySpending(
		ctx context.Context,
		userID uuid.UUID,
		startDate time.Time,
		endDate time.Time,
	) ([]CategorySpending, error)
}

type analyzerRepositoryImpl struct {
	db     database.DBManager
	logger *zap.Logger
}

func NewRepository(
	db database.DBManager,
	logger *zap.Logger,
) AnalyzerRepository {
	return &analyzerRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (repo *analyzerRepositoryImpl) GetStatistics(
	ctx context.Context,
	userID uuid.UUID,
	startDate time.Time,
	endDate time.Time,
) (totalIncome int64, totalExpense int64, err error) {
	query := `
		SELECT 
			COALESCE(
				SUM(
					CASE 
						WHEN type = 'INCOME' 
						THEN amount 
						ELSE 0 
					END
				),
				0
			) as total_income,
			COALESCE(
				SUM(
					CASE
						WHEN type = 'EXPENSE'
						THEN amount
						ELSE 0
					END
				),
				0
			) as total_expense
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.user_id = $1
			AND t.created_at BETWEEN $2 AND $3
	`

	var result struct {
		TotalIncome  int64 `db:"total_income"`
		TotalExpense int64 `db:"total_expense"`
	}

	err = repo.db.GetDB().GetContext(ctx, &result, query, userID, startDate, endDate)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"failed to get statistics for uid %s: %w",
			userID.String(),
			err,
		)
	}

	return result.TotalIncome, result.TotalExpense, nil
}

func (repo *analyzerRepositoryImpl) GetPeriodBalances(
	ctx context.Context,
	userID uuid.UUID,
	startDate time.Time,
	endDate time.Time,
	groupBy string,
) ([]PeriodBalance, error) {
	query := `
		SELECT 
			date_trunc($1, t.created_at) as period_start,
			date_trunc($1, t.created_at) + interval '1 ' || $1 as period_end,

			COALESCE(
				SUM(
					CASE 
						WHEN type = 'INCOME' THEN amount
						ELSE 0
					END
				),
				0
			) as income,
			
			COALESCE(
				SUM(
					CASE
						WHEN type = 'EXPENSE' THEN amount
						ELSE 0
					END
				),
				0
			) as expense,
			
			COALESCE(
				SUM(
					CASE
						WHEN type = 'INCOME' THEN amount
						WHEN type = 'EXPENSE' THEN -amount
						ELSE 0
					END
				),
				0
			) as balance,
			
			t.currency

		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.user_id = $2
			AND t.created_at BETWEEN $3 AND $4
		GROUP BY date_trunc($1, t.created_at), t.currency
		ORDER BY period_start
	`

	var balances []PeriodBalance
	err := repo.db.GetDB().SelectContext(ctx, &balances, query, groupBy, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get period balances for uid %s: %w",
			userID.String(),
			err,
		)
	}

	return balances, nil
}

func (repo *analyzerRepositoryImpl) GetCategorySpending(
	ctx context.Context,
	userID uuid.UUID,
	startDate time.Time,
	endDate time.Time,
) ([]CategorySpending, error) {
	query := `
		SELECT 
			COALESCE(
				t.category_id::text,
				'uncategorized'
			) as category_id,
			SUM(t.amount) as total_amount,
			t.currency
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.user_id = $1
			AND t.type = 'EXPENSE'
			AND t.created_at BETWEEN $2 AND $3
		GROUP BY t.category_id, t.currency
		ORDER BY total_amount DESC
	`

	var spending []CategorySpending
	err := repo.db.GetDB().SelectContext(ctx, &spending, query, userID, startDate, endDate)
	if err != nil {
		repo.logger.Error(
			"failed to get category spending",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf(
			"failed to get category spending for uid %s: %w",
			userID.String(),
			err,
		)
	}

	return spending, nil
}
