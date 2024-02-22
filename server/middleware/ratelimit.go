package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

var limiter = rate.NewLimiter(rate.Every(time.Second*10), 3)

func RateLimitMiddleware(context *gin.Context) {
	if limiter.Allow() == false {
		context.AbortWithStatus(http.StatusTooManyRequests)
	}
}
