package services

import (
	"errors"
	"testing"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Claims), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	t.Run("User not found", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userRepo.On("FindByEmail", "nonexistent@email.com").Return(nil, gorm.ErrRecordNotFound)

		token, err := authService.Login("nonexistent@email.com", "password")

		assert.Empty(t, token)
		assert.EqualError(t, err, "user not found")
		userRepo.AssertExpectations(t)
	})

	t.Run("Error fetching user", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userRepo.On("FindByEmail", "test@email.com").Return(nil, errors.New("connection error"))

		token, err := authService.Login("test@email.com", "password")

		assert.Empty(t, token)
		assert.EqualError(t, err, "failed to fetch user")
		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userID := uuid.New()
		
		correctPassword := "correct_password"
		hash, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
		assert.NoError(t, err)

		user := &models.User{
			ID:           userID,
			Email:        "test@email.com",
			PasswordHash: string(hash),
		}

		userRepo.On("FindByEmail", "test@email.com").Return(user, nil)

		token, err := authService.Login("test@email.com", "wrong_password")

		assert.Empty(t, token)
		assert.EqualError(t, err, "invalid credentials")
		userRepo.AssertExpectations(t)
	})

	t.Run("Successful login", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userID := uuid.New()
		user := &models.User{
			ID:    userID,
			Email: "test@email.com",
		}
		
		err := user.SetPassword("correct_password")
		assert.NoError(t, err)

		userRepo.On("FindByEmail", "test@email.com").Return(user, nil)
		jwtService.On("GenerateToken", userID).Return("generated_token", nil)

		token, err := authService.Login("test@email.com", "correct_password")

		assert.NoError(t, err)
		assert.Equal(t, "generated_token", token)
		userRepo.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
}

func TestAuthService_Signup(t *testing.T) {
	t.Run("User already exists", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		existingUser := &models.User{Email: "existing@email.com"}
		userRepo.On("FindByEmail", "existing@email.com").Return(existingUser, nil)

		newUser := &models.User{Email: "existing@email.com", Password: "password123"}
		token, err := authService.Signup(newUser)

		assert.Empty(t, token)
		assert.EqualError(t, err, "user already exists")
		userRepo.AssertExpectations(t)
	})

	t.Run("Error checking user existence", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userRepo.On("FindByEmail", "test@email.com").Return(nil, errors.New("connection error"))

		user := &models.User{Email: "test@email.com", Password: "password123"}
		token, err := authService.Signup(user)

		assert.Empty(t, token)
		assert.EqualError(t, err, "failed to check user existence")
		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid password", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userRepo.On("FindByEmail", "test@email.com").Return(nil, gorm.ErrRecordNotFound)

		user := &models.User{Email: "test@email.com", Password: "123"}
		token, err := authService.Signup(user)

		assert.Empty(t, token)
		assert.ErrorContains(t, err, "password")
		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid email", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userRepo.On("FindByEmail", "invalid-email").Return(nil, gorm.ErrRecordNotFound)

		user := &models.User{Email: "invalid-email", Password: "validPassword123"}
		token, err := authService.Signup(user)

		assert.Empty(t, token)
		assert.ErrorContains(t, err, "email")
		userRepo.AssertExpectations(t)
	})

	t.Run("Error creating user", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userRepo.On("FindByEmail", "test@email.com").Return(nil, gorm.ErrRecordNotFound)
		userRepo.On("Create", mock.Anything).Return(errors.New("save error"))

		user := &models.User{Email: "test@email.com", Password: "validPassword123"}
		token, err := authService.Signup(user)

		assert.Empty(t, token)
		assert.EqualError(t, err, "failed to create user")
		userRepo.AssertExpectations(t)
	})

	t.Run("Successful signup", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		jwtService := new(MockJWTService)
		authService := NewAuthService(userRepo, jwtService)

		userID := uuid.New()
		
		userRepo.On("FindByEmail", "test@email.com").Return(nil, gorm.ErrRecordNotFound)
		userRepo.On("Create", mock.Anything).Run(func(args mock.Arguments) {
			user := args.Get(0).(*models.User)
			user.ID = userID
		}).Return(nil)
		
		jwtService.On("GenerateToken", userID).Return("generated_token", nil)

		user := &models.User{
			Email:    "test@email.com",
			Password: "validPassword123",
		}
		token, err := authService.Signup(user)

		assert.NoError(t, err)
		assert.Equal(t, "generated_token", token)
		userRepo.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
}