package main

import (
	"sync"
)

func main() {
	// testNewService()
	RunServicesConcurrently()
}

func RunServicesConcurrently() {
	var wg sync.WaitGroup

	wg.Add(1)
	go StartApiServer()

	wg.Add(1)
	go RunSuperVisor()

	wg.Add(1)
	go testNewService()

	wg.Wait()
}

func testNewService() {
	// securities.UpdateNextIPDatesFDs()
	// fd, _ := securities.GetFDByID(2)
	// nipd := fd.CalculateNextIPDate()
	// t1 := fd.IPDate.AddDate(0, 1, 0)
	// fmt.Println(fd.IPDate, fd.IPFreq, t1)
	// t1 := time.Now()
	// fmt.Println(t1)
	// tom := t1.AddDate(0, 0, 1)
	// mon := t1.AddDate(0, 1, 0)
	// fmt.Printf("tommorow is %s, month after %s", tom, mon)

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
