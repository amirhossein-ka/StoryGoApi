package echo

import (
	"StoryGoAPI/config"
	"StoryGoAPI/controller"
	"StoryGoAPI/service"
	"context"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type rest struct {
	cfg     *config.Config
	echo    *echo.Echo
	handler *handler
}

// Start implements controller.Rest
func (r *rest) Start() error {
	//r.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Format: "method=${method}, URI=${uri}, code=${status}, latency=${latency_human}\n",
	//}))
	//r.echo.Pre(middleware.RemoveTrailingSlash())

	r.echo.Use(r.handler.costumeContext)

	r.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
			http.MethodDelete,
		},
		AllowHeaders: []string{
			"X-Requested-With",
			"X-Guest-Token",
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}))

	// set our costumeContext err handler
	r.echo.HTTPErrorHandler = r.ErrHandler
	if r.cfg.Server.Debug != "" {
		r.echo.Debug = true
	}

	r.routing()

	return r.echo.Start(r.cfg.Server.Listen)
}

// Stop implements controller.Rest
func (r *rest) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return r.echo.Shutdown(ctx)
}

type handler struct {
	srv service.Service
}

func NewRest(cfg *config.Config, service service.Service) controller.Rest {
	return &rest{
		echo: echo.New(),
		cfg:  cfg,
		handler: &handler{
			srv: service,
		},
	}
}
