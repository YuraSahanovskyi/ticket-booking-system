package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type bookingInput struct {
	SeatID uuid.UUID `json:"seat_id" binding:"required"`
}

func (h *Handler) initBookingRoutes(api *gin.RouterGroup) {
	bookings := api.Group("/bookings", h.userIdentity)
	{
		bookings.POST("/", h.createBooking)
		bookings.GET("/", h.getUserBookings)
		bookings.DELETE("/:id", h.cancelBooking)
	}
}

func (h *Handler) createBooking(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input bookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	booking, err := h.bookingService.BookSeat(c.Request.Context(), userID, input.SeatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *Handler) getUserBookings(c *gin.Context) {
	userID, _ := getUserID(c)
	bookings, err := h.bookingService.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *Handler) cancelBooking(c *gin.Context) {
	userID, _ := getUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.bookingService.CancelBooking(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
