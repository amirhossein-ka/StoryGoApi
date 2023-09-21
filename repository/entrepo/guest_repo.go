package entrepo

import (
	"StoryGoAPI/ent"
	"StoryGoAPI/ent/guestuser"
	"StoryGoAPI/ent/story"
	"StoryGoAPI/ent/user"
	"context"
	"time"
)

func (r *Repo) NewGuestUser(ctx context.Context, gu *ent.GuestUser) (*ent.GuestUser, error) {
	createdGuestUser, err := r.client.GuestUser.
		Create().
		SetToken(gu.Token).
		SetVersionNumber(gu.VersionNumber).
		SetOperationSystem(gu.OperationSystem).
		SetUserAgent(gu.UserAgent).
		SetDisplayDetails(gu.DisplayDetails).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return createdGuestUser, nil
}

func (r *Repo) GuestUserByToken(ctx context.Context, token string) (*ent.GuestUser, error) {
	foundGuestUser, err := r.client.GuestUser.
		Query().
		Where(guestuser.And(guestuser.Token(token), guestuser.DeletedAtIsNil())).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return foundGuestUser, nil
}

func (r *Repo) GuestUserByID(ctx context.Context, id int) (*ent.GuestUser, error) {
	foundGuestUser, err := r.client.GuestUser.Query().
		Where(
			guestuser.And(
				guestuser.DeletedAtIsNil(), guestuser.ID(id),
			),
		).First(ctx)

	if err != nil {
		return nil, err
	}
	return foundGuestUser, nil
}

func (r *Repo) DeleteGuestUser(ctx context.Context, token string) (bool, error) {
	err := r.client.GuestUser.Update().Where(
		guestuser.And(
			guestuser.Token(token), guestuser.DeletedAtIsNil(),
		),
	).SetDeletedAt(time.Now()).Exec(ctx)

	return err == nil, err
}

// ScanStory adds a relation like a guest followed a user.
func (r *Repo) ScanStory(ctx context.Context, token string, id int) error {
	u, err := r.client.Story.Query().Where(story.ID(id), story.StatusEQ(story.StatusPublic)).QueryPostedby().Only(ctx)
	if err != nil {
		return err
	}
	gu, err := r.client.GuestUser.Query().Where(guestuser.Token(token)).Only(ctx)
	if err != nil {
		return err
	}

	u, err = u.Update().AddFollowedBy(gu).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// AllScannedStories used for getting stories from users that a guest followed.
func (r *Repo) AllScannedStories(ctx context.Context, option *GuestStoryOption) (ent.Stories, error) {
	followedUserIDs, err := r.client.GuestUser.Query().Where(guestuser.Token(option.Token)).
		QueryFollowed().Select(user.FieldID).Ints(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch all timed stories posted by the followed users, sorted by creation time.
	query := r.client.Story.Query().
		Where(
			story.HasPostedbyWith(user.IDIn(followedUserIDs...)), // Filter stories by followed users.
			story.DeletedAtIsNil(),
			story.StatusEQ(story.StatusPublic),
		).
		Offset(option.Offset). // Apply offset for pagination.
		Limit(option.Limit)    // Apply limit for pagination.

	// Date filters
	if !option.FromTime.IsZero() {
		query.Where(story.FromTimeGTE(option.FromTime))
	}

	if !option.ToTime.IsZero() {
		query.Where(story.FromTimeLTE(option.ToTime))
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

	timedStories, err := query.All(ctx)

	if err != nil {
		return nil, err
	}

	return timedStories, nil
}

// GuestStoryOption represents the options for fetching guest stories.
type GuestStoryOption struct {
	Token    string    // Required
	Limit    int       // Optional
	Offset   int       // Optional
	SortBy   string    // Optional
	FromTime time.Time // Optional
	ToTime   time.Time // Optional
}
