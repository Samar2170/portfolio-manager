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
	DtFormat       = "2006-01-02"
	FDTemplateFile = "Assets/templates/FDBUTemp.csv"
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

func RegisterMFTrades(c echo.Context) error {
	var err error
	mfId := c.FormValue("mutual_fund_id")
	// mfIdInt, err := strconv.ParseInt(mfId, 10, 64)
	// if err != nil {
	// return c.JSON(http.StatusBadRequest, map[string]string{
	// "message": "mutual_fund_id should be a number",
	// })
	// }
	dematAccCode := c.FormValue("demat")
	quantity := c.FormValue("quantity")
	price := c.FormValue("price")
	tradeType := c.FormValue("trade_type")
	tradeDate := c.FormValue("trade_date")
	if mfId == "" || dematAccCode == "" || tradeType == "" || quantity == "" || price == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "all fields should be non empty (mutual_fund_id,demat,quantity,price,trade_type)",
		})
	}
	_, err = portfolio.CreateMFTrade(mfId, dematAccCode, quantity, price, tradeType, tradeDate)
	if err != nil {
		return c.JSON(
			http.StatusConflict, Response{
				Message: err.Error(),
			},
		)
	}
	return c.JSON(
		http.StatusConflict, Response{
			Message: "trade registered successfully",
		})

	//  kept for contigency purposes
	// var tradeDateParsed time.Time
	// if tradeDate == "" {
	// 	tradeDateParsed = time.Now()
	// } else {
	// 	tradeDateParsed, err = time.Parse(DtFormat, tradeDate)
	// 	if err != nil {
	// 		return c.JSON(http.StatusBadRequest, map[string]string{
	// 			"message": "Trade date should be in format 2022-11-22",
	// 		})
	// 	}
	// }

	// quantityFloat, err := strconv.ParseFloat(quantity, 64)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"message": "quantity should be a number",
	// 	})
	// }
	// priceFloat, err := strconv.Atoi(price)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"message": "price should be a number",
	// 	})
	// }

	// mfTrade, err := portfolio.NewMFTrade(uint(mfIdInt), tradeType, dematAccCode, quantityFloat, float64(priceFloat), tradeDateParsed)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"message": err.Error(),
	// 	})
	// }
	// err = portfolio.RegisterMFTrade(*mfTrade)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"message": err.Error(),
	// 	})
	// }
	// return c.JSON(http.StatusAccepted, map[string]string{
	// 	"message": "trade registered successfully",
	// })
}

func RegisterFD(c echo.Context) error {
	bankName := c.FormValue("bank")
	amount := c.FormValue("amount")
	mtAmount := c.FormValue("maturity_amount")

	if mtAmount == "" {
		mtAmount = amount
	}
	ipRate := c.FormValue("ip_rate")
	ipFreq := c.FormValue("ip_frequency")
	startDate := c.FormValue("start_date")
	ipDate := c.FormValue("ip_date")
	if ipDate == "" {
		ipDate = startDate
	}
	mtDate := c.FormValue("maturity_date")
	accNumber := c.FormValue("account_number")

	if ipDate == "" {
		ipDate = startDate
	}
	if amount == "" || ipRate == "" || ipFreq == "" || startDate == "" || mtDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "amount , ip_rate, ip_frequency, start_date, account_number cant be empty",
		})
	}
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "make sure amount is a number",
		})
	}

	mtAmountFloat, err := strconv.ParseFloat(mtAmount, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "make sure Mtamount is a number",
		})
	}

	ipRateFloat, err := strconv.ParseFloat(ipRate, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "make sure ipRate is in 0.0x Format",
		})
	}
	fdh, err := portfolio.CreateFDHolding(bankName, amountFloat, mtAmountFloat, ipRateFloat, ipFreq, startDate, ipDate, mtDate, accNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"message": fmt.Sprintf("FD created successfully with ID %d", fdh.ID),
	})

}

func DownloadFDBulkUploadFile(c echo.Context) error {
	return c.File(FDTemplateFile)
}
