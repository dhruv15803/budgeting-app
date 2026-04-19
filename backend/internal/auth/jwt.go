package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSigner struct {
	secret []byte
	expiry time.Duration
}

func NewJWTSigner(secret string, expiry time.Duration) *JWTSigner {
	return &JWTSigner{secret: []byte(secret), expiry: expiry}
}

type Claims struct {
	UserID int    `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *JWTSigner) SignAccessToken(userID int, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   strconv.Itoa(userID),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *JWTSigner) ParseAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
