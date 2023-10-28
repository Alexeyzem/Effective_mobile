package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"local/EffectiveMobile/config"
)

type Controller interface {
	InitRoutes(e *echo.Group)
}

type EchoServer struct {
	e   *echo.Echo
	cfg config.Server
}

func NewServer(cfg config.Server, v1 Controller) *EchoServer {
	e := echo.New()
	e.Use(
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: cfg.Timeout,
		}),
		middleware.RequestID(),
		middleware.Logger(),
	)
	apiGroup := e.Group("/my")
	{
		v1.InitRoutes(apiGroup)
	}

	return &EchoServer{e: e, cfg: cfg}
}

func (s *EchoServer) Run() (func(), error) {
	if err := s.e.Start(fmt.Sprintf(":%v", s.cfg.Port)); err != nil {
		return func() {}, errors.Wrap(err, "echo server start")
	}
	return func() {
		s.e.Close()
	}, nil
}
