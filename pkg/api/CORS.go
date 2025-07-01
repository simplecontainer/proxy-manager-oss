package api

import "github.com/gin-gonic/gin"

func (api *Api) CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", api.Config.AllowOrigin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Upstream, Authorization, ")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
}
