package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Signup(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful login", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.On("Login", "test@example.com", "password123").Return("jwt-token-123", nil)

		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response dto.TokenResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "jwt-token-123", response.Token)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.On("Login", "test@example.com", "wrongpassword").Return("", errors.ErrInvalidCredentials)

		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "invalid credentials", response.Message)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.On("Login", "nonexistent@example.com", "password123").Return("", errors.ErrUserNotFound)

		loginData := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "user not found", response.Message)
		mockAuthService.AssertExpectations(t)
	})
}

func TestAuthHandler_Signup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful signup", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			Email:    "newuser@example.com",
			Password: "securePassword123",
		}

		mockAuthService.On("Signup", expectedUser).Return("jwt-token-456", nil)

		signupData := map[string]string{
			"email":    "newuser@example.com",
			"password": "securePassword123",
		}
		jsonData, _ := json.Marshal(signupData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response dto.TokenResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "jwt-token-456", response.Token)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("User already exists", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			Email:    "existing@example.com",
			Password: "password123",
		}

		mockAuthService.On("Signup", expectedUser).Return("", errors.ErrUserAlreadyExists)

		signupData := map[string]string{
			"email":    "existing@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(signupData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "user already exists", response.Message)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid email format", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			Email:    "invalid-email",
			Password: "password123",
		}

		mockAuthService.On("Signup", expectedUser).Return("", errors.ErrInvalidEmail)

		signupData := map[string]string{
			"email":    "invalid-email",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(signupData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "invalid email format", response.Message)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid password", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			Email:    "test@example.com",
			Password: "short",
		}

		mockAuthService.On("Signup", expectedUser).Return("", errors.ErrWeakPassword)

		signupData := map[string]string{
			"email":    "test@example.com",
			"password": "short",
		}
		jsonData, _ := json.Marshal(signupData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "password must be at least 8 characters", response.Message)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Failed to create user", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			Email:    "test@example.com",
			Password: "validPassword123",
		}

		mockAuthService.On("Signup", expectedUser).Return("", errors.ErrFailedToCreateUser)

		signupData := map[string]string{
			"email":    "test@example.com",
			"password": "validPassword123",
		}
		jsonData, _ := json.Marshal(signupData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "failed to create user", response.Message)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Failed to check user existence", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			Email:    "test@example.com",
			Password: "validPassword123",
		}

		mockAuthService.On("Signup", expectedUser).Return("", errors.ErrFailedToCheckUser)

		signupData := map[string]string{
			"email":    "test@example.com",
			"password": "validPassword123",
		}
		jsonData, _ := json.Marshal(signupData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "failed to check user existence", response.Message)
		mockAuthService.AssertExpectations(t)
	})
}