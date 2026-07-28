package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/model"
)

type UserClaims struct {
	jwt.RegisteredClaims

	ClientCode string `json:"client_code"`
}

type JWTService struct {
	config config.AuthConfig
}

func NewJWTService(config config.AuthConfig) *JWTService {
	return &JWTService{config: config}
}

func (s *JWTService) GenerateToken(clientCode string, user *model.User) (string, time.Time, UserClaims, error) {
	expiresAt := time.Now().Add(s.config.JWTExpiration)
	claims := UserClaims{
		ClientCode: clientCode,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.JWTIssuer,
			Subject:   user.ID,
			Audience:  jwt.ClaimStrings{s.config.JWTAudience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", time.Time{}, UserClaims{}, err
	}

	return fmt.Sprintf("Bearer %s", tokenString), expiresAt, claims, nil
}

func (s *JWTService) VerifyToken(tokenString string) (*UserClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return &UserClaims{}, err
	}
	return token.Claims.(*UserClaims), nil
}
