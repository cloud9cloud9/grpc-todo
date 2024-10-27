package utils

import (
	"errors"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)
import "time"

var (
	tokenTTL                = 12 * time.Hour
	signingKey              = viper.GetString("JWT_SECRET")
	issuer                  = viper.GetString("JWT_ISSUER")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}

func GenerateToken(user domain.User) (string, error) {
	claims := jwt.Claims(
		&TokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(jwt.NewNumericDate(time.Now()).Add(tokenTTL)),
				Issuer:    issuer,
			},
			UserId: user.Id,
			Email:  user.Email,
		})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(token string) (claims *TokenClaims, err error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidSigningMethod
			}
			return []byte(signingKey), nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
