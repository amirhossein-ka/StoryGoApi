package entrepo

import (
	"StoryGoAPI/ent"
	"StoryGoAPI/ent/guestuser"
	"StoryGoAPI/ent/story"
	"StoryGoAPI/ent/user"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserExists = errors.New("user already exists")
)

func (r *Repo) NewUser(ctx context.Context, u ent.User) (*ent.User, error) {
	exists, err := r.client.User.Query().Where(user.Email(u.Email)).Exist(ctx)
	if err != nil {
		return nil, ErrUserExists
	}

	if !exists {
		addedUser, err := r.client.User.Create().
			SetName(u.Name).
			SetEmail(u.Email).
			SetPassword(u.Password).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("error creating user: %w", err)
		}

		return addedUser, nil
	}

	return nil, ErrUserExists
}

func (r *Repo) GetUser(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().Where(
		user.And(user.Email(email), user.DeletedAtIsNil()),
	).First(ctx)
}

func (r *Repo) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
	return r.client.User.Query().Where(
		user.And(user.ID(id), user.DeletedAtIsNil()),
	).First(ctx)
}

func (r *Repo) DeleteUser(ctx context.Context, id int) error {
	return r.client.User.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (r *Repo) NewStory(ctx context.Context, s *ent.Story, ownerID int) (*ent.Story, []string, error) {
	u, err := r.client.User.Get(ctx, ownerID)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := r.client.User.QueryFollowedBy(u).Select(guestuser.FieldToken).Strings(ctx)
	if err != nil {
		return nil, nil, err
	}

	_story, err := r.client.Story.Create().SetPostedby(u).
		SetStoryName(s.StoryName).
		SetBackgroundColor(s.BackgroundColor).
		SetBackgroundImage(s.BackgroundImage).
		SetIsShareable(s.IsShareable).
		SetExternalWebLink(s.ExternalWebLink).
		SetFromTime(s.FromTime).
		SetToTime(s.ToTime).
		SetStatus(s.Status).
		Save(ctx)

	if err != nil {
		return nil, nil, err
	}
	return _story, tokens, nil
}

func (r *Repo) DeleteStory(ctx context.Context, storyID int) (*ent.Story, error) {
	s, err := r.client.Story.UpdateOneID(storyID).Where(story.DeletedAtIsNil()).SetDeletedAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, err
	}

	return s, err
}

func (r *Repo) EditStory(ctx context.Context, s *ent.Story) (*ent.Story, error) {
	_story, err := r.client.Story.UpdateOneID(s.ID).Where(story.DeletedAtIsNil()).
		SetStoryName(s.StoryName).
		SetBackgroundColor(s.BackgroundColor).
		SetBackgroundImage(s.BackgroundImage).
		SetIsShareable(s.IsShareable).
		SetExternalWebLink(s.ExternalWebLink).
		SetFromTime(s.FromTime).SetToTime(s.ToTime).
		SetStatus(s.Status).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return _story, nil
}

func (r *Repo) AllPostedStories(ctx context.Context, id int, option *UserStoryOption) (ent.Stories, error) {
	u, err := r.client.User.Query().Where(
		user.And(user.ID(id), user.DeletedAtIsNil()),
	).Only(ctx)
	if err != nil {
		return nil, err
	}

	query := r.client.Story.Query().
		Where(
			story.HasPostedbyWith(user.ID(u.ID)), // Filter stories by followed users.
			story.DeletedAtIsNil(),
		).
		Offset(option.Offset). // Apply offset for pagination.
		Limit(option.Limit)

	// Date filters
	if !option.FromTime.IsZero() {
		query.Where(story.FromTimeGTE(option.FromTime))
	}

	if !option.ToTime.IsZero() {
		query.Where(story.ToTimeLTE(option.ToTime))
	}

	// sort by the given option
	switch option.SortBy {
	case "created":
		query.Order(ent.Desc(story.FieldCreatedAt))
	case "updated":
		query.Order(ent.Desc(story.FieldUpdatedAt))
	case "from_time":
		query.Order(ent.Desc(story.FieldFromTime))
	case "to_time":
		query.Order(ent.Desc(story.FieldToTime))
	}

	res, err := query.All(ctx)

	if err != nil {
		return nil, err
	}
	return res, nil
}

type UserStoryOption struct {
	Token    string    // Required
	Limit    int       // Optional
	Offset   int       // Optional
	SortBy   string    // Optional
	FromTime time.Time // Optional
	ToTime   time.Time // Optional
	Status   story.Status
}
