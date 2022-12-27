package securities

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func LoadNiftyStocks() error {
	f, err := excelize.OpenFile(NSE500StocksFile)
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			handleError(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	entries := 0
	for _, row := range rows {
		entries += 1
		if entries > 1 {
			s := Stock{Name: row[0], Industry: row[1], Symbol: row[2], SecurityCode: row[4], Exchange: "NSE"}
			err := s.create()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println(entries)
	return nil
}

func LoadMutualFunds() error {
	f, err := excelize.OpenFile(MFSchemesFile)
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			handleError(err)
		}
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	entries := 0
	var activeRows [][]string
	for _, row := range rows {
		entries += 1
		fmt.Println(entries, row)
		if entries > 1 {
			if len(row) < 9 || row[8] != "" {
				activeRows = append(activeRows, row)
			}
		}
	}
	fmt.Println(len(rows))
	fmt.Println(len(activeRows))

	for _, activeRow := range activeRows {
		mf := MutualFund{AMC: activeRow[0], SchemeName: activeRow[2], SchemeNAVName: activeRow[5], Category: activeRow[4]}
		err := mf.create()
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				continue
			}
			return err
		}
	}

	return nil

}
