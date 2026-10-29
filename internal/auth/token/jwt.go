package token

import (
	"fmt"
	"time"

	"github.com/EmotionlessDev/GoEcho/internal/auth/apperror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret []byte, ttl time.Duration) *JWTService {
	return &JWTService{
		secret: secret,
		ttl:    ttl,
	}
}

func (s *JWTService) Create(userID uuid.UUID) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) Parse(tokenString string) (uuid.UUID, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, apperror.ErrInvalidToken
			}

			return s.secret, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !token.Valid {
		return uuid.Nil, apperror.ErrInvalidToken
	}

	return claims.UserID, nil
}
