package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// type handler struct {
// 	db map[string]*account.User
// }

// var (
// 	mockDB = map[string]*account.User{
// 		"jon@labstack.com": &account.User{Username: "Jon Snow", Password: "test@123", Email: "jon@labstack.com"},
// 	}
// 	userJSON = `{"name":"Jon Snow","email":"jon@labstack.com","password":"test@123"}`
// )

type LoginResponse struct {
	Token string
}

func TestLogin(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("username", "Jon Snow")
	f.Set("password", "test@123")
	// f.Set("email", "test@gmail.com")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var resp LoginResponse
	if assert.NoError(t, login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		fmt.Println(resp)
	}

}
