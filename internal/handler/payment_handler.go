package handler

import (
	"net/http"

	"github.com/YuraSahanovskyi/booking-system/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initPaymentRoutes(api *gin.RouterGroup) {
	payments := api.Group("/payments")
	{
		payments.POST("/webhook", h.handlePaymentWebhook)
	}
}

func (h *Handler) handlePaymentWebhook(c *gin.Context) {
	var input dto.PaymentWebhookRequest
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
