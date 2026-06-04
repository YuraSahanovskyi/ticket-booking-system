package dto

import "github.com/google/uuid"

type PaymentWebhookRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	Status  string    `json:"status" binding:"required"`
}
