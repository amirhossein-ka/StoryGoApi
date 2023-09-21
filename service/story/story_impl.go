package story

import (
	"StoryGoAPI/DTO"
	"StoryGoAPI/repository/redis"
	"context"
	"errors"
	"github.com/go-redis/cache/v9"
)

func (s *SrvImpl) StoryInfo(ctx context.Context, id int) (*DTO.StoryResponse, error) {
	resp := &DTO.StoryResponse{}
	key := redis.GenStoryKey(id)
	err := s.repo.Redis.Get(ctx, key, &resp)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			err = s.repo.Redis.Set(ctx, key, "stories", &resp, func(item *cache.Item) (any, error) {
				info, err := s.repo.Psql.StoryInfo(ctx, id)
				if err != nil {
					return nil, err
				}
				r := &DTO.StoryResponse{}
				r.FromModel(info)
				return r, nil
			})

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return resp, nil
}
