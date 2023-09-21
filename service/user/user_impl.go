package user

import (
	"StoryGoAPI/DTO"
	myErrors "StoryGoAPI/errors"
	"StoryGoAPI/repository/redis"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/cache/v9"
	"net/http"
)

func (s *SrvImpl) Delete(ctx context.Context, id int) (bool, error) {
	err := s.repo.Psql.DeleteUser(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SrvImpl) Login(ctx context.Context, email, pass string) (*DTO.LoginResponse, error) {
	user, err := s.repo.Psql.GetUser(ctx, email)
	if err != nil {
		return nil, err
	}

	var token string
	if comparePassword(user.Password, pass) {
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "id", user.ID)
		token, err = s.auth.GenerateToken(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid password")
	}

	return &DTO.LoginResponse{
		Email: email,
		Token: token,
	}, nil
}

func (s *SrvImpl) Register(ctx context.Context, request *DTO.RegisterRequest) (*DTO.LoginResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, myErrors.NewErr(http.StatusBadRequest, err)
	}
	//ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	//defer cancel()
	pass, err := hashPassword(request.Password)
	if err != nil {
		return nil, err
	}
	res, err := s.repo.Psql.NewUser(ctx, request.ToModel(pass))

	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "email", request.Email)
	ctx = context.WithValue(ctx, "id", res.ID)
	token, err := s.auth.GenerateToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &DTO.LoginResponse{
		Email: res.Email,
		Token: token,
	}, nil
}

func (s *SrvImpl) DeleteStory(ctx context.Context, sid, uid int) error {
	lKey := redis.GenUserListKey(uid)
	_, err := s.repo.Psql.DeleteStory(ctx, sid)
	if err != nil {
		return err
	}
	if err := s.repo.Redis.InvalidateCache(ctx, lKey); err != nil {
		return err
	}
	return nil
}

func (s *SrvImpl) EditStory(ctx context.Context, req *DTO.StoryRequest) (*DTO.StoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	model := req.ToModel()
	story, err := s.repo.Psql.EditStory(ctx, model)
	if err != nil {
		return nil, err
	}

	lKey := redis.GenUserListKey(req.CreatorUserID)
	if err := s.repo.Redis.InvalidateCache(ctx, lKey); err != nil {
		return nil, err
	}

	sKey := redis.GenStoryKey(req.StoryID)
	if err = s.repo.Redis.Delete(ctx, sKey); err != nil {
		return nil, err
	}
	resp := &DTO.StoryResponse{}
	resp.FromModel(story)
	return resp, nil

}

func (s *SrvImpl) NewStory(ctx context.Context, req *DTO.StoryRequest) (*DTO.StoryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	lKey := redis.GenUserListKey(req.CreatorUserID)
	err := s.repo.Redis.InvalidateCache(ctx, lKey)
	if err != nil {
		return nil, err
	}

	model := req.ToModel()
	story, tokens, err := s.repo.Psql.NewStory(ctx, model, req.CreatorUserID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tokens); i++ {
		lKey = redis.GenGuestListKey(tokens[i])
		if err = s.repo.Redis.InvalidateCache(ctx, lKey); err != nil {
			return nil, err
		}
	}

	resp := &DTO.StoryResponse{}
	resp.FromModel(story)

	return resp, nil
}

func (s *SrvImpl) AllPostedStories(ctx context.Context, uid int, opt *DTO.UserStoryOption) (*DTO.StoriesResponse, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	resp := &DTO.StoriesResponse{}
	redisKey := redis.GenUserRedisKey(uid, opt.Limit, opt.Offset, opt.FromTime, opt.ToTime, opt.SortBy)
	listKey := redis.GenUserListKey(uid)
	err := s.repo.Redis.Get(ctx, redisKey, &resp)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			err = s.repo.Redis.Set(ctx, redisKey, listKey, &resp, func(_ *cache.Item) (any, error) {
				dbOpts := opt.ToModel()
				stories, err := s.repo.Psql.AllPostedStories(ctx, uid, dbOpts)
				if err != nil {
					return nil, err
				}

				var r = &DTO.StoriesResponse{
					Stories: make([]*DTO.StoryResponse, len(stories)),
				}
				for i := 0; i < len(stories); i++ {
					r.Stories[i] = &DTO.StoryResponse{}
					r.Stories[i].FromModel(stories[i])
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
	resp.SetRelevance()
	return resp, err
}
