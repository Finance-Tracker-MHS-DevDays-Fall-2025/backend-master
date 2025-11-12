package database

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const (
	startDataBaseQuery = "start database query"
	endDataBaseQuery   = "end database query"
)

type PgLogger struct {
	Logger *zap.Logger
}

func NewPgxLogger(logger *zap.Logger) *PgLogger {
	return &PgLogger{
		Logger: logger,
	}
}

func (p *PgLogger) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {

	p.Logger.Info(startDataBaseQuery, zap.String("query", data.SQL))
	return ctx
}

func (p *PgLogger) TraceQueryEnd(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	var e *pgconn.PgError

	if data.Err != nil {
		if errors.As(data.Err, &e) && e.Code == pgerrcode.UniqueViolation {
			p.Logger.Warn(endDataBaseQuery, zap.Error(data.Err))
			return
		}
		p.Logger.Error(endDataBaseQuery, zap.Error(data.Err))
	}

	p.Logger.Info(endDataBaseQuery, zap.String("status", data.CommandTag.String()))
}
