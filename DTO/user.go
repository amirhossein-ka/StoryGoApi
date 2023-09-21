package DTO

import (
	"StoryGoAPI/ent"
	"StoryGoAPI/repository/entrepo"
	"errors"
	"regexp"
	"time"
)

type (
	LoginResponse struct {
		Email string `json:"email" example:"user@example.com"`
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	}

	RegisterRequest struct {
		Email    string `json:"email" example:"user@example.com"`
		Password string `json:"password" example:"password123"`
		Name     string `json:"name" example:"John Doe"`
	}
)

func (l LoginResponse) FromModel(token string, u ent.User) {
	l.Email = u.Email
	// l.Password = u.Password
	l.Token = token
}

// ToModel converts RegisterRequest to ent.User, receive hashed password.
func (ur RegisterRequest) ToModel(pass string) ent.User {
	return ent.User{
		Email:    ur.Email,
		Password: pass,
		Name:     ur.Name,
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (ur RegisterRequest) Validate() error {
	// Email validation using a simple regular expression
	if !emailRegex.MatchString(ur.Email) {
		return errors.New("invalid email format")
	}

	// Password length validation
	if len(ur.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Name length and emptiness validation
	if len(ur.Name) < 4 {
		return errors.New("name must be at least 6 characters long")
	}

	return nil // No errors, validation passed successfully
}

type UserStoryOption struct {
	SortBy   string    `query:"sort_by"`   // Optional
	Limit    int       `query:"limit"`     // Optional
	Offset   int       `query:"offset"`    // Optional
	FromTime time.Time `query:"from_date"` // Optional
	ToTime   time.Time `query:"to_date"`   // Optional
	Status   Status    `query:"status"`
}

func (o *UserStoryOption) ToModel() *entrepo.UserStoryOption {
	return &entrepo.UserStoryOption{
		Limit:    o.Limit,
		Offset:   o.Offset,
		SortBy:   o.SortBy,
		FromTime: o.FromTime,
		ToTime:   o.ToTime,
	}
}

var userSortOptions = map[string]struct{}{
	"updated":   {},
	"created":   {},
	"from_time": {},
	"to_time":   {},
	"":          {},
}

// Validate checks if the GuestStoryOption is valid.
// It returns an error if the Token is empty, or if the FromTime is after the ToTime.
func (o *UserStoryOption) Validate() error {
	if _, valid := userSortOptions[o.SortBy]; !valid {
		return errors.New("sortBy field must be equal to 'updated', 'created', 'from_time', 'to_time', or an empty string")
	}

	switch o.Status {
	case StatusPrivate, StatusPublic:
	default:
		return errors.New("invalid status")
	}

	// Set default values if optional fields are not provided.
	if o.Limit == 0 || o.Limit < 0 {
		o.Limit = 10 // Default limit to 10 if not provided.
	}

	if o.Offset < 0 {
		o.Offset = 0 // Default offset to 0 if not provided or negative.
	}

	// Check if FromTime is before or equal to ToTime.
	if !o.FromTime.IsZero() && !o.ToTime.IsZero() && o.FromTime.After(o.ToTime) {
		return errors.New("FromTime must be before or equal to ToTime")
	}

	return nil
}
