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
	"google.golang.org/protobuf/types/known/timestamppb"
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
	client *analyzer.AnalyzerClient
	logger *zap.Logger
}

func NewController(
	client *analyzer.AnalyzerClient,
	logger *zap.Logger,
) AnalyzerController {
	return &analyzerControllerImpl{
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
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	resp, err := cont.client.GetStatistics(
		ctx,
		&pb.GetStatisticsRequest{
			UserId:    uid.String(),
			StartDate: timestamppb.New(startDate),
			EndDate:   timestamppb.New(endDate),
			GroupBy:   groupBy,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics from analyzer: %w", err)
	}
	return resp, nil
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

	resp, err := cont.client.GetForecast(
		ctx,
		&pb.GetForecastRequest{
			UserId:       uid.String(),
			Period:       period,
			PeriodsAhead: periodsAhead,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast from analyzer: %w", err)
	}

	return resp, nil
}
