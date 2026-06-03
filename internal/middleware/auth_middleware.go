package middleware

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"momentia-be/repository"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, secret string, expireHours int) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthMiddleware(secret string, sessionRepo repository.UserSessionRepository) gin.HandlerFunc {
	if secret == "" {
		panic("AuthMiddleware: JWT Secret must not be empty")
	}
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		hash := sha256.Sum256([]byte(tokenStr))
		tokenHash := fmt.Sprintf("%x", hash)
		if _, err := sessionRepo.FindByTokenHash(tokenHash); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token has been revoked"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
