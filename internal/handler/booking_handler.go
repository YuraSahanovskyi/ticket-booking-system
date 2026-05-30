package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
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

// @Summary      Create a booking
// @Description  Reserve a seat for a specific event
// @Tags         bookings
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateBookingRequest  true  "Seat ID"
// @Success      201    {object}  dto.CreateBookingResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Failure      401    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
// @Router       /bookings [post]
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
		switch {
		case errors.Is(err, domain.ErrSeatAlreadyBooked):
			c.JSON(http.StatusConflict, dto.NewErrorResponse("seat already booked"))
		case strings.Contains(err.Error(), "already booked") ||
			strings.Contains(err.Error(), "locked"):
			c.JSON(http.StatusConflict, dto.NewErrorResponse(err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("failed to create booking"))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.ToCreateBookingResponse(*booking))
}

// @Summary      Get my bookings
// @Description  Get a list of all bookings for the current user
// @Tags         bookings
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}   dto.BookingResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /bookings [get]
func (h *Handler) getUserBookings(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
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

// @Summary      Cancel booking
// @Description  Cancel an existing booking by its ID
// @Tags         bookings
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Booking ID (UUID)"
// @Success      204  "No Content"
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /bookings/{id} [delete]
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
