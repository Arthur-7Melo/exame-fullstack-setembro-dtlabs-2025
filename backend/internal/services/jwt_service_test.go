package services

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

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
		userID := uuid.New()

		expiredTime := time.Now().Add(-1 * time.Hour)
		claims := &Claims{
			UserID: userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiredTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test-secret-key"))
		assert.NoError(t, err)

		_, err = jwtService.ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Equal(t, "invalid token", err.Error())
	})

	t.Run("Fail to validate token with wrong secret", func(t *testing.T) {
		jwtService1 := NewJWTService("secret-key-1")
		userID := uuid.New()
		token, err := jwtService1.GenerateToken(userID)
		assert.NoError(t, err)

		jwtService2 := NewJWTService("secret-key-2")
		_, err = jwtService2.ValidateToken(token)
		assert.Error(t, err)
		assert.Equal(t, "invalid token", err.Error())
	})

	t.Run("Fail to validate malformed token", func(t *testing.T) {
		jwtService := NewJWTService("test-secret-key")

		_, err := jwtService.ValidateToken("malformed.token.here")
		assert.Error(t, err)
		assert.Equal(t, "invalid token", err.Error())
	})

	t.Run("Fail to validate token with wrong signing method", func(t *testing.T) {
		jwtService := NewJWTService("test-secret-key")
		
		// Generate RSA key for testing wrong signing method
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		assert.NoError(t, err)

		userID := uuid.New()
		claims := &Claims{
			UserID: userID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
		}

		// Create token with RSA method instead of HS256
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenString, err := token.SignedString(privateKey)
		assert.NoError(t, err)

		_, err = jwtService.ValidateToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("Extract user ID from valid token", func(t *testing.T) {
		jwtService := NewJWTService("test-secret-key")
		userID := uuid.New()

		token, err := jwtService.GenerateToken(userID)
		assert.NoError(t, err)

		claims, err := jwtService.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})
}

func TestJWTService_TokenExpiration(t *testing.T) {
	t.Run("Token has correct expiration time", func(t *testing.T) {
		jwtService := NewJWTService("test-secret-key")
		userID := uuid.New()

		beforeGeneration := time.Now()
		token, err := jwtService.GenerateToken(userID)
		assert.NoError(t, err)

		claims, err := jwtService.ValidateToken(token)
		assert.NoError(t, err)

		// Check if expiration is approximately 24 hours from now
		expectedExpiration := beforeGeneration.Add(24 * time.Hour)
		actualExpiration := time.Unix(claims.ExpiresAt, 0)
		
		// Allow for small time difference due to execution time
		assert.InDelta(t, expectedExpiration.Unix(), actualExpiration.Unix(), 2)
	})
}