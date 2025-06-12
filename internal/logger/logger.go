package logger

import (
	"io"
	"log"
	"os"
)

func InitLogger() {
	// Create (or append to) log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Write logs to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Optional: set log flags for better formatting
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
