// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import "time"

// IContext is the context for service
// Interface of Context for Scheduler service
type IContext interface {
	Log(message string)
	Param(name string) string
	Response(responseCode int, responseData interface{})
	ReadInput() string
	ReadInputs() []string

	Now() time.Time
}
