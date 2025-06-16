package services

import (
	"fmt"
	"strconv"
	"time"
	"ya41-56/internal/gophermart/customerror"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	secretKey     string
	tokenLifetime time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"id"`
}

func NewTokenService(secretKey string, tokenLifetime time.Duration) *TokenService {
	return &TokenService{
		secretKey:     secretKey,
		tokenLifetime: tokenLifetime,
	}
}

func (s *TokenService) BuildJWTString(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenLifetime)),
		},
		UserID: strconv.Itoa(int(userID)),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", customerror.ErrJWTToken
	}

	return tokenString, nil
}

func (s *TokenService) ParseToken(claims *Claims, tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(s.secretKey), nil
		})

	if err != nil || !token.Valid {
		return nil, customerror.ErrJWTToken
	}

	return token, nil
}
