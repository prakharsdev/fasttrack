package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"fasttrack/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initializes the MySQL connection with retry logic and creates required tables
func InitDB(cfg config.Config) {
	var err error

	// Build connection string from config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDatabase,
	)

	// Retry logic
	for i := 1; i <= 5; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err == nil {
			err = DB.Ping()
		}

		if err == nil {
			log.Println("Successfully connected to MySQL")
			break
		}

		log.Printf("MySQL connection failed (attempt %d/5): %v", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to MySQL after retries: %v", err)
	}

	createTables()
	log.Println("MySQL tables ensured.")
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS payment_events (
			user_id INT,
			payment_id INT PRIMARY KEY,
			deposit_amount INT
		);`,
		`CREATE TABLE IF NOT EXISTS skipped_messages (
			user_id INT,
			payment_id INT PRIMARY KEY,
			deposit_amount INT
		);`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
}
