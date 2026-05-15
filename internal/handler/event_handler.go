package handler

import (
	"net/http"

	"github.com/YuraSahanovskyi/booking-system/internal/handler/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) initEventRoutes(api *gin.RouterGroup) {
	events := api.Group("/events")
	{
		events.GET("/", h.getAllEvents)
		events.GET("/:id/seats", h.getEventSeats)
	}
}

// @Summary      Get all events
// @Description  Get a list of all upcoming events
// @Tags         events
// @Produce      json
// @Success      200  {array}   dto.EventResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /events [get]
func (h *Handler) getAllEvents(c *gin.Context) {
	events, err := h.eventService.GetAllEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.ToEventsResponse(events))
}

// @Summary      Get event seats
// @Description  Get detailed information about an event and its seats availability
// @Tags         events
// @Produce      json
// @Param        id   path      string  true  "Event ID (UUID)"
// @Success      200  {object}  dto.EventWithSeatsResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /events/{id}/seats [get]
func (h *Handler) getEventSeats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("invalid id"))
		return
	}

	event, seats, err := h.eventService.GetEventWithSeats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.ToEventWithSeatsResponse(*event, seats))
}
