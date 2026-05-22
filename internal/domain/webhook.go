package domain

import "github.com/google/uuid"

type PaymentWebhookPayload struct {
	OrderID uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}
