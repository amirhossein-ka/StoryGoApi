package echo

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func getToken(c echo.Context) (string, error) {
	auth := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(auth, "Bearer")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("error getting jwt token from header: token is not in proper fomrat")
	}
	auth = strings.TrimSpace(splitToken[1])
	return auth, nil
}

func (h *handler) addTokenToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := getToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
			c.Set("token", token)
			return next(c)
		}
	}
}

func (h *handler) checkBlackList() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if blacklisted := h.srv.Auth().CheckToken(c.Request().Context()); blacklisted {
				return echo.NewHTTPError(http.StatusUnauthorized, "token blacklisted")
			}
			return next(c)
		}
	}
}

// this type and its methods are for echo.Context so it Set() and Get() methods operate on internal context.Context
type contextValue struct {
	echo.Context
}

func (ctx contextValue) Get(key string) any {
	val := ctx.Context.Get(key)
	if val != nil {
		return val
	}
	return ctx.Request().Context().Value(key)
}

func (ctx contextValue) Set(key string, val interface{}) {
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), key, val)))
}

func (h *handler) costumeContext(fn echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return fn(contextValue{ctx})
	}
}
