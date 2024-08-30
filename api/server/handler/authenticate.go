package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	u "github.com/hookenz/moneygo/api/services/user"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//
// TODO: https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
//

type User struct {
	Username string
	Password string
}

func (h *Handler) Authenticate(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	sess, err := session.Get("id", c)
	if err != nil {
		return err
	}

	fmt.Printf("id: %v\n", sess)

	u, err := u.Authenticate(h.db, user.Username, user.Password)
	if err != nil {
		return err
	}

	log.Debug().Msgf("User authenticated %v", u.Name)

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.Redirect(200, "/")
}
