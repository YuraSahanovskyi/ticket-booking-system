package handler

import (
	"github.com/YuraSahanovskyi/booking-system/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService    service.AuthService
	eventService   service.EventService
	bookingService service.BookingService
}

func NewHandler(auth service.AuthService, event service.EventService, booking service.BookingService) *Handler {
	return &Handler{
		authService:    auth,
		eventService:   event,
		bookingService: booking,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	// Додаємо CORS або Logger за потреби тут

	api := router.Group("/api")
	{
		h.initAuthRoutes(api)
		h.initEventRoutes(api)
		h.initBookingRoutes(api)
		h.initPaymentRoutes(api)
	}

	return router
}
