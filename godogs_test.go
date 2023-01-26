package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/cucumber/godog"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type godogsCtxKey struct{}

type apiFeature struct {
	resp *httptest.ResponseRecorder
}

func iLogin(ctx context.Context) error {
	fmt.Println(ctx.Value(godogsCtxKey{}))
	return nil
}

func iSignup(ctx context.Context) (context.Context, error) {
	mockUser := account.User{Username: "test", Password: "tests@123", Email: "test@email.com"}
	f := make(url.Values)
	f.Set("username", mockUser.Username)
	f.Set("password", mockUser.Password)
	f.Set("email", mockUser.Email)
	req, err := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(f.Encode()))
	if err != nil {
		return ctx, err
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	e := echo.New()
	_ = e.NewContext(req, httptest.NewRecorder())
	return context.WithValue(ctx, godogsCtxKey{}, mockUser), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// api := apiFeature{}

	// ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	// 	api.newServer(sc)
	// 	return ctx, nil
	// })
	ctx.Step(`^i login$`, iLogin)
	ctx.Step(`^i signup$`, iSignup)
}

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
