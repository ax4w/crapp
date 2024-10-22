package main

import (
	"crapp/internal"
	"crapp/internal/front"
	"crapp/internal/middle"
	"sync"
)

func main() {
	r := internal.Scan()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		front.Run()
		wg.Done()
	}(&wg)
	internal.Run(r)
	middle.BackToFront <- middle.Packet{Id: middle.CMD_DONE}
	wg.Wait()
}
