package utils

import "log"

func LogError(module string, logMessage string, errorValue error) {
	log.Printf("[%s] %s: %v \n", module, logMessage, errorValue)
}

func LogErrorNoValue(module string, logMessage string) {
	log.Printf("[%s] %s \n", module, logMessage)
}

func LogFatal(module string, logMessage string, errorValue error) {
	log.Fatalf("[%s] %s: %v \n", module, logMessage, errorValue)
}

func LogFatalNoValue(module string, logMessage string) {
	log.Fatalf("[%s] %s \n", module, logMessage)
}

func LogInfo(module string, logMessage string) {
	log.Printf("[%s] %s \n", module, logMessage)
}

func LogWarn(module string, logMessage string) {
	log.Printf("[%s] %s \n", module, logMessage)
}
