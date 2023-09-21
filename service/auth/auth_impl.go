package auth

import (
	myErrors "StoryGoAPI/errors"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// BlackList implements service.Auth
func (s *SrvImpl) BlackList(ctx context.Context) error {
	tokenString := ctx.Value("token").(string)
	return s.repo.Redis.BlackList(ctx, tokenString)
}

func (s *SrvImpl) CheckToken(ctx context.Context) bool {
	tokenString := ctx.Value("token").(string)
	return s.repo.Redis.Exists(ctx, tokenString)
}

// ClaimsFromToken deprecated
func (s *SrvImpl) ClaimsFromToken(ctx context.Context) (any, error) {
	tokenString := ctx.Value("token").(string)
	// check if token in banned.
	// I guess this part have no use since we have no logout function. BUT we have a delete function !
	c, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	if exists := s.repo.Redis.Exists(c, tokenString); exists {
		// jwt key is blacklisted.
		return nil, myErrors.NewErr(http.StatusUnauthorized, fmt.Errorf("token is blacklisted"))
	}

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.secrets.JwtSecret), nil
	}, jwt.WithLeeway(time.Second*5))

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return token.Claims, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, fmt.Errorf("token is malformed")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		// Invalid signature
		return nil, fmt.Errorf("token signature is invalid")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, fmt.Errorf("token is expired/not active")
	} else {
		return nil, fmt.Errorf("could not handle this token: %w", err)
	}
}

func (s *SrvImpl) GenerateToken(ctx context.Context) (string, error) {
	id := ctx.Value("id").(int)

	expTime := time.Now().Add(s.secrets.ExpTime)
	claims := &JwtClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	jc := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jc.SignedString([]byte(s.secrets.JwtSecret))
}
