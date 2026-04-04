package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey  string
	Issuer     string
	Subject    string
	ExpireDays int
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func ParseToken(tokenStr, secretKey string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func GenerateToken(userID int64, email string, cfg JWTConfig) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   cfg.Subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}
