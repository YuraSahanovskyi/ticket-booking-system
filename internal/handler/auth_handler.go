package handler

import (
	"net/http"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}
}

func (h *Handler) register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewValidationErrorResponse(err))
		return
	}

	id, err := h.authService.Register(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{ID: id.String()})
}

func (h *Handler) login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewValidationErrorResponse(err))
		return
	}

	token, err := h.authService.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, dto.NewErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal error"))
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}