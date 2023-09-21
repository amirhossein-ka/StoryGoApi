package echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *handler) storyInfo(c echo.Context) error {
	storyID := c.Param("id")
	id, err := strconv.Atoi(storyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}
	storyInfo, err := h.srv.Story().StoryInfo(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, storyInfo)
}
