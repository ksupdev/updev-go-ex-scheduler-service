package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// -- Defind interface and function type for Scheduler service --

// IMicroservice is interface for centralized service management
type IMicroservice interface {
	Start() error
	Stop()
	Cleanup() error
	Log(message string)

	// Scheduler Services
	Schedule(timer time.Duration, h ServiceHandleFunc) chan bool /*exit channel*/
}

// ServiceHandleFunc is Function types
// IContext use fore manage Context (read and response)
type ServiceHandleFunc func(ctx IContext) error

// Microservice is the centralized service management
type Microservice struct {
	// exitChannel is way for exit from Microservie. (if exitChannel = true then exit from this service )
	exitChannel chan bool
}

func NewMicroservice() *Microservice {
	return &Microservice{}
}

// Start start all registered services
func (ms *Microservice) Start() error {
	// There are 2 ways to exit from  Microservices
	// 1. The SigTerm can be send from outside program such as from K8S
	// 2. Send true to ms.exitChannel

	// Make Chan and specific size = 1
	// the osQuit is channel for handle incoming signals to stop this microservice
	osQuit := make(chan os.Signal, 1)
	ms.exitChannel = make(chan bool, 1)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)
	exit := false

	// Infinity loop
	for {
		if exit {
			break
		}

		select {
		case <-osQuit:
			// Handle The Signal for request close service from out side program
			exit = true
		case <-ms.exitChannel:
			// Handle Close service inside program
			exit = true
		}

	}

	return nil
}

// Stop stop the service
func (ms *Microservice) Stop() {
	if ms.exitChannel == nil {
		// Force exit func
		return
	}

	ms.exitChannel <- true
}

// Cleanup clean resources up from every registered services before exit
func (ms *Microservice) Cleanup() error {
	return nil
}

func (ms *Microservice) Log(tag string, message string) {
	fmt.Println(tag+": ", message)
}

func (ms *Microservice) Schedule(timer time.Duration, h ServiceHandleFunc) chan bool {
	// exitChan must be call exitChan <- true from caller to exit scheduler
	exitChan := make(chan bool, 1)

	go func() {
		t := time.NewTicker(timer)
		done := make(chan bool, 1)
		isExit := false
		isExitMutex := sync.Mutex{}

		// This routine handle when process has done
		go func() {
			fmt.Println("--- wait.. exitChan")
			<-exitChan // Block wait..
			fmt.Printf("--- wait.. exitChan %v \n", exitChan)
			isExitMutex.Lock()
			isExit = true
			isExitMutex.Unlock()
			// Stop Tick() and send done message to exit for loop below
			// Ref: From the documentation http://golang.org/pkg/time/#Ticker.Stop
			// Stop turns off a ticker. After stop, no more ticks will be sent.
			// Stop does not close the Channel, to prevent a read from the channel succeeding incorrectly
			t.Stop()
			done <- true
		}()

		for {
			select {
			// t.C is Channel of timer, Timer will send value every duration
			case execTime := <-t.C:
				// Execute schedule
				isExitMutex.Lock()
				if isExit {
					isExitMutex.Unlock()
					// Done in the next round
					continue
				}

				isExitMutex.Unlock()

				now := time.Now()
				// The schedule that older than 10s, will be skip, because t.C is buffer at size 1
				diff := now.Sub(execTime).Seconds()
				if diff > 10 {
					continue
				}
				h(NewSchedulerContext(ms))
			case <-done:
				// This schedule has finished
				return
			}
		}
	}()
	return exitChan
}
