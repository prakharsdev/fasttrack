package consumer

import (
	"encoding/json"
	"log"
	"strings"

	"fasttrack/internal/db"

	"github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

const MySQLDuplicateEntryCode = 1062

// Payment represents the payment message structure
type Payment struct {
	UserID        int `json:"user_id"`
	PaymentID     int `json:"payment_id"`
	DepositAmount int `json:"deposit_amount"`
}

// Consume reads messages from RabbitMQ and stores them in MySQL
func Consume(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for msg := range msgs {
		var p Payment
		if err := json.Unmarshal(msg.Body, &p); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Try inserting into payment_events
		_, err := db.DB.Exec("INSERT INTO payment_events VALUES (?, ?, ?)", p.UserID, p.PaymentID, p.DepositAmount)
		if err != nil {
			if isDuplicateKeyError(err) {
				// On primary key conflict insert into skipped_messages
				_, skipErr := db.DB.Exec("INSERT INTO skipped_messages VALUES (?, ?, ?)", p.UserID, p.PaymentID, p.DepositAmount)
				if skipErr != nil {
					log.Printf("Failed to insert into skipped_messages: %v", skipErr)
				} else {
					log.Printf("Duplicate detected, moved to skipped_messages: %+v", p)
				}
			} else {
				log.Printf("Failed to insert into payment_events: %v", err)
			}
		} else {
			log.Printf("Inserted payment: %+v", p)
		}
	}
}

// Explicit check for MySQL duplicate key error code 1062
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == MySQLDuplicateEntryCode
	}
	// fallback generic detection
	return strings.Contains(err.Error(), "Error 1062")
}
