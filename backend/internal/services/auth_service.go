package services

import (
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
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
	userRepo userRepository
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
        if err == gorm.ErrRecordNotFound {
            return "", errors.ErrUserNotFound
        }
        return "", errors.ErrFailedToCheckUser
    }

    if !user.CheckPassword(password) {
        return "", errors.ErrInvalidCredentials
    }

    return s.jwtService.GenerateToken(user.ID)
}

func (s *authService) Signup(user *models.User) (string, error) {
    _, err := s.userRepo.FindByEmail(user.Email)
    if err == nil {
        return "", errors.ErrUserAlreadyExists
    }
    if err != gorm.ErrRecordNotFound {
        return "", errors.ErrFailedToCheckUser
    }

    if err := utils.ValidatePassword(user.Password); err != nil {
        return "", errors.ErrWeakPassword
    }
    if err := utils.ValidateEmail(user.Email); err != nil {
        return "", errors.ErrInvalidEmail
    }

    if err := user.SetPassword(user.Password); err != nil {
        return "", errors.ErrFailedToSetPassword
    }

    if err := s.userRepo.Create(user); err != nil {
        return "", errors.ErrFailedToCreateUser
    }

    return s.jwtService.GenerateToken(user.ID)
}