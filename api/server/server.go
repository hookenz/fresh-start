package server

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"github.com/hookenz/moneygo/api/db"
	"github.com/hookenz/moneygo/api/server/handler"
	"github.com/hookenz/moneygo/api/server/middleware/logging"
	"github.com/hookenz/moneygo/web/pages"
)

type Server struct {
	e        *echo.Echo
	address  string
	staticfs embed.FS
	db       db.Database
}

func New(address string, db db.Database, staticfs embed.FS) *Server {
	s := &Server{
		e:        echo.New(),
		address:  address,
		staticfs: staticfs,
		db:       db,
	}

	s.e.HideBanner = true

	s.setupMiddleware()
	s.setupHandlers()
	s.setupStaticHandler()
	return s
}

func (s *Server) setupMiddleware() {
	logging.NewLogger()
	s.e.Use(logging.LoggingMiddleware)
	s.e.Use(middleware.Recover())
	s.e.Use(session.Middleware(s.db.SessionStore()))
}

func (s *Server) setupHandlers() {
	s.e.GET("/", IndexHandler)
	s.e.GET("/login", LoginHandler)
	s.e.GET("/home", HomeHandler)

	api := s.e.Group("/api")

	h := handler.NewHandler(s.db)
	api.POST("/authenticate", h.Authenticate)

	s.e.GET("/create-session", func(c echo.Context) error {
		sess, err := session.Get("id", c)
		if err != nil {
			return err
		}

		log.Debug().Msgf("CreateSession")

		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["foo"] = "bar"
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})

	s.e.GET("/read-session", func(c echo.Context) error {
		sess, err := session.Get("id", c)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, fmt.Sprintf("foo=%v\n", sess.Values["foo"]))
	})
}

func IndexHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Index())
}

func LoginHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Login())
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Home())
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

func (s *Server) setupStaticHandler() {
	// Serve the frontend at "/"
	fs := echo.MustSubFS(s.staticfs, "")
	s.e.StaticFS("/", fs)
}

func (s *Server) Start() {
	s.e.Start(s.address)
}
