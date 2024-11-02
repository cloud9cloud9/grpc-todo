package security

import (
	"errors"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//go:generate mockgen -source=helper.go -destination=mocks/mock.go

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

type AuthHelper interface {
	CompareHashAndPassword(hashedPass string, password []byte) bool
	GenerateToken(user *domain.User) (string, error)
	HashPassword(password string) string
	ValidateToken(token string) (claims *TokenClaims, err error)
}

type AuthUtil struct{}

func NewAuthUtil() *AuthUtil {
	return &AuthUtil{}
}

func (a *AuthUtil) CompareHashAndPassword(hashedPass string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), password)
	return err == nil
}
func (a *AuthUtil) GenerateToken(user *domain.User) (string, error) {
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

func (a *AuthUtil) HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashBytes)
}

func (a *AuthUtil) ValidateToken(token string) (claims *TokenClaims, err error) {
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
