package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	OTP_MIN = 100000
	OTP_MAX = 999999
)

func chatLogin(c echo.Context) error {
	chatId := c.FormValue("chat_id")
	chatIdInt, _ := strconv.ParseInt(chatId, 10, 64)

	user, _ := account.GetUserByChatId(int64(chatIdInt))
	claims := &account.JwtCustomClaims{
		user.ID, user.Username, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
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
	user, err2 := user.GetOrCreate()
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User already exists",
		})
	}
	createAccountStatusJob := CreateUserAccountStatusJob{UserId: user.ID}
	createGeneralAccountJob := CreateGeneralAccountJob{UserId: user.ID}
	JobQueue <- createAccountStatusJob
	JobQueue <- createGeneralAccountJob

	return c.JSON(http.StatusOK, user)
}

func RegisterDematAccounts(c echo.Context) error {
	code := c.FormValue("code")
	broker := c.FormValue("broker")

	if code == "" || broker == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Code and Broker must not be empty",
		})
	}
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
	da.Create()
	return c.JSON(
		http.StatusAccepted, map[string]string{
			"message": "Demat Account Created",
		})
}
func RegisterBankAccounts(c echo.Context) error {
	accountNo := c.FormValue("account_number")
	bank := c.FormValue("bank")
	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User Id not found",
		})
	}
	if accountNo == "" || bank == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bank and account_number must not be empty",
		})
	}
	ba := account.BankAccount{
		AccountNo: accountNo,
		Bank:      bank,
		UserId:    user.Id,
	}
	ba.Create()
	return c.JSON(
		http.StatusAccepted, map[string]string{
			"message": "Bank Account Created",
		})
}

func ViewAccounts(c echo.Context) error {
	accountType := c.FormValue("type")
	// if accountType != "bank" && accountType != "demat" {
	// 	return c.JSON(http.StatusBadRequest, Response{
	// 		Message: "type parameter should bank or demat",
	// 	})
	// }

	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusForbidden, Response{Message: "Please log in"})
	}
	dataResponse := make(map[string]interface{})

	switch accountType {
	case "bank":
		bankAccounts, err := account.GetBankAccountsByUser(user.Id)
		if err != nil {
			return c.JSON(http.StatusForbidden, Response{Message: err.Error()})
		}
		dataResponse["bank_accounts"] = bankAccounts
	case "demat":
		dematAccounts, err := account.GetDematAccountsByUser(user.Id)
		if err != nil {
			return c.JSON(http.StatusForbidden, Response{Message: err.Error()})
		}
		dataResponse["demat_accounts"] = dematAccounts
	default:
		dematAccounts, err := account.GetDematAccountsByUser(user.Id)
		if err != nil {
			return c.JSON(http.StatusForbidden, Response{Message: err.Error()})
		}
		bankAccounts, err := account.GetBankAccountsByUser(user.Id)
		if err != nil {
			return c.JSON(http.StatusForbidden, Response{Message: err.Error()})
		}
		dataResponse["bank_accounts"] = bankAccounts
		dataResponse["demat_accounts"] = dematAccounts
	}

	return c.JSON(http.StatusOK, Response{
		Message: "Good",
		Data:    dataResponse,
	})
}

func GetTelegramOTP(c echo.Context) error {
	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User Id not found",
		})
	}
	randomOTP := rand.Intn(OTP_MAX-OTP_MIN) + OTP_MIN
	telegramOtp := account.TelegramOTP{
		UserId:    user.Id,
		OTP:       uint(randomOTP),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = telegramOtp.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong")
	}
	responseMap := Response{
		Data: map[string]int{
			"OTP": randomOTP,
		},
		Message: "valid for 1 min 30 seconds",
	}
	return c.JSON(
		http.StatusOK, responseMap,
	)
}
