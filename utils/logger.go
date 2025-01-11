package utils

import "log"

func Info(msg string) {
	log.Println("[INFO]:", msg)
}

func Error(msg string, err error) {
	log.Printf("[ERROR]: %s - %v\n", msg, err)
}
