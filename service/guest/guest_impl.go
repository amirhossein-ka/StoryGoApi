package guest

import (
	"StoryGoAPI/DTO"
	"StoryGoAPI/ent"
	myErrors "StoryGoAPI/errors"
	"StoryGoAPI/repository/redis"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-redis/cache/v9"
	"net/http"

	"github.com/google/uuid"
)

// CheckGuestToken implements service.Guest
func (s *SrvImpl) CheckGuestToken(ctx context.Context) (bool, error) {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return false, fmt.Errorf("repo: token is invalid")
	}
	_, err := s.repo.Psql.GuestUserByToken(ctx, token)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteGuest implements service.Guest
func (s *SrvImpl) DeleteGuest(ctx context.Context, token string) (bool, error) {
	return s.repo.Psql.DeleteGuestUser(ctx, token)
}

// NewGuest implements service.Guest
func (s *SrvImpl) NewGuest(ctx context.Context, gr DTO.GuestRequest) (DTO.GuestResponse, error) {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return DTO.GuestResponse{}, fmt.Errorf("failed to generate uuid: %w", err)
	}

	hash := sha256.Sum256([]byte(randomUUID.String()))
	token := fmt.Sprintf("%x", hash)

	gu := gr.ToModel(gr.UserAgent, token[:])
	guestUser, err := s.repo.Psql.NewGuestUser(ctx, &gu)
	if err != nil {
		return DTO.GuestResponse{}, fmt.Errorf("failed to create guest user: %w", err)
	}

	return DTO.GuestResponse{
		UserAgent: guestUser.UserAgent,
		Token:     guestUser.Token,
	}, nil
}

func (s *SrvImpl) StoryFeed(ctx context.Context, gso *DTO.GuestStoryOption) (*DTO.StoriesResponse, error) {
	if err := gso.Validate(); err != nil {
		return nil, myErrors.NewErr(http.StatusBadRequest, err)
	}
	resp := &DTO.StoriesResponse{}
	redisKey := redis.GenGuestRedisKey(gso.Token, gso.SortBy, gso.Limit, gso.Offset, gso.FromTime, gso.ToTime)
	listKey := redis.GenGuestListKey(gso.Token)
	err := s.repo.Redis.Get(ctx, redisKey, &resp)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			err = s.repo.Redis.Set(ctx, redisKey, listKey, &resp, func(_ *cache.Item) (any, error) {
				r := &DTO.StoriesResponse{}

				opts := gso.ToModel()

				sr, err := s.repo.Psql.AllScannedStories(ctx, opts)
				if err != nil {
					return nil, err
				}
				r.Stories = make([]*DTO.StoryResponse, len(sr))
				for i := 0; i < len(sr); i++ {
					r.Stories[i] = &DTO.StoryResponse{}
					r.Stories[i].FromModel(sr[i])
					r.Stories[i].Status = "" // guests should not see the status of posts
				}
				return r, nil
			})

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// cache is ok
	resp.SetRelevance()
	return resp, err
}

// ScanStory implements service.Story
// if error is nil, operation done successfully.
func (s *SrvImpl) ScanStory(ctx context.Context, token string, id int) error {
	err := s.repo.Psql.ScanStory(ctx, token, id)
	if err != nil {
		var notfound *ent.NotFoundError
		if errors.As(err, &notfound) {
			return myErrors.NewErr(http.StatusNotFound, notfound)
		}
		return err
	}
	key := redis.GenGuestListKey(token)
	if err = s.repo.Redis.InvalidateCache(ctx, key); err != nil {
		return err
	}

	return nil
}
