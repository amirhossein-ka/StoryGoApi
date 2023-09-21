package user

import (
	"StoryGoAPI/config"
	"StoryGoAPI/repository"
	"StoryGoAPI/service/auth"
)

type SrvImpl struct {
	repo *repository.Repo
	auth *auth.SrvImpl
}

func New(r *repository.Repo, auth *auth.SrvImpl, _ *config.Config) *SrvImpl {
	return &SrvImpl{
		repo: r,
		auth: auth,
	}
}
