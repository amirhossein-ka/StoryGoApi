package entrepo

import (
	"StoryGoAPI/ent"
	"StoryGoAPI/ent/story"
	"context"
)

func (r *Repo) StoryInfo(ctx context.Context, id int) (*ent.Story, error) {
	return r.client.Story.Query().Where(
		story.And(story.ID(id), story.DeletedAtIsNil()),
	).Only(ctx)
}
