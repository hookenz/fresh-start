package server

import (
	"embed"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hookenz/moneygo/api/server/handler"
	"github.com/hookenz/moneygo/api/server/logging"
	"github.com/hookenz/moneygo/web/pages"
)

type Server struct {
	api      *echo.Echo
	address  string
	staticfs embed.FS
}

func New(address string, staticfs embed.FS) *Server {
	s := &Server{
		api:      echo.New(),
		address:  address,
		staticfs: staticfs,
	}

	s.api.HideBanner = true

	s.setupMiddleware()
	s.setupHandlers()
	s.setupStaticHandler()
	return s
}

func (s *Server) setupMiddleware() {
	logging.NewLogger()
	s.api.Use(logging.LoggingMiddleware)
	s.api.Use(middleware.Recover())
}

func (s *Server) setupHandlers() {
	s.api.POST("/api/authenticate", handler.Authenticate)
	s.api.GET("/", IndexHandler)
	s.api.GET("/login", LoginHandler)
}

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func IndexHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Index())
}

func LoginHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Login())
}

func (s *Server) setupStaticHandler() {
	// Serve the frontend at "/"
	fs := echo.MustSubFS(s.staticfs, "")
	s.api.StaticFS("/", fs)
}

func (s *Server) Start() {
	s.api.Start(s.address)
}
