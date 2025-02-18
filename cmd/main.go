package main

import (
	"Go_Assignment/app"
	"Go_Assignment/loggerutil"
	"os"
)

func main() {
	// Clear App.log file content
	os.Truncate("App.log", 0)

	// Start API routing
	app.HandleRoutes()

	// Wait for all go routines to end and close channel
	loggerutil.WaitForAllRoutineToEnd()
	loggerutil.CloseChannel()
}
