package main

import (
	"fmt"
	"sync"

	"github.com/Samar2170/portfolio-manager/portfolio"
)

func main() {
	RunServicesConcurrently()
	// testNewService()
}

func RunServicesConcurrently() {
	var wg sync.WaitGroup

	wg.Add(1)
	go StartApiServer()

	wg.Add(1)
	go RunSuperVisor()

	wg.Add(1)
	go StartCronServer()

	wg.Wait()

}

func testNewService() {
	err := portfolio.ParseFDFile(2)
	if err != nil {
		fmt.Println(err)
	}
}

// func loadScripts() {
// 	err := securities.LoadNiftyStocks()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	stocks, err := securities.GetAllStocks()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	for _, stock := range stocks {
// 		fmt.Println(stock.ID)
// 	}
// 	err = securities.LoadMutualFunds()
// 	if err != nil {
// 		log.Println(err)
// 	}

// }
