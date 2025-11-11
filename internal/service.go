package internal

import (
	"backend-master/internal/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Service struct {
	echo   *echo.Echo
	logger *zap.Logger
}

func NewService(zap *zap.Logger) *Service {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(logger.NewLoggerMiddleware(zap))

	s := &Service{
		echo:   e,
		logger: zap,
	}

	s.registerRoutes()

	return s
}

func (s *Service) Start(addr string) error {
	s.logger.Info("starting server", zap.String("addr", addr))
	return s.echo.Start(addr)
}

func (s *Service) Shutdown() error {
	s.logger.Info("shutting down server")
	return s.echo.Close()
}

func (s *Service) registerRoutes() {
	s.echo.GET("/ping", s.handlePing)
}

func (s *Service) handlePing(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"ping": "pong",
	})
}
