// internal/services/jwt_service.go
package services

import (
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JWTService interface {
    GenerateToken(userID uuid.UUID) (string, error)
    ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    jwt.StandardClaims
}

type jwtService struct {
    jwtKey []byte
}

func NewJWTService(jwtSecret string) JWTService {
    return &jwtService{
        jwtKey: []byte(jwtSecret),
    }
}

func (s *jwtService) GenerateToken(userID uuid.UUID) (string, error) {
	 if len(s.jwtKey) == 0 {
		return "", errors.ErrTokenGeneration
	}
	
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", errors.ErrTokenGeneration
	}

	return tokenString, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.ErrInvalidSigningMethod
        }
        return s.jwtKey, nil
    })

    if err != nil {
        if ve, ok := err.(*jwt.ValidationError); ok {
            if ve.Errors&jwt.ValidationErrorExpired != 0 {
                return nil, errors.ErrTokenExpired
            }
        }
        return nil, errors.ErrInvalidToken
    }

    if !token.Valid {
        return nil, errors.ErrInvalidToken
    }

    return claims, nil
}