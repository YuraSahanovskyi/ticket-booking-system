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

func (h *Handler) getAllEvents(c *gin.Context) {
	events, err := h.eventService.GetAllEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.ToEventsResponse(events))
}

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
