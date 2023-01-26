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
	"github.com/stretchr/testify/suite"
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
type TestSuite struct {
	suite.Suite
	Token string
}

func (ts TestSuite) TestLogin() error {
	e := echo.New()
	f := make(url.Values)
	f.Set("username", "Jon Snow")
	f.Set("password", "test@123")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var resp LoginResponse
	if assert.NoError(ts.Suite.T(), login(c)) {
		assert.Equal(ts.Suite.T(), http.StatusOK, rec.Code)

		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		if err != nil {
			return err
		}
		ts.Token = resp.Token
		assert.Greater(ts.Suite.T(), len(ts.Token), 0)
	}
	return nil
}

func (ts TestSuite) TestHelloWorld() error {
	token := ts.Token
	fmt.Println(token)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(ts.Suite.T(), hello(c)) {
		assert.Equal(ts.Suite.T(), http.StatusOK, rec.Code)
		fmt.Println(rec.Body)
	}

	return nil
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
