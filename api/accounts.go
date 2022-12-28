package api

import (
	"net/http"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	user, err := account.GetUserByUsername(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if user.Password != password {
		return c.String(http.StatusInternalServerError, "Invalid Password")
	}
	claims := &account.JwtCustomClaims{
		user.ID,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})

}
func signup(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")
	if username == "" || password == "" || email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "All fields are required",
		})
	}

	user := account.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	_, err2 := user.GetOrCreate()
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User already exists",
		})
	}
	return c.JSON(http.StatusOK, user)
}

func RegisterDematAccounts(c echo.Context) error {
	code := c.FormValue("code")
	broker := c.FormValue("broker")

	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User Id not found",
		})
	}
	da := account.DematAccount{
		UserId: user.Id,
		Code:   code,
		Broker: broker,
	}
	da.GetOrCreate()
	return c.JSON(
		http.StatusAccepted, map[string]string{
			"message": "Demat Account Created",
		})
}
