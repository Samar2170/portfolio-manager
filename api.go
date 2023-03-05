package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type Response struct {
	Message string
	Data    interface{}
}

func StartApiServer() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000/", "http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowCredentials,
			echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.POST("/signup", signup)
	e.POST("/login", login)

	e.POST("/register-account/demat", RegisterDematAccounts)
	e.POST("/register-account/bank", RegisterBankAccounts)
	e.POST("/register-trade/stock", RegisterStockTrades)
	e.POST("/register-trade/mf", RegisterMFTrades)
	e.POST("/register-trade/listed-ncd", RegisterUnistedNCDTrade)
	e.POST("/register-trade/unlisted-ncd", RegisterUnistedNCDTrade)
	e.POST("/register-fd", RegisterFD)

	e.GET("/view-accounts", ViewAccounts)

	e.POST("/bulk-upload/:security", UploadFile)
	e.POST("bulk-upload-template/:security", DownloadTemplateFile)

	e.GET("/securities/mutual-funds/search", SearchMutualFunds)
	e.GET("/securities/stocks/search", SearchStocks)
	e.GET("/securities/stocks-list", StocksList)
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
		ErrorHandler: func(err error) error {
			log.Println(err.Error())
			return err
		},
	}))

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":1323"))

}
