package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/simplecontainer/proxy-manager-oss/pkg/configuration"
	"net/http"
)

func SetCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie, Upstream, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func CORS(config *configuration.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetCORSHeaders(c.Writer, config.AllowOrigin)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
