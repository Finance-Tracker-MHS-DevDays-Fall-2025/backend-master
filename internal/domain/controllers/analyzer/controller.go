package analyzer

import (
	"context"
	"fmt"
	"time"

	pb "backend-master/internal/api-gen/proto/analyzer"
	"backend-master/internal/api-gen/proto/common"
	"backend-master/internal/data/repositories/analyzer"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AnalyzerController interface {
	GetStatistics(
		ctx context.Context,
		userID string,
		startDate time.Time,
		endDate time.Time,
		groupBy common.TimePeriod,
	) (*pb.GetStatisticsResponse, error)

	GetForecast(
		ctx context.Context,
		userID string,
		period common.TimePeriod,
		periodsAhead int32,
	) (*pb.GetForecastResponse, error)
}

type analyzerControllerImpl struct {
	repo   analyzer.AnalyzerRepository
	client *analyzer.AnalyzerClient
	logger *zap.Logger
}

func NewController(
	repo analyzer.AnalyzerRepository,
	client *analyzer.AnalyzerClient,
	logger *zap.Logger,
) AnalyzerController {
	return &analyzerControllerImpl{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

func (cont *analyzerControllerImpl) GetStatistics(
	ctx context.Context,
	userID string,
	startDate time.Time,
	endDate time.Time,
	groupBy common.TimePeriod,
) (*pb.GetStatisticsResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		cont.logger.Error(
			"invalid user ID",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	totalIncome, totalExpense, err := cont.repo.GetStatistics(ctx, uid, startDate, endDate)
	if err != nil {
		cont.logger.Error(
			"failed to get statistics from repository",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	groupByStr := mapTimePeriodToString(groupBy)
	periodBalances, err := cont.repo.GetPeriodBalances(ctx, uid, startDate, endDate, groupByStr)
	if err != nil {
		cont.logger.Error(
			"failed to get period balances from repository",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get period balances: %w", err)
	}

	categorySpending, err := cont.repo.GetCategorySpending(ctx, uid, startDate, endDate)
	if err != nil {
		cont.logger.Error(
			"failed to get category spending from repository",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get category spending: %w", err)
	}

	pbCategorySpending := make([]*pb.CategorySpending, 0, len(categorySpending))
	for _, cs := range categorySpending {
		pbCategorySpending = append(pbCategorySpending, cs.ToProto())
	}

	pbPeriodData := make([]*pb.PeriodBalance, 0, len(periodBalances))
	for _, period := range periodBalances {
		pbPeriodData = append(pbPeriodData, period.ToProto(pbCategorySpending))
	}

	return &pb.GetStatisticsResponse{
		TotalIncome: &common.Money{
			Amount:   totalIncome,
			Currency: "RUB",
		},
		TotalExpense: &common.Money{
			Amount:   totalExpense,
			Currency: "RUB",
		},
		PeriodData: pbPeriodData,
	}, nil
}

func (cont *analyzerControllerImpl) GetForecast(
	ctx context.Context,
	userID string,
	period common.TimePeriod,
	periodsAhead int32,
) (*pb.GetForecastResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		cont.logger.Error(
			"invalid user ID",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	totalIncome, totalExpense, err := cont.repo.GetStatistics(ctx, uid, startDate, endDate)
	if err != nil {
		cont.logger.Error(
			"failed to get historical statistics for forecast",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get historical statistics: %w", err)
	}

	avgIncome := totalIncome / 12
	avgExpense := totalExpense / 12

	forecasts := make([]*pb.Forecast, 0, periodsAhead)
	for i := int32(0); i < periodsAhead; i++ {
		forecastStart := addPeriod(endDate, period, i)
		forecastEnd := addPeriod(endDate, period, i+1)

		forecast := analyzer.Forecast{
			PeriodStart:     forecastStart,
			PeriodEnd:       forecastEnd,
			ExpectedIncome:  avgIncome,
			ExpectedExpense: avgExpense,
			ExpectedBalance: avgIncome - avgExpense,
			Currency:        "RUB",
		}

		forecasts = append(forecasts, forecast.ToProto(nil))
	}

	return &pb.GetForecastResponse{
		Forecasts: forecasts,
	}, nil
}

func mapTimePeriodToString(period common.TimePeriod) string {
	switch period {
	case common.TimePeriod_TIME_PERIOD_QUARTER:
		return "quarter"
	case common.TimePeriod_TIME_PERIOD_MONTH:
		return "month"
	case common.TimePeriod_TIME_PERIOD_YEAR:
		return "year"
	default:
		return "month"
	}
}

func addPeriod(t time.Time, period common.TimePeriod, count int32) time.Time {
	switch period {
	case common.TimePeriod_TIME_PERIOD_QUARTER:
		return t.AddDate(0, int(count)*4, 0)
	case common.TimePeriod_TIME_PERIOD_MONTH:
		return t.AddDate(0, int(count), 0)
	case common.TimePeriod_TIME_PERIOD_YEAR:
		return t.AddDate(int(count), 0, 0)
	default:
		return t.AddDate(0, int(count), 0)
	}
}
