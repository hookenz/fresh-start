package server

import (
	"embed"

	"github.com/hookenz/moneygo/api/server/handler"
	"github.com/hookenz/moneygo/api/server/logging"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	mux      *echo.Echo
	address  string
	staticfs embed.FS
}

func New(address string, staticfs embed.FS) *Server {
	s := &Server{
		mux:      echo.New(),
		address:  address,
		staticfs: staticfs,
	}

	s.mux.HideBanner = true

	s.setupMiddleware()
	s.setupHandlers()
	s.setupStaticHandler()

	return s
}

func (s *Server) setupMiddleware() {
	logging.NewLogger()
	s.mux.Use(logging.LoggingMiddleware)
	s.mux.Use(middleware.Recover())
}

func (s *Server) setupHandlers() {
	s.mux.POST("/api/login", handler.Login)
}

func (s *Server) setupStaticHandler() {
	// Serve the frontend at "/"
	fs := echo.MustSubFS(s.staticfs, "dist/public")
	s.mux.StaticFS("/", fs)
}

func (s *Server) Start() {
	s.mux.Start(s.address)
}
