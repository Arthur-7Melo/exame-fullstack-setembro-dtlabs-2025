package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
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
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "jwt-token-123", response["token"])
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
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.On("Login", "test@example.com", "wrongpassword").Return("", errors.New("invalid credentials"))

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
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "Invalid credentials", response["error"])
		mockAuthService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		mockAuthService.On("Login", "nonexistent@example.com", "password123").Return("", errors.New("user not found"))

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
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "Invalid credentials", response["error"])
		mockAuthService.AssertExpectations(t)
	})
}

func TestAuthHandler_Signup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful signup", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			Email:    "newuser@example.com",
			Password: "securePassword123",
		}

		mockAuthService.On("Signup", user).Return("jwt-token-456", nil)

		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "jwt-token-456", response["token"])
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
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("User already exists", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			Email:    "existing@example.com",
			Password: "password123",
		}

		mockAuthService.On("Signup", user).Return("", errors.New("user already exists"))

		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "user already exists", response["error"])
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid email format", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			Email:    "invalid-email",
			Password: "password123",
		}

		mockAuthService.On("Signup", user).Return("", errors.New("invalid email format"))

		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "invalid email format", response["error"])
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid password", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			Email:    "test@example.com",
			Password: "short",
		}

		mockAuthService.On("Signup", user).Return("", errors.New("password must be at least 8 characters"))

		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "password must be at least 8 characters", response["error"])
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Failed to create user", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			Email:    "test@example.com",
			Password: "validPassword123",
		}

		mockAuthService.On("Signup", user).Return("", errors.New("failed to create user"))

		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "failed to create user", response["error"])
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Failed to check user existence", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		handler := NewAuthHandler(mockAuthService)

		user := &models.User{
			Email:    "test@example.com",
			Password: "validPassword123",
		}

		mockAuthService.On("Signup", user).Return("", errors.New("failed to check user existence"))

		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Signup(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		
		assert.Equal(t, "failed to check user existence", response["error"])
		mockAuthService.AssertExpectations(t)
	})
}