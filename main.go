package main

import (
	"log"

	"github.com/Samar2170/portfolio-manager/securities"
)

func main() {
	// err := securities.LoadNiftyStocks()
	// if err != nil {
	// 	log.Println(err)
	// }
	// stocks, err := securities.GetAllStocks()
	// if err != nil {
	// 	log.Println(err)
	// }
	// for _, stock := range stocks {
	// 	fmt.Println(stock.ID)
	// }
	err := securities.LoadMutualFunds()
	if err != nil {
		log.Println(err)
	}
}
