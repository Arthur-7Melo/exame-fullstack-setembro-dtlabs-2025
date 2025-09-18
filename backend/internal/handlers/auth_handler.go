package handlers

import (
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := h.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(403, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := h.authService.Signup(&user)
	if err != nil {
		statusCode := 500
		if err.Error() == "user already exists" {
			statusCode = 409
		} else if err.Error() == "invalid email format" || err.Error() == "password must be at least 8 characters" {
			statusCode = 400
		}
		
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"token": token})
}