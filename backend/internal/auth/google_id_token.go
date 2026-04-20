package auth

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/idtoken"
)

type GoogleIDTokenClaims struct {
	Sub           string
	Email         string
	EmailVerified bool
	Picture       string
}

func ParseGoogleIDToken(ctx context.Context, rawToken, audience string) (*GoogleIDTokenClaims, error) {
	if rawToken == "" || audience == "" {
		return nil, errors.New("missing token or audience")
	}
	payload, err := idtoken.Validate(ctx, rawToken, audience)
	if err != nil {
		return nil, err
	}
	sub, _ := payload.Claims["sub"].(string)
	if sub == "" {
		return nil, fmt.Errorf("invalid token: missing sub")
	}
	email, _ := payload.Claims["email"].(string)
	if email == "" {
		return nil, fmt.Errorf("invalid token: missing email")
	}
	emailVerified := googleClaimBool(payload.Claims, "email_verified")
	picture, _ := payload.Claims["picture"].(string)
	return &GoogleIDTokenClaims{
		Sub:           sub,
		Email:         email,
		EmailVerified: emailVerified,
		Picture:       picture,
	}, nil
}

func googleClaimBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t == "true" || t == "1"
	default:
		return false
	}
}
