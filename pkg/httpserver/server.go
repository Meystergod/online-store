package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	ErrEmptyHTTPHandler  = errors.New("empty http handler")
	ErrUnknownHTTPMethod = errors.New("unknown http method")
)

type ServerDeps struct {
	Address string `envconfig:"HTTP_ADDR" default:"0.0.0.0:9100"`
}

type Server struct {
	address    string
	echoServer *echo.Echo
}

func NewServer(deps *ServerDeps) *Server {
	echoServer := echo.New()
	echoServer.Use(middleware.Recover())
	echoServer.Debug = true
	echoServer.DisableHTTP2 = true
	echoServer.HideBanner = true
	echoServer.HidePort = true

	s := &Server{
		address:    deps.Address,
		echoServer: echoServer,
	}

	s.echoServer = echoServer

	return s
}

func (s *Server) Server() *echo.Echo {
	return s.echoServer
}

func (s *Server) Start(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Str("bind_addr", s.address).Msg("listen and serve http api")
	if err := s.echoServer.Start(s.address); err != nil {
		return errors.Wrap(err, "start echo server")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("shutdown http api server")
	if err := s.echoServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shudown echo server")
	}

	return nil
}
