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

// @title           Booking System API
// @version         1.0
// @description     API Server for Cinema/Event Booking System.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Yura Sahanovskyi
// @contact.url    https://github.com/YuraSahanovskyi

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer <your_jwt_token>" to authenticate
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
