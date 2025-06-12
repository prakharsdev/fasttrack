package test

import (
	"encoding/json"
	"testing"

	"fasttrack/internal/consumer"
)

func TestConsumerUnmarshalValid(t *testing.T) {
	input := `{"user_id": 1, "payment_id": 10, "deposit_amount": 50}`
	var payment consumer.Payment
	err := json.Unmarshal([]byte(input), &payment)
	if err != nil {
		t.Fatalf("Failed to unmarshal valid JSON: %v", err)
	}

	if payment.UserID != 1 || payment.PaymentID != 10 || payment.DepositAmount != 50 {
		t.Errorf("Unexpected unmarshalled result: %+v", payment)
	}
}

func TestConsumerUnmarshalInvalid(t *testing.T) {
	input := `{"user_id": "invalid", "payment_id": 10, "deposit_amount": 50}`
	var payment consumer.Payment
	err := json.Unmarshal([]byte(input), &payment)
	if err == nil {
		t.Errorf("Expected unmarshal error but got none")
	}
}
