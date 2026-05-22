package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type paymentWebhookInput struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	Status  string    `json:"status" binding:"required"`
}

func (h *Handler) initPaymentRoutes(api *gin.RouterGroup) {
	payments := api.Group("/payments")
	{
		payments.POST("/webhook", h.handlePaymentWebhook)
	}
}

func (h *Handler) handlePaymentWebhook(c *gin.Context) {
	var input paymentWebhookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment payload"})
		return
	}

	if input.Status != "success" {
		c.JSON(http.StatusOK, gin.H{"message": "non-success status ignored"})
		return
	}

	err := h.bookingService.ConfirmPayment(c.Request.Context(), input.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not confirm payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment confirmed"})
}