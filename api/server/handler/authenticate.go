package handler

import (
	"fmt"
	"net/http"

	"github.com/hookenz/moneygo/api/db"
	u "github.com/hookenz/moneygo/api/services/user"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

//
// TODO: https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
//

type User struct {
	Username string
	Password string
}

func Authenticate(c echo.Context, db db.Database) error {
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

	u.Authenticate(db, user.Username, user.Password)

	return c.Redirect(200, "/")
}
