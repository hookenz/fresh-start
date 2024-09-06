package handler

import (
	"net/http"

	u "github.com/hookenz/moneygo/api/services/user"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//
// TODO: https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
//

type User struct {
	Username string `query:"username" json:"username"`
	Password string `query:"password" json:"password"`
}

func (h *Handler) Authenticate(c echo.Context) error {
	var user User

	user.Username = c.FormValue("username")
	user.Password = c.FormValue("password")

	// Bind seems to be very problematic! might come back to it
	// err := c.Bind(&user)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "bad request")
	// }

	log.Debug().Msgf("user.name=%s, user.password=%s", user.Username, user.Password)

	u, err := u.Authenticate(h.db, user.Username, user.Password)
	if err != nil {
		return err
	}

	log.Debug().Msgf("User authenticated %v", u.Name)

	// Generate an id
	id, err := h.db.CreateSession()
	if err != nil {
		return err
	}

	// Create a session cookie
	writeSessionCookie(c, id)
	return c.Redirect(302, "/home")
}

func (h *Handler) Logout(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	return c.Redirect(302, "/")
}

func writeSessionCookie(c echo.Context, sessionid string) {
	log.Debug().Msg("Set Session Cookie")
	cookie := new(http.Cookie)
	cookie.Name = "id"
	cookie.Value = sessionid
	cookie.Path = "/"
	cookie.MaxAge = 24 * 60 * 60
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}
