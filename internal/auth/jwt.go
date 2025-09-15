package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "user"

var (
	jwtKey       []byte
	ErrKeyNotSet = errors.New("JWT key not set")
)

// SetJWTKey sets key for signing JWT tokens
func SetJWTKey(key string) {
	if key == "" {
		panic("JWT key cannot be empty")
	}
	jwtKey = []byte(key)
	log.Printf("JWT key set with length: %d", len(jwtKey))
}

// GenerateToken creates new JWT token
func GenerateToken(username string) (string, error) {
	if len(jwtKey) == 0 {
		return "", ErrKeyNotSet
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	log.Printf("Generated token for user %s", username)
	return tokenString, nil
}

// ValidateToken validates JWT token
func ValidateToken(tokenStr string) (*Claims, error) {
	if len(jwtKey) == 0 {
		return nil, ErrKeyNotSet
	}

	log.Printf("Validating token: %s", tokenStr[:10]) // Log only beginning of token

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Check that correct signing method is used
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Method)
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		log.Printf("Token validation error: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Printf("Token is invalid")
		return nil, errors.New("token is invalid")
	}

	log.Printf("Token validated successfully for user: %s", claims.Username)
	return claims, nil
}

// Claims represents JWT token data
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// SetUserContext adds user to context
func SetUserContext(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, userContextKey, username)
}

// GetUserFromContext gets user from context
func GetUserFromContext(ctx context.Context) string {
	if username, ok := ctx.Value(userContextKey).(string); ok {
		return username
	}
	return ""
}
