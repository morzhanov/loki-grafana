package rest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type rest struct {
	log *zap.Logger
	v   string
}

type REST interface {
	Listen()
}

func (r *rest) handleVersion(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("version: %s", r.v))
}

func (r *rest) handleHello(c *gin.Context) {
	c.String(http.StatusOK, "Hi!")
}

func getDurationInMilliseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}

func getClientIP(c *gin.Context) string {
	requester := c.Request.Header.Get("X-Forwarded-For")
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}
	return requester
}

func (r *rest) jsonLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := getDurationInMilliseconds(start)
		r.log.Info(
			"request handled",
			zap.String("client_ip", getClientIP(c)),
			zap.Float64("duration", duration),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.RequestURI),
			zap.Int("status", c.Writer.Status()),
			zap.String("referrer", c.Request.Referer()),
		)
	}
}

func (r *rest) Listen() {
	router := gin.Default()
	router.Use(r.jsonLogMiddleware())
	router.GET("/version", r.handleVersion)
	router.GET("/hello", r.handleHello)
	if err := router.Run(); err != nil {
		r.log.Fatal("error during rest controller execution", zap.Error(err))
	}
}

func NewREST(l *zap.Logger, version string) REST {
	return &rest{l, version}
}
