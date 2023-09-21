package echo

import (
	"StoryGoAPI/DTO"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// getGuestToken godoc
//
//	@Summary		Generate a guest token
//	@Description	Generate a guest token for anonymous access.
//	@Tags			guest
//	@Accept			json
//	@Produce		json
//	@Param			guest	body		DTO.NewGuestReq	true	"Guest details"
//	@Success		201		{object}	DTO.GuestResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/guest/new [post]
func (h *handler) getGuestToken(c echo.Context) error {
	req := new(DTO.GuestRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	req.UserAgent = c.Request().UserAgent()
	gr, err := h.srv.Guest().NewGuest(c.Request().Context(), *req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, gr)
}

// checkToken godoc
//
//	@Summary		Check a guest token
//	@Description	Check the validity of a guest token.
//	@Tags			guest
//	@Accept			json
//	@Produce		json
//	@Param			token	body		DTO.VerifyReq	true	"Guest token to be checked"
//	@Success		200		{object}	DTO.SuccessResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		404		{object}	DTO.SuccessResponse
//	@Router			/guest/verify [post]
func (h *handler) checkToken(c echo.Context) error {
	req := new(DTO.GuestRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	c.Set("token", req.Token)

	b, err := h.srv.Guest().CheckGuestToken(c.Request().Context())
	if err != nil {
		c.Logger().Error("error", err)
		return c.JSON(http.StatusNotFound, DTO.SuccessResponse{Success: false})
	}

	return c.JSON(http.StatusOK, DTO.SuccessResponse{Success: b})
}

// deleteGuest godoc
//
//	@Summary	Delete a guest token
//	@Tags		guest
//	@Accept		json
//	@Security	GuestAuth
//	@Produce	json
//	@Success	200		{object}	DTO.SuccessResponse
//	@Failure	400		{object}	DTO.ErrorResponse
//	@Router		/guest/delete [delete]
func (h *handler) deleteGuest(c echo.Context) error {
	options := &DTO.GuestStoryOption{}
	err := (&echo.DefaultBinder{}).BindHeaders(c, options)
	if err != nil {
		return err
	}

	if options.Token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "no token found")
	}

	b, err := h.srv.Guest().DeleteGuest(c.Request().Context(), options.Token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, DTO.SuccessResponse{Success: b})
}

// scanStory godoc
//
//	@Summary		Scan a story
//	@Description	Scan a story using a guest token.
//	@Tags			story, guest
//	@Security		GuestAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Story ID to be scanned"
//	@Param			token	header		string	true	"Guest token for authentication"
//	@Success		200		{object}	DTO.SuccessResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Router			/guest/scan/{id} [post]
func (h *handler) scanStory(c echo.Context) error {
	options := &DTO.GuestStoryOption{}
	err := (&echo.DefaultBinder{}).BindHeaders(c, options)
	if err != nil {
		return err
	}

	if options.Token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "no token found")
	}

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "story id is not valid")
	}

	err = h.srv.Guest().ScanStory(c.Request().Context(), options.Token, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, DTO.SuccessResponse{Success: true})
}

// storyFeed godoc
//
//	@Summary		Get story feed
//	@Description	Retrieve a story feed based on specified options.
//	@Tags			story, guest
//	@Security		GuestAuth
//	@Accept			json
//	@Produce		json
//	@Param			sort_by		query		string	false	"Sort the stories by a field (e.g., 'created')"
//	@Param			limit		query		int		false	"Limit the number of returned stories"
//	@Param			offset		query		int		false	"Offset the returned stories"
//	@Param			from_date	query		string	false	"Filter stories by start date (e.g.: 2006-01-02T15:04:05)"
//	@Param			to_date		query		string	false	"Filter stories by end date (e.g: 2006-01-02T15:04:05)"
//	@Param			token		header		string	true	"Guest token for authentication"
//	@Success		200			{object}	DTO.StoriesResponse
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Router			/guest/stories [get]
func (h *handler) storyFeed(c echo.Context) error {
	options := &DTO.GuestStoryOption{}
	err := (&echo.DefaultBinder{}).BindQueryParams(c, options)
	if err != nil {
		return err
	}

	err = (&echo.DefaultBinder{}).BindHeaders(c, options)
	if err != nil {
		return err
	}

	feed, err := h.srv.Guest().StoryFeed(c.Request().Context(), options)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, feed)
}
