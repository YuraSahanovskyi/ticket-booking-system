package handler

import (
	"net/http"

	"github.com/YuraSahanovskyi/booking-system/internal/handler/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("unauthorized"))
		return
	}

	var input dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewValidationErrorResponse(err))
		return
	}

	seatID, err := uuid.Parse(input.SeatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("invalid seat id"))
		return
	}

	booking, err := h.bookingService.BookSeat(c.Request.Context(), userID, seatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.ToCreateBookingResponse(*booking))
}

func (h *Handler) getUserBookings(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok{
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("unauthorized"))
		return
	}

	bookings, err := h.bookingService.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("failed to fetch bookings"))
		return
	}

	c.JSON(http.StatusOK, dto.ToBookingsResponse(bookings))
}

func (h *Handler) cancelBooking(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("unauthorized"))
		return
	}
	
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("invalid id"))
		return
	}

	if err := h.bookingService.CancelBooking(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
