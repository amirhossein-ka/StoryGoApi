package DTO

import (
	"StoryGoAPI/ent"
	"StoryGoAPI/repository/entrepo"
	"errors"
	"time"
)

type GuestRequest struct {
	Token           string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	VersionNumber   int    `json:"versionNumber" example:"1"`
	OperatingSystem string `json:"operatingSystem" example:"Linux"`
	DisplayDetails  string `json:"displayDetails" example:"1920x1080"`
	UserAgent       string `json:"userAgent" example:"Mozilla/5.0"`
}

func (gr GuestRequest) ToModel(ua string, token string) ent.GuestUser {
	return ent.GuestUser{
		VersionNumber:   gr.VersionNumber,
		OperationSystem: gr.OperatingSystem,
		DisplayDetails:  gr.DisplayDetails,
		UserAgent:       ua,
		Token:           token,
	}
}

type NewGuestReq struct {
	VersionNumber   int    `json:"versionNumber" example:"1"`
	OperatingSystem string `json:"operatingSystem" example:"Linux"`
	DisplayDetails  string `json:"displayDetails" example:"1920x1080"`
}

type VerifyReq struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

type GuestResponse struct {
	UserAgent string `json:"userAgent" example:"Mozilla/5.0"`
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

func (gr GuestResponse) FromModel(gu ent.GuestUser) {
	gr.UserAgent = gu.UserAgent
	gr.Token = gu.Token
}

type GuestStoryOption struct {
	Token    string    `query:"token" header:"X-Guest-Token"` // Required
	SortBy   string    `query:"sort_by"`                      // Optional
	Limit    int       `query:"limit"`                        // Optional
	Offset   int       `query:"offset"`                       // Optional
	FromTime time.Time `query:"from_date"`                    // Optional
	ToTime   time.Time `query:"to_date"`                      // Optional
}

func (o *GuestStoryOption) ToModel() *entrepo.GuestStoryOption {
	return &entrepo.GuestStoryOption{
		Token:    o.Token,
		Limit:    o.Limit,
		Offset:   o.Offset,
		SortBy:   o.SortBy,
		FromTime: o.FromTime,
		ToTime:   o.ToTime,
	}
}

var sortOptions = map[string]struct{}{
	"updated":   {},
	"created":   {},
	"from_time": {},
	"to_time":   {},
	"":          {},
}

// Validate checks if the GuestStoryOption is valid.
// It returns an error if the Token is empty, or if the FromTime is after the ToTime.
func (o *GuestStoryOption) Validate() error {
	if o.Token == "" {
		return errors.New("token cannot be empty")
	}

	if _, valid := sortOptions[o.SortBy]; !valid {
		return errors.New("SortBy field must be equal to 'updated', 'created', 'from_time', 'to_time', or an empty string")
	}

	// Set default values if optional fields are not provided.
	if o.Limit == 0 {
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
