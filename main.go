package main

import "sync"

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go RunServicesConcurrently()
	// testNewService()
	ac := APIClient{Token: ""}
	wg.Add(1)
	go ac.getToken(1592798840)

	wg.Wait()
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
	StartBot()
}
