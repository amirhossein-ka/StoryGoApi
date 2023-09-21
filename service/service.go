package service

import (
	"StoryGoAPI/DTO"
	"StoryGoAPI/config"
	"StoryGoAPI/repository"
	"StoryGoAPI/service/auth"
	"StoryGoAPI/service/guest"
	"StoryGoAPI/service/story"
	"StoryGoAPI/service/user"
	"context"
)

type (
	Service interface {
		Auth() Auth
		User() User
		Guest() Guest
		Story() Story
	}
	Story interface {
		StoryInfo(ctx context.Context, id int) (*DTO.StoryResponse, error)
	}
	User interface {
		Register(ctx context.Context, request *DTO.RegisterRequest) (*DTO.LoginResponse, error)
		Login(ctx context.Context, email, pass string) (*DTO.LoginResponse, error)
		Delete(ctx context.Context, id int) (bool, error)

		NewStory(ctx context.Context, req *DTO.StoryRequest) (*DTO.StoryResponse, error)
		EditStory(ctx context.Context, req *DTO.StoryRequest) (*DTO.StoryResponse, error)
		AllPostedStories(ctx context.Context, uid int, opt *DTO.UserStoryOption) (*DTO.StoriesResponse, error)
		DeleteStory(ctx context.Context, sid, uid int) error
	}
	Guest interface {
		NewGuest(ctx context.Context, gu DTO.GuestRequest) (DTO.GuestResponse, error)
		CheckGuestToken(ctx context.Context) (bool, error)
		DeleteGuest(ctx context.Context, token string) (bool, error)
		ScanStory(ctx context.Context, token string, id int) error
		StoryFeed(ctx context.Context, gso *DTO.GuestStoryOption) (*DTO.StoriesResponse, error)
	}
	Auth interface {
		CheckToken(ctx context.Context) bool
		GenerateToken(context.Context) (string, error)
		ClaimsFromToken(context.Context) (any, error)
		BlackList(context.Context) error
	}
	srv struct {
		user  *user.SrvImpl
		auth  *auth.SrvImpl
		guest *guest.SrvImpl
		story *story.SrvImpl
	}
)

// Guest implements Service
func (s *srv) Guest() Guest {
	return s.guest
}

// Auth implements Service
func (s *srv) Auth() Auth {
	return s.auth
}

func (s *srv) User() User {
	return s.user
}

func (s *srv) Story() Story {
	return s.story
}

func NewService(cfg *config.Config, repo *repository.Repo) Service {
	authImpl := auth.NewAuth(cfg.Secrets, repo)
	srv := &srv{
		user:  user.New(repo, authImpl, cfg),
		auth:  authImpl,
		guest: guest.NewGuestSrv(repo),
		story: story.NewStorySrv(repo),
	}

	return srv
}
