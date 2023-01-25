package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/portfolio"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string
	Data    interface{}
}

const (
	UPLOADFILES_DIR  = "upload_files/"
	FD_TEMPLATE_FILE = "FDBUTemp.csv"
	MF_TEMPLATE_FILE = "MFTradeBUTemp.csv"
)

func DownloadTemplateFile(c echo.Context) error {
	security := c.Param("security")
	if _, ok := portfolio.ValidSecurities[security]; !ok {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Sprintf("Bad Security Parameter, Valid choices are %s", portfolio.ValidSecurities.Keys()),
		})
	}
	switch security {
	case "fd":
		return c.File(UPLOADFILES_DIR + FDTemplateFile)
	case "mf":
		return c.File(UPLOADFILES_DIR + MF_TEMPLATE_FILE)

	}
	return nil
}

// Save the file -> trigger task to parse -> save entries
//
//	-> take wrong entries and email it

func UploadFile(c echo.Context) error {
	security := c.Param("security")
	if _, ok := portfolio.ValidSecurities[security]; !ok {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Sprintf("Bad Security Parameter, Valid choices are %s", portfolio.ValidSecurities.Keys()),
		})
	}
	file, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "File not found, file parameter missing",
		})
	}
	user, err := account.UnwrapToken(c.Get("user").(*jwt.Token))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "User not found, Are you logged in",
		})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "Error while opening file",
		})
	}
	// fileSizeinMB := file.Size / (1024 * 1024)
	ogFileName, fileExtension := strings.Split(file.Filename, ".")[0], file.Filename[strings.LastIndex(file.Filename, ".")+1:]
	fileName := fmt.Sprintf("%s_%s.%s", user.Username, ogFileName, fileExtension)

	if fileExtension != "csv" {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "Only csv supported",
		})
	}
	defer src.Close()
	dst, err := os.Create(UPLOADFILES_DIR + fileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "Error while opening dst file",
		})
	}
	defer dst.Close()

	io.Copy(dst, src)
	switch security {
	case "fd":
		fileData := portfolio.FDFile{
			UserId:   user.Id,
			FileName: fileName,
			FilePath: UPLOADFILES_DIR + fileName,
			Parsed:   false,
		}
		fileData.Create()
	case "mf":
		fileData := portfolio.MFFile{
			UserId:   user.Id,
			FileName: fileName,
			FilePath: UPLOADFILES_DIR + fileName,
			Parsed:   false,
		}
		fileData.Create()
	case "stock", "stocks", "shares":
		fileData := portfolio.StockFile{
			UserId:   user.Id,
			FileName: fileName,
			FilePath: UPLOADFILES_DIR + fileName,
			Parsed:   false,
		}
		fileData.Create()
	}
	return c.JSON(http.StatusAccepted, Response{
		Message: "File Successfully Uploaded",
	})
}
