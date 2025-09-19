package handlers

import (
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.TokenResponse "Returns JWT token"
// @Failure 400 {object} dto.DetailedErrorResponse "Invalid request body"
// @Failure 403 {object} dto.ForbiddenErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var credentials dto.LoginRequest

    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(400, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    token, err := h.authService.Login(credentials.Email, credentials.Password)
    if err != nil {
        if customErr, ok := err.(errors.CustomError); ok {
            c.JSON(customErr.StatusCode(), dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeFromStatusCode(customErr.StatusCode()),
                Message: customErr.Message(),
                Details: "Please check your credentials and try again",
            })
        } else {
            c.JSON(500, dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeInternalError,
                Message: "Internal server error",
                Details: err.Error(),
            })
        }
        return
    }

    c.JSON(200, dto.TokenResponse{Token: token})
}

// Signup godoc
// @Summary User registration
// @Description Register a new user and return JWT token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param request body dto.SignupRequest true "User registration information"
// @Success 201 {object} dto.TokenResponse "Returns JWT token"
// @Failure 400 {object} dto.BadRequestErrorResponse "Invalid email format or password too short"
// @Failure 409 {object} dto.ConflictErrorResponse "User already exists"
// @Failure 500 {object} dto.InternalServerErrorResponse "Internal server error"
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
    var signupReq dto.SignupRequest
    if err := c.ShouldBindJSON(&signupReq); err != nil {
        c.JSON(400, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    user := models.User{
        Email:    signupReq.Email,
        Password: signupReq.Password,
    }

    token, err := h.authService.Signup(&user)
    if err != nil {
        if customErr, ok := err.(errors.CustomError); ok {
            c.JSON(customErr.StatusCode(), dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeFromStatusCode(customErr.StatusCode()),
                Message: customErr.Message(),
                Details: "Please check your input and try again",
            })
        } else {
            c.JSON(500, dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeInternalError,
                Message: "Internal server error",
                Details: err.Error(),
            })
        }
        return
    }

    c.JSON(201, dto.TokenResponse{Token: token})
}