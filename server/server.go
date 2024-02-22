package server

import (
	"fmt"
	"net/http"
	"strings"

	"verifydog/common"
	"verifydog/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OnAuthenticate func(authenticated common.Authenticated)

type Server struct {
	Environment    *common.Environment
	Router         *gin.Engine
	OnAuthenticate OnAuthenticate
}

func NewServer(environment *common.Environment, onAuthenticate OnAuthenticate) *Server {
	gin.SetMode(gin.ReleaseMode)

	return &Server{
		Environment:    environment,
		Router:         gin.Default(),
		OnAuthenticate: onAuthenticate,
	}
}

func (server *Server) Startup() error {
	server.Router.ForwardedByClientIP = false
	if err := server.Router.SetTrustedProxies(nil); err != nil {
		return err
	}

	server.Router.Use(middleware.RateLimitMiddleware)
	server.Router.GET("/", server.StatusHandler())
	server.Router.GET("/authorize", server.AuthorizeHandler())
	server.Router.NoRoute(middleware.NoRouteMiddleware)

	runPort := ":8080"
	if server.Environment.DevMode {
		runPort = ":3030"
	}

	return server.Router.Run(runPort)
}

func (server *Server) StatusHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "running",
		})
	}
}

func (server *Server) AuthorizeHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizeCode := context.Query("code")
		if len(strings.TrimSpace(authorizeCode)) == 0 {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		accessToken, refreshToken := server.OAuth2Token(authorizeCode)
		if len(accessToken) == 0 || len(refreshToken) == 0 {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		logrus.Info(fmt.Sprintf("access_token=%s refresh_token=%s", accessToken, refreshToken))
		snowflake, username := server.Profile(accessToken)
		if len(snowflake) == 0 || len(username) == 0 {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		authenticated := common.Authenticated{
			Snowflake:    snowflake,
			Username:     username,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		if err := common.Create(authenticated); err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		server.OnAuthenticate(authenticated)
		context.Redirect(http.StatusPermanentRedirect, "discord://")
	}
}
