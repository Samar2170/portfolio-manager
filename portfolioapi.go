package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/portfolio"
	"github.com/labstack/echo/v4"
)

const (
	DtFormat = "2022-11-22"
)

func RegisterStockTrades(c echo.Context) error {
	var err error
	stockSymbol := c.FormValue("symbol")
	stockSymbol = strings.ToUpper(stockSymbol)
	dematAccCode := c.FormValue("demat")
	quantity := c.FormValue("quantity")
	price := c.FormValue("price")
	tradeType := c.FormValue("trade_type")
	tradeDate := c.FormValue("trade_date")
	if stockSymbol == "" || dematAccCode == "" || tradeType == "" || quantity == "" || price == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "all fields should be non empty (symbol,demat,quantity,price,trade_type)",
		})
	}
	var tradeDateParsed time.Time
	if tradeDate == "" {
		tradeDateParsed = time.Now()
	} else {
		tradeDateParsed, err = time.Parse(DtFormat, tradeDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Trade date should be in format 2022-11-22",
			})
		}
	}

	fmt.Println(tradeDate)
	quantityInt, err := strconv.ParseInt(quantity, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "quantity should be a number",
		})
	}
	priceFloat, err := strconv.Atoi(price)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "price should be a number",
		})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "trade_date should be in format 2022-11-22",
		})
	}

	stockTrade, err := portfolio.NewStockTrade(stockSymbol, tradeType, dematAccCode, uint(quantityInt), float64(priceFloat), tradeDateParsed)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	err = portfolio.RegisterTrade(*stockTrade)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "trade registered successfully",
	})
}
