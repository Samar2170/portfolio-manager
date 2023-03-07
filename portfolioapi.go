package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Samar2170/portfolio-manager/portfolio"
	"github.com/Samar2170/portfolio-manager/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	DtFormat = "2006-01-02"
)

func RegisterStockTrades(c echo.Context) error {
	var err error
	st := new(StockTrade)
	if err := c.Bind(st); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	stockSymbol, dematAccCode := st.Symbol, st.Demat
	quantity, price := st.Quantity, st.Price
	tradeType, tradeDate := st.TradeType, st.TradeDate
	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "User detail not Found",
		})
	}
	if stockSymbol == "" || dematAccCode == "" || tradeType == "" || quantity == "" || price == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "all fields should be non empty (symbol,demat,quantity,price,trade_type)",
		})
	}
	_, err = portfolio.CreateStockTrade(stockSymbol, dematAccCode, quantity, price, tradeType, tradeDate, user.Id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest, Response{
				Message: err.Error(),
			},
		)
	}
	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "trade registered successfully",
	})
}

func RegisterMFTrades(c echo.Context) error {
	var err error
	mft := new(MFTrade)
	if err := c.Bind(mft); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	mfId, dematAccCode := mft.MutualFundId, mft.Demat
	quantity, price := mft.Quantity, mft.Price
	tradeType, tradeDate := mft.TradeType, mft.TradeDate
	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "User detail not Found",
		})
	}
	if mfId == "" || dematAccCode == "" || tradeType == "" || quantity == "" || price == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "all fields should be non empty (mutual_fund_id,demat,quantity,price,trade_type)",
		})
	}
	_, err = portfolio.CreateMFTrade(mfId, dematAccCode, quantity, price, tradeType, tradeDate, user.Id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest, Response{
				Message: err.Error(),
			},
		)
	}
	return c.JSON(
		http.StatusOK, Response{
			Message: "trade registered successfully",
		})

}

func RegisterListedNCDTrade(c echo.Context) error {
	lncdt := new(ListedNCDTrade)
	if err := c.Bind(lncdt); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	symbol, quantity, price := lncdt.Symbol, lncdt.Quantity, lncdt.Price
	dematAccount, tradeDate := lncdt.Demat, lncdt.TradeDate

	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "User detail not Found",
		})
	}
	_, err = portfolio.CreateListedNCDHolding(symbol, quantity, price, tradeDate, dematAccount, user.Id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest, Response{
				Message: err.Error(),
			},
		)
	}
	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "trade registered successfully",
	})
}

func RegisterUnistedNCDTrade(c echo.Context) error {
	symbol := strings.ToUpper(c.FormValue("symbol"))
	quantity := c.FormValue("quantity")
	price := c.FormValue("price")
	tradeDate := c.FormValue("trade_date")
	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "User detail not Found",
		})
	}
	_, err = portfolio.CreateUnlistedNCDHolding(symbol, quantity, price, tradeDate, user.Id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest, Response{
				Message: err.Error(),
			},
		)
	}
	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "trade registered successfully",
	})
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
	user, err := utils.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "User detail not Found",
		})
	}
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
	fdh, err := portfolio.CreateFDHolding(bankName, amountFloat, mtAmountFloat, ipRateFloat, ipFreq, startDate, ipDate, mtDate, accNumber, user.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"message": fmt.Sprintf("FD created successfully with ID %d", fdh.ID),
	})

}
