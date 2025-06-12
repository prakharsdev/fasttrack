package publisher

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// Payment represents the payment message structure
type Payment struct {
	UserID        int `json:"user_id"`
	PaymentID     int `json:"payment_id"`
	DepositAmount int `json:"deposit_amount"`
}

// Publish sends initial messages to RabbitMQ
func Publish(ch *amqp.Channel, queueName string, done chan bool) {
	payments := []Payment{
		{1, 1, 10},
		{1, 2, 20},
		{2, 3, 20},
	}

	for _, p := range payments {
		body, err := json.Marshal(p)
		if err != nil {
			log.Printf("Failed to marshal payment: %v", err)
			continue
		}

		err = ch.Publish("", queueName, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			log.Printf("Published message: %+v", p)
		}
	}

	done <- true
}
