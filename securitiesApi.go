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
