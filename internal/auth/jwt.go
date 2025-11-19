package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT payload for API usage (mobile app).
type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken builds a signed JWT and returns the raw string.
func GenerateToken(signingKey []byte, userID, role string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return token.SignedString(signingKey)
}
