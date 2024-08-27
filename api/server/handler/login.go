package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Username string
	Password string
}

func Authenticate(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	return c.Redirect(200, "/")
}
