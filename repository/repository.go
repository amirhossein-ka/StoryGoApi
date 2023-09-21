package repository

import (
	"StoryGoAPI/config"
	"StoryGoAPI/ent"
	"StoryGoAPI/repository/entrepo"
	"StoryGoAPI/repository/redis"
	"context"
)

type (
	Repository interface {
		User() User
		GuestUser() GuestUser
		Story() Story
		Cache() Cache
	}

	User interface {
		NewUser(ctx context.Context, u ent.User) (*ent.User, error)
		GetUser(ctx context.Context, email string) (*ent.User, error)
		GetUserByID(ctx context.Context, id int) (*ent.User, error)
		DeleteUser(ctx context.Context, id int) error
	}
	GuestUser interface {
		NewGuestUser(ctx context.Context, gu *ent.GuestUser) (*ent.GuestUser, error)
		GuestUserByToken(ctx context.Context) (*ent.GuestUser, error)
		GuestUserByID(ctx context.Context, id int) (*ent.GuestUser, error)
		DeleteGuestUser(ctx context.Context) (bool, error)
	}
	Story interface{}

	Cache interface {
		Set(ctx context.Context, key string, value any) error
		Get(ctx context.Context, key string) (any, error)
	}
)
type Repo struct {
	Psql  *entrepo.Repo
	Redis *redis.CacheImpl
}

func (r *Repo) Close() error {
	if err := r.Psql.Close(); err != nil {
		return err
	}
	if err := r.Redis.Close(); err != nil {
		return err
	}

	return nil
}

func NewRepo(cfg *config.Config) (*Repo, error) {
	r, err := entrepo.NewRepo(cfg.DataBase)
	if err != nil {
		return nil, err
	}

	ci, err := redis.NewCache(cfg.Redis)
	if err != nil {
		return nil, err
	}
	return &Repo{Psql: r, Redis: ci}, nil
}
