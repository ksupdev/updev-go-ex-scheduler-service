package main

import (
	"fmt"
	"time"
)

// Implementation of IContext in context_scheduler.go
type SchedulerContext struct {
	// Ref to Microservice Struct in microservice.go
	ms *Microservice
}

// NewSchedulerContext is the function constructor function for create instance SchedulerContext
func NewSchedulerContext(ms *Microservice) *SchedulerContext {
	return &SchedulerContext{ms: ms}
}

// New return time.Now
func (ctx *SchedulerContext) Now() time.Time {
	return time.Now()
}

// Log will log a message
func (ctx *SchedulerContext) Log(message string) {
	fmt.Println("Scheduler: ", message)
}

// Param return parameter by name (empty in scheduler)
func (ctx *SchedulerContext) Param(name string) string {
	return ""
}

// ReadInput return message (return empty in Scheduler)
func (ctx *SchedulerContext) ReadInput() string {
	return ""
}

// ReadInputs return message in batch (return nil in scheduler)
func (ctx *SchedulerContext) ReadInputs() []string {
	return nil
}

// Response return response to client
func (ctx *SchedulerContext) Response(responseCode int, responseData interface{}) {
	return
}
