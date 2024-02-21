package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoRouteMiddleware(context *gin.Context) {
	context.AbortWithStatus(http.StatusNotFound)
}
