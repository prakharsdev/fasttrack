package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"fasttrack/config"
	"fasttrack/internal/consumer"
	"fasttrack/internal/db"
	"fasttrack/internal/logger"
	"fasttrack/internal/publisher"

	"github.com/streadway/amqp"
)

func main() {
	// Initialize centralized logger
	logger.InitLogger()
	log.Println("Starting FastTrack Data Ingestion App")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MySQL database
	db.InitDB(cfg)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQUri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	q, err := ch.QueueDeclare("payments", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	log.Println("RabbitMQ queue declared")

	// Purge queue before publishing to ensure clean test state
	_, err = ch.QueuePurge(q.Name, false)
	if err != nil {
		log.Fatalf("Failed to purge queue: %v", err)
	}
	log.Println("Queue purged for clean publisher start")

	// Start publisher in goroutine
	done := make(chan bool)
	go publisher.Publish(ch, q.Name, done)
	<-done
	log.Println("Initial messages published")

	// Start consumer in goroutine
	go consumer.Consume(ch, q.Name)
	log.Println("Consumer started")

	// Publish duplicate message to test primary key violation handling
	duplicate := publisher.Payment{UserID: 1, PaymentID: 1, DepositAmount: 10}
	body, err := json.Marshal(duplicate)
	if err != nil {
		log.Fatalf("Failed to marshal duplicate message: %v", err)
	}
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		log.Printf("Failed to publish duplicate message: %v", err)
	} else {
		log.Println("Duplicate message published for testing")
	}

	// Start HTTP server for health check
	go startHealthCheckServer()

	// Graceful shutdown: block until CTRL+C pressed
	log.Println("Application is running. Press CTRL+C to exit.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down gracefully.")
}

func startHealthCheckServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Health check server running on :8080/health")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start health check server: %v", err)
	}
}
