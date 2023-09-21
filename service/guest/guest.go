package guest

import (
	"StoryGoAPI/repository"
)

type SrvImpl struct {
	repo *repository.Repo
}

func NewGuestSrv(repo *repository.Repo) *SrvImpl {
	return &SrvImpl{repo: repo}
}
