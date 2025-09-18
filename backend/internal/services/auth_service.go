package services

import (
	"errors"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
}

type AuthService interface {
	Login(email, password string) (string, error)
	Signup(user *models.User) (string, error)
}

type authService struct {
	userRepo userRepository  // Use a interface local
	jwtService JWTService
}

func NewAuthService(userRepo userRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", errors.New("failed to fetch user")
	}

	if !user.CheckPassword(password) {
		return "", errors.New("invalid credentials")
	}

	return s.jwtService.GenerateToken(user.ID)
}

func (s *authService) Signup(user *models.User) (string, error) {
	_, err := s.userRepo.FindByEmail(user.Email)
	if err == nil {
		return "", errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("failed to check user existence")
	}

	if err := utils.ValidatePassword(user.Password); err != nil {
		return "", err
	}
	if err := utils.ValidateEmail(user.Email); err != nil {
		return "", err
	}

	if err := user.SetPassword(user.Password); err != nil {
		return "", errors.New("failed to set password")
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", errors.New("failed to create user")
	}

	return s.jwtService.GenerateToken(user.ID)
}
