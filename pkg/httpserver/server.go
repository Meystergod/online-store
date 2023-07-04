package httpserver

import (
	"context"
	"net/http"

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

type MiddleWareFunc func(http.Handler) http.Handler

func (s *Server) RegisterEndpoint(method, endpoint string, handler http.Handler, m ...MiddleWareFunc) error {
	if handler == nil {
		return ErrEmptyHTTPHandler
	}

	echoHandler := echo.WrapHandler(handler)

	echoMiddleware := make([]echo.MiddlewareFunc, len(m))
	for i, v := range m {
		echoMiddleware[i] = echo.WrapMiddleware(v)
	}

	switch method {
	case echo.GET:
		s.echoServer.GET(endpoint, echoHandler, echoMiddleware...)
	case echo.POST:
		s.echoServer.POST(endpoint, echoHandler, echoMiddleware...)
	case echo.PUT:
		s.echoServer.PUT(endpoint, echoHandler, echoMiddleware...)
	case echo.DELETE:
		s.echoServer.DELETE(endpoint, echoHandler, echoMiddleware...)
	case echo.PATCH:
		s.echoServer.PATCH(endpoint, echoHandler, echoMiddleware...)
	case echo.CONNECT:
		s.echoServer.CONNECT(endpoint, echoHandler, echoMiddleware...)
	case echo.OPTIONS:
		s.echoServer.OPTIONS(endpoint, echoHandler, echoMiddleware...)
	case echo.TRACE:
		s.echoServer.TRACE(endpoint, echoHandler, echoMiddleware...)
	case echo.HEAD:
		s.echoServer.HEAD(endpoint, echoHandler, echoMiddleware...)
	default:
		return ErrUnknownHTTPMethod
	}

	return nil
}
