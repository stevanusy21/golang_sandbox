package utils

import "log"

func LogError(location string, logMessage string, errorValue error) {
	log.Printf("[%s] %s: %v \n", location, logMessage, errorValue)
}

func LogErrorNoValue(location string, logMessage string) {
	log.Printf("[%s] %s \n", location, logMessage)
}

func LogFatal(location string, logMessage string, errorValue error) {
	log.Fatalf("[%s] %s: %v \n", location, logMessage, errorValue)
}

func LogFatalNoValue(location string, logMessage string) {
	log.Fatalf("[%s] %s \n", location, logMessage)
}

func LogInfo(location string, logMessage string) {
	log.Printf("[%s] %s \n", location, logMessage)
}

func LogWarn(location string, logMessage string) {
	log.Printf("[%s] %s \n", location, logMessage)
}
