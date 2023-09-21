package DTO

import (
	"StoryGoAPI/ent"
	"StoryGoAPI/ent/story"
	myErrors "StoryGoAPI/errors"
	"context"
	"fmt"
	"net/http"
	"time"
)

type (
	StoriesResponse struct {
		Stories []*StoryResponse `json:"stories"`
	}

	StoryResponse struct {
		StoryID         int       `json:"storyID"`
		CreatorUserID   int       `json:"creatorUserId"`
		FromTime        time.Time `json:"fromTime,omitempty"`
		ToTime          time.Time `json:"toTime,omitempty"`
		StoryName       string    `json:"storyName"`
		BackgroundColor string    `json:"backgroundColor"`
		BackgroundImage string    `json:"backgroundImage"`
		IsShareable     bool      `json:"isShareable"`
		AttachedFile    string    `json:"attachedFile,omitempty"`
		ExternalWebLink string    `json:"externalWebLink,omitempty"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
		Relevance       bool      `json:"relevance"`
		Status          Status    `json:"status,omitempty"`
	}
	StoryRequest struct {
		StoryID         int       `json:"storyID"`
		CreatorUserID   int       `json:"creatorUserId"`
		FromTime        time.Time `json:"fromTime,omitempty"`
		ToTime          time.Time `json:"toTime,omitempty"`
		StoryName       string    `json:"storyName"`
		BackgroundColor string    `json:"backgroundColor"`
		BackgroundImage string    `json:"backgroundImage"`
		IsShareable     bool      `json:"isShareable"`
		AttachedFile    string    `json:"attachedFile,omitempty"`
		ExternalWebLink string    `json:"externalWebLink,omitempty"`
		Status          Status    `json:"status"`
	}

	CreateStoryRequest struct {
		FromTime        time.Time `json:"fromTime,omitempty" example:"2020-01-01T00:00:00Z"`
		ToTime          time.Time `json:"toTime,omitempty" example:"2020-12-31T23:59:59Z"`
		StoryName       string    `json:"storyName" example:"My first story"`
		BackgroundColor string    `json:"backgroundColor" example:"#RRGGBB"`
		BackgroundImage string    `json:"backgroundImage" example:"http://example.com/myimage.jpg"`
		IsShareable     bool      `json:"isShareable" example:"true"`
		AttachedFile    string    `json:"attachedFile,omitempty" example:"http://example.com/myfile.pdf"`
		ExternalWebLink string    `json:"externalWebLink,omitempty" example:"http://example.com"`
		Status          Status    `json:"status" example:"private/public"`
	}
)

type Status string

const (
	StatusPrivate Status = "private"
	StatusPublic  Status = "public"
)

func (s *StoriesResponse) SetRelevance() {
	if len(s.Stories) <= 0 {
		return
	}

	currentTime := time.Now()
	var nearestIndex int
	nearestDuration := time.Duration(0)

	for i := 0; i < len(s.Stories); i++ {
		if s.Stories[i].FromTime.After(currentTime) {
			duration := s.Stories[i].FromTime.Sub(currentTime)
			if nearestDuration == 0 || duration < nearestDuration {
				nearestIndex = i
				nearestDuration = duration
			}
		}
	}

	s.Stories[nearestIndex].Relevance = true
}

func (s *StoryResponse) FromModel(story *ent.Story) {
	s.StoryID = story.ID
	s.CreatorUserID = story.QueryPostedby().FirstIDX(context.Background())
	s.FromTime = story.FromTime
	s.FromTime = story.FromTime
	s.ToTime = story.ToTime
	s.ToTime = story.ToTime
	s.StoryName = story.StoryName
	s.BackgroundColor = story.BackgroundColor
	s.BackgroundImage = story.BackgroundImage
	s.IsShareable = story.IsShareable
	s.AttachedFile = story.AttachedFile
	s.ExternalWebLink = story.ExternalWebLink
	s.CreatedAt = story.CreatedAt
	s.UpdatedAt = story.UpdatedAt
	s.Status = Status(story.Status)
}

func (s *StoryRequest) ToModel() *ent.Story {
	return &ent.Story{
		ID:              s.StoryID,
		FromTime:        s.FromTime,
		ToTime:          s.ToTime,
		StoryName:       s.StoryName,
		BackgroundColor: s.BackgroundColor,
		BackgroundImage: s.BackgroundImage,
		IsShareable:     s.IsShareable,
		AttachedFile:    s.AttachedFile,
		ExternalWebLink: s.ExternalWebLink,
		Status:          story.Status(s.Status),
	}
}

func (s *StoryRequest) Validate() error {
	// Validate required fields
	if s.StoryName == "" {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("story name is required"))
	}
	if s.BackgroundColor == "" {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("background color is required"))
	}
	if s.BackgroundImage == "" {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("background image is required"))
	}
	if s.CreatorUserID == 0 {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("creator user ID is required"))
	}
	switch s.Status {
	case StatusPrivate, StatusPublic:
	default:
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("invalid status %s", s.Status))
	}

	// Validate date and time fields for timed stories
	if s.ToTime.IsZero() && s.ToTime.IsZero() && s.FromTime.IsZero() && s.FromTime.IsZero() {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("times cannot be zero when creating stories"))
	}
	if s.FromTime.IsZero() {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("from date is required for timed stories"))
	}
	if s.ToTime.IsZero() {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("to date is required for timed stories"))
	}
	if s.FromTime.IsZero() {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("from time is required for timed stories"))
	}
	if s.ToTime.IsZero() {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("to time is required for timed stories"))
	}

	if s.FromTime.After(s.ToTime) || s.FromTime.Equal(s.ToTime) {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("from date must be before to date"))
	}
	if s.FromTime.After(s.ToTime) || s.FromTime.Equal(s.ToTime) {
		return myErrors.NewErr(http.StatusBadRequest, fmt.Errorf("from time must be before to time"))
	}

	return nil
}
