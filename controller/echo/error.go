package echo

import (
	"StoryGoAPI/ent"
	"StoryGoAPI/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *rest) ErrHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var code int
	var msg any
	switch e := err.(type) {
	case *echo.HTTPError:
		if e.Internal != nil {
			if herr, ok := e.Internal.(*echo.HTTPError); ok {
				e = herr
			}
		}
		code = e.Code
		msg = echo.Map{"error": e.Message}

	case *errors.HttpErr:
		code = e.Code
		msg = echo.Map{"error": e.Error()}

	case *errors.ValidationErr:
		code = e.Code
		msg = echo.Map{"errors": err}

	case *ent.NotFoundError:
		code = http.StatusNotFound
		msg = echo.Map{"error": e.Error()}

	default:
		code = 500
		msg = echo.Map{"error": err.Error()}
	}

	// Send response
	if c.Request().Method == http.MethodHead { // Issue #608
		err = c.NoContent(code)
	} else {
		err = c.JSON(code, msg)
	}
	if err != nil {
		r.echo.Logger.Error(err)
	}
}
