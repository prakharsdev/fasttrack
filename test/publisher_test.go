package test

import (
	"encoding/json"
	"testing"

	"fasttrack/internal/publisher"
)

func TestPaymentMarshalling(t *testing.T) {
	payment := publisher.Payment{
		UserID:        1,
		PaymentID:     42,
		DepositAmount: 100,
	}

	data, err := json.Marshal(payment)
	if err != nil {
		t.Errorf("Failed to marshal payment struct: %v", err)
	}

	var result publisher.Payment
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON data: %v", err)
	}

	if result != payment {
		t.Errorf("Expected %+v, got %+v", payment, result)
	}
}

func TestPublisherPaymentsPayload(t *testing.T) {
	payments := []publisher.Payment{
		{1, 1, 10},
		{1, 2, 20},
		{2, 3, 20},
	}

	if len(payments) != 3 {
		t.Errorf("Expected 3 payments, got %d", len(payments))
	}

	expected := publisher.Payment{1, 1, 10}
	if payments[0] != expected {
		t.Errorf("First payment payload mismatch. Got %+v, want %+v", payments[0], expected)
	}
}
