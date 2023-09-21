package auth

import (
	"StoryGoAPI/config"
	"StoryGoAPI/repository"

	"github.com/golang-jwt/jwt/v5"
)

type SrvImpl struct {
	secrets config.Secrets
	repo    *repository.Repo
}

type JwtClaims struct {
	ID int
	jwt.RegisteredClaims
}

func NewAuth(scr config.Secrets, r *repository.Repo) *SrvImpl {
	return &SrvImpl{
		secrets: scr,
		repo:    r,
	}
}
