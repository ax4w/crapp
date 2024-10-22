package main

import (
	"context"
	"crapp/internal"
	"crapp/internal/parse"
	"crapp/internal/ui"
	"sync"
)

func main() {
	r := internal.Scan()
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func(wg *sync.WaitGroup, c context.Context) {
		parse.Run(c)
		wg.Done()
	}(&wg, ctx)
	ui.Show(r)
	cancel()
	println("cancel called")
	wg.Wait()
}
