package echo

import (
	"StoryGoAPI/service/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var echoJwtConf = echojwt.Config{
	// SigningKey:    []byte(r.cfg.Secrets.JwtSecret),
	ContextKey:    "claims",
	SigningMethod: "HS256",
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return &auth.JwtClaims{}
	},
}

func (r *rest) routing() {
	if r.cfg.Server.Docs.Path != "" {
		r.echo.GET(r.cfg.Server.Docs.Path, echoSwagger.WrapHandler)
	}

	echoJwtConf.SigningKey = []byte(r.cfg.Secrets.JwtSecret)
	{
		user := r.echo.Group("/api/v1/user")

		user.POST("/login", r.handler.login)
		user.POST("/register", r.handler.register)
		{
			authUsers := user.Group(
				"",
				echojwt.WithConfig(echoJwtConf), // to add claims in context
				r.handler.addTokenToContext(),   // to add token(string) for blacklist
				r.handler.checkBlackList(),      // check if a token blacklisted
			)

			authUsers.POST("/new_story", r.handler.newStory)
			authUsers.GET("/stories", r.handler.allPostedStories)
			authUsers.PUT("/edit_story/:story_id", r.handler.editStory)
			authUsers.DELETE("/delete_story/:id", r.handler.deleteStory)
			authUsers.DELETE("/delete", r.handler.delete)
		}
	}

	{
		guest := r.echo.Group("/api/v1/guest")

		guest.POST("/new", r.handler.getGuestToken)
		guest.POST("/verify", r.handler.checkToken)
		guest.DELETE("/delete", r.handler.deleteGuest)
		guest.POST("/scan/:id", r.handler.scanStory)
		guest.GET("/stories", r.handler.storyFeed) // should be able to filter by dates
	}

	story := r.echo.Group("/api/v1/story")
	{
		story.GET("/:id", r.handler.storyInfo)
	}
}
