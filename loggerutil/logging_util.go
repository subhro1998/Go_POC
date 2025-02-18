package loggerutil

import (
	"fmt"
	"os"
	"sync"
)

type LogMessage struct {
	Level   string
	Message string
}

var (
	logChannel   chan LogMessage
	logWaitGroup sync.WaitGroup
)

func init() { // Initialize the channel with a buffer size of 100
	logChannel = make(chan LogMessage, 100)
}

func WaitForAllRoutineToEnd() {
	logWaitGroup.Wait()
}

func CloseChannel() {
	close(logChannel)
}

func LogProcessor() {
	defer logWaitGroup.Done()
	logWaitGroup.Add(1)

	file, err := os.OpenFile("App.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write Log content to file
	for logMessage := range logChannel {
		message := "Log:: " + logMessage.Level + " ::" + logMessage.Message + "\n"
		fmt.Println("Retriving Log messages and writing in file - ", message)
		_, fileWriteErr := file.WriteString(message)
		if fileWriteErr != nil {
			panic(fileWriteErr)
		}
	}
}

func PostLogMessages(logMessages []LogMessage) {
	defer logWaitGroup.Done()
	logWaitGroup.Add(1)
	for _, logMsg := range logMessages {
		fmt.Println("Posting Log messages - ", logMsg)
		logChannel <- logMsg
	}
}
