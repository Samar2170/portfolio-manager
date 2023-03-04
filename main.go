package main

import "sync"

func main() {
	RunServicesConcurrently()
}

func RunServicesConcurrently() {
	var wg sync.WaitGroup

	wg.Add(1)
	go StartApiServer()

	wg.Add(1)
	go RunSuperVisor()

	wg.Add(1)
	go StartCronServer()

	wg.Add(1)
	go StartBot()

	wg.Wait()

}

func testNewService() {
	StartBot()
}
