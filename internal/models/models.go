package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusPending PaymentStatus = "pending"
	StatusSuccess PaymentStatus = "success"
	StatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID        uuid.UUID     `json:"id"`
	BookingID uuid.UUID     `json:"booking_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Amount    float64       `json:"amount"`
	Status    PaymentStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type CreatePaymentInput struct {
	BookingID uuid.UUID
	UserID    uuid.UUID
	Amount    float64
}
