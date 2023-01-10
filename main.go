package main

import (
	"sync"
)

func main() {
	RunServicesConcurrently()
}

func RunServicesConcurrently() {
	var wg sync.WaitGroup

	wg.Add(1)
	go StartApiServer()

	wg.Add(1)
	go RunSuperVisor()

	wg.Wait()
}

func testNewService() {

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
