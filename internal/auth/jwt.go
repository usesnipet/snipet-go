package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/model"
)

type CUserClaims struct {
	jwt.RegisteredClaims

	ClientCode string `json:"client_code"`
}

type JWTService struct {
	config config.AuthConfig
}

func NewJWTService(config config.AuthConfig) *JWTService {
	return &JWTService{config: config}
}

func (s *JWTService) GenerateToken(
	clientCode string,
	cUser *model.CUser,
) (string, CUserClaims, error) {
	claims := CUserClaims{
		ClientCode: clientCode,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.JWTIssuer,
			Subject:   cUser.ID.String(),
			Audience:  jwt.ClaimStrings{s.config.JWTAudience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWTExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", CUserClaims{}, err
	}

	return tokenString, claims, nil
}

func (s *JWTService) VerifyToken(tokenString string) (*CUserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CUserClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return &CUserClaims{}, err
	}
	return token.Claims.(*CUserClaims), nil
}
