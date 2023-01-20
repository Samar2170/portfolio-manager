package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartApiServer() {
	e := echo.New()
	e.POST("/signup", signup)
	e.POST("/login", login)

	e.POST("/register-account/demat", RegisterDematAccounts)
	e.POST("/register-account/bank", RegisterBankAccounts)
	e.POST("/register-trade/stock", RegisterStockTrades)
	e.POST("/register-trade/mf", RegisterMFTrades)
	e.POST("/register-fd", RegisterFD)

	e.POST("/bulk-upload/:security", UploadFile)

	e.GET("/securities/mutual-funds/search", SearchMutualFunds)
	e.GET("/securities/stocks/search", SearchStocks)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "header:Authorization",
		ContextKey:  "user",
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/signup" || c.Path() == "/login" {
				return true
			}
			return false
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))

}
