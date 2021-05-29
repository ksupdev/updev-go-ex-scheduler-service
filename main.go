package main

import (
	"fmt"
	"time"
)

func main() {
	ms := NewMicroservice()

	timer := 1 + time.Second
	/*
		timer := 1*time.Second // every 1 second
		timer := 60*time.Second // every 1 Minute
		timer := time.Minute * 60 // every 60 Minute
	*/
	exitScheduler := ms.Schedule(timer, func(ctx IContext) error {
		now := ctx.Now()
		ctx.Log(fmt.Sprintf("Tick at %s", now.Format("15:04:05")))
		return nil
	})

	defer func() { exitScheduler <- true }()
	// When main() has completed ,the defer func will set exitScheduler (Channel) <- true for stop timer in Schedule func
	defer ms.Cleanup()

	// Defer will run last-in-first-out order
	// ms.Cleanup() , exitScheduler <- true

	ms.Start()

}
