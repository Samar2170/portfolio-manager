package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Samar2170/portfolio-manager/securities"
	"github.com/Samar2170/portfolio-manager/utils"
	"github.com/labstack/echo/v4"
)

func SearchStocks(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	log.Println(symbol, page, limit)
	if page == "" {
		page = "1"
	}
	if symbol == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"messsage": "symbol cant be empty",
		})
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 20
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"messsage": "page must be a number",
		})
	}
	pagination := utils.Pagination{
		Page:  pageInt,
		Limit: limitInt,
	}
	pages, stocks, err := securities.SearchStockSymbol(symbol, pagination)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Something wrong",
		})
	}
	pages.Rows = stocks
	return c.JSON(http.StatusOK, pages)
}

func SearchMutualFunds(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	log.Println(symbol, page, limit)
	if page == "" {
		page = "1"
	}
	if symbol == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"messsage": "symbol cant be empty",
		})
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 20
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"messsage": "page must be a number",
		})
	}
	pagination := utils.Pagination{
		Page:  pageInt,
		Limit: limitInt,
	}
	pages, stocks, err := securities.SearchMutualFunds(symbol, pagination)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Something wrong",
		})
	}
	pages.Rows = stocks
	return c.JSON(http.StatusOK, pages)
}

func RegisterBond(c echo.Context) error {
	name, symbol := c.QueryParam("name"), c.QueryParam("symbol")
	securityCode, exchange := c.QueryParam("security_code"), c.QueryParam("exchange")
	ipFreq, ipDate, ipRate := c.QueryParam("ip_frequency"), c.QueryParam("ip_date"), c.QueryParam("ip_rate")
	mtDate, faceValue, mtValue := c.QueryParam("maturity_date"), c.QueryParam("face_value"), c.QueryParam("maturity_value")
	if name == "" || symbol == "" || ipRate == "" || ipFreq == "" || mtDate == "" || faceValue == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "all fields should be non empty (symbol,name,ip_frequency,maturity_date,ip_date,ip_rate)",
		})
	}
	_, err := securities.CreateListedNCD(name, symbol, securityCode, exchange, ipRate, ipFreq,
		ipDate, mtDate, faceValue, mtValue)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusBadRequest, Response{
		Message: "Bond Registered Successfully",
	})
}
