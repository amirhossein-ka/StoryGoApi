package echo

import (
	"StoryGoAPI/DTO"
	"StoryGoAPI/service/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate a user with their credentials and retrieve an access token.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		DTO.RegisterRequest	true	"User login credentials"
//
//	@Success		201		{object}	DTO.LoginResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		500		{object}	DTO.ErrorResponse
//
//	@Router			/user/login [post]
func (h *handler) login(c echo.Context) error {
	req := new(DTO.RegisterRequest)

	err := c.Bind(req)
	if err != nil {
		return err
	}
	response, err := h.srv.User().Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

// register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with the provided details.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		DTO.RegisterRequest	true	"User registration details"
//	@Success		200		{object}	DTO.LoginResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		500		{object}	DTO.ErrorResponse
//	@Router			/user/register [post]
func (h *handler) register(c echo.Context) error {
	req := new(DTO.RegisterRequest)

	if err := c.Bind(req); err != nil {
		return err
	}
	lr, err := h.srv.User().Register(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, lr)
}

// Delete godoc
//
//	@Summary		Delete user account
//	@Description	Delete the user account associated with the provided JWT token.
//	@Tags			user
//	@Security		JwtAuth
//	@Accept			json
//	@Produce		json
//
//	@Success		200	{object}	DTO.SuccessResponse
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Failure		500	{object}	DTO.ErrorResponse
//
//	@Router			/user/delete [delete]
func (h *handler) delete(c echo.Context) error {
	token, ok := c.Get("claims").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, DTO.ErrorResponse{Error: "token is invalid"})
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, DTO.ErrorResponse{Error: "token is invalid"})

	}

	deleted, err := h.srv.User().Delete(c.Request().Context(), claims.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, DTO.ErrorResponse{Error: "token is invalid"})
	}

	if err = h.srv.Auth().BlackList(c.Request().Context()); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, DTO.SuccessResponse{Success: deleted})
}

// newStory godoc
//
//	@Summary		Create a new story
//	@Description	Create a new story with the provided details.
//	@Tags			user, story
//	@Security		JwtAuth
//	@Accept			json
//	@Produce		json
//	@Param			story	body		DTO.CreateStoryRequest	true	"New story details"
//	@Success		201		{object}	DTO.StoryResponse
//	@Failure		400		{object}	DTO.ErrorResponse
//	@Failure		401		{object}	DTO.ErrorResponse
//	@Router			/user/new_story [post]
func (h *handler) newStory(c echo.Context) error {
	req := new(DTO.StoryRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	token, ok := c.Get("claims").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get token from header")
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get claims from token")
	}

	req.CreatorUserID = claims.ID

	sr, err := h.srv.User().NewStory(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, sr)
}

// allPostedStories godoc
//
//	@Summary		Get all posted stories of the user
//	@Description	Retrieve all posted stories of the user based on specified options.
//	@Tags			story, user
//	@Security		JwtAuth
//	@Accept			json
//	@Produce		json
//	@Param			sort_by		query		string	false	"Sort the stories by a field (e.g., 'created')"
//	@Param			limit		query		int		false	"Limit the number of returned stories"
//	@Param			offset		query		int		false	"Offset the returned stories"
//	@Param			from_date	query		string	false	"Filter stories by start date (2006-01-02T15:04:05Z07:00)"
//	@Param			to_date		query		string	false	"Filter stories by end date (2006-01-02T15:04:05Z07:00)"
//	@Success		200			{object}	DTO.StoriesResponse
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Failure		401			{object}	DTO.ErrorResponse
//	@Router			/user/stories [get]
//	@x-order		5
func (h *handler) allPostedStories(c echo.Context) error {
	options := &DTO.UserStoryOption{}
	err := (&echo.DefaultBinder{}).BindQueryParams(c, options)
	if err != nil {
		return err
	}

	token, ok := c.Get("claims").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get token from header")
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get claims from token")
	}

	sr, err := h.srv.User().AllPostedStories(c.Request().Context(), claims.ID, options)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sr)
}

// editStory godoc
//
//	@Summary		Edit a story
//	@Description	Edit an existing story with the provided details.
//	@Tags			story, user
//	@Security		JwtAuth
//	@Accept			json
//	@Produce		json
//	@Param			story_id	path		int					true	"Story ID to be edited"
//	@Param			story		body		DTO.StoryRequest	true	"Story details to be edited"
//	@Success		200			{object}	DTO.StoryResponse
//	@Failure		400			{object}	DTO.ErrorResponse
//	@Failure		401			{object}	DTO.ErrorResponse
//	@Failure		500			{object}	DTO.ErrorResponse
//	@Router			/api/v1/user/edit_story/{story_id} [put]
func (h *handler) editStory(c echo.Context) error {
	req := new(DTO.StoryRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	token, ok := c.Get("claims").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get token from header")
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get claims from token")
	}
	req.CreatorUserID = claims.ID

	sid, err := strconv.Atoi(c.Param("story_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid story id")
	}
	req.StoryID = sid

	sr, err := h.srv.User().EditStory(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sr)
}

// deleteStory godoc
//
//	@Summary		Delete a story
//	@Description	Delete a story with the provided ID.
//	@Tags			story, user
//	@Security		JwtAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Story ID to be deleted"
//	@Success		200	{object}	DTO.SuccessResponse
//	@Failure		400	{object}	DTO.ErrorResponse
//	@Failure		401	{object}	DTO.ErrorResponse
//	@Router			/user/delete_story/{id} [delete]
//	@x-order		8
func (h *handler) deleteStory(c echo.Context) error {
	var id int
	idStr := c.Param("id")
	if idStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, DTO.ErrorResponse{Error: "id not set"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, DTO.ErrorResponse{Error: "id is not int"}).SetInternal(err)
	}
	token, ok := c.Get("claims").(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get token from header")
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get claims from token")
	}

	if err = h.srv.User().DeleteStory(c.Request().Context(), id, claims.ID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, DTO.SuccessResponse{Success: true})
}
