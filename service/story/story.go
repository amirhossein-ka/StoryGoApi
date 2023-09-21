package story

import (
	"StoryGoAPI/repository"
)

type SrvImpl struct {
	repo *repository.Repo
}

func NewStorySrv(r *repository.Repo) *SrvImpl {
	return &SrvImpl{
		repo: r,
	}
}
