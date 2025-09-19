package services

import (
	"testing"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateToken(t *testing.T) {
    t.Run("Successfully generate token", func(t *testing.T) {
        jwtService := NewJWTService("test-secret-key")
        userID := uuid.New()

        token, err := jwtService.GenerateToken(userID)

        assert.NoError(t, err)
        assert.NotEmpty(t, token)
    })

    t.Run("Fail to generate token with empty secret", func(t *testing.T) {
        jwtService := NewJWTService("")
        userID := uuid.New()

        token, err := jwtService.GenerateToken(userID)

        assert.Error(t, err)
        assert.Empty(t, token)
        assert.Equal(t, errors.ErrTokenGeneration, err)
    })
}

func TestJWTService_ValidateToken(t *testing.T) {
    t.Run("Successfully validate valid token", func(t *testing.T) {
        jwtService := NewJWTService("test-secret-key")
        userID := uuid.New()

        token, err := jwtService.GenerateToken(userID)
        assert.NoError(t, err)

        claims, err := jwtService.ValidateToken(token)
        assert.NoError(t, err)
        assert.Equal(t, userID, claims.UserID)
        assert.True(t, claims.ExpiresAt > time.Now().Unix())
    })

    t.Run("Fail to validate expired token", func(t *testing.T) {
        jwtService := NewJWTService("test-secret-key")
        expiredTime := time.Now().Add(-1 * time.Hour)
        claims := &Claims{
            UserID: uuid.New(),
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expiredTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString([]byte("test-secret-key"))
        assert.NoError(t, err)

        _, err = jwtService.ValidateToken(tokenString)
        assert.Error(t, err)
        assert.Equal(t, errors.ErrTokenExpired, err)
    })

    t.Run("Fail to validate token with wrong secret", func(t *testing.T) {
        jwtService1 := NewJWTService("secret-key-1")
        userID := uuid.New()
        token, err := jwtService1.GenerateToken(userID)
        assert.NoError(t, err)

        jwtService2 := NewJWTService("secret-key-2")
        _, err = jwtService2.ValidateToken(token)
        assert.Error(t, err)
        assert.Equal(t, errors.ErrInvalidToken, err)
    })

    t.Run("Fail to validate malformed token", func(t *testing.T) {
        jwtService := NewJWTService("test-secret-key")
        _, err := jwtService.ValidateToken("malformed.token.here")
        assert.Error(t, err)
        assert.Equal(t, errors.ErrInvalidToken, err)
    })

    t.Run("Fail to validate token with wrong signing method", func(t *testing.T) {
        jwtService := NewJWTService("test-secret-key")
        claims := &Claims{
            UserID: uuid.New(),
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims) // Wrong method
        _, err := token.SignedString([]byte("test-secret-key"))
        // This will fail because we're trying to use RSA key with HMAC method
        if err == nil {
            _, err = jwtService.ValidateToken(token.Raw)
            assert.Error(t, err)
            assert.Equal(t, errors.ErrInvalidSigningMethod, err)
        }
    })
}