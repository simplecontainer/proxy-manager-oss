package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/logger"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/proxy"
	"go.uber.org/zap"
	"net/http"
)

func (api *Api) AddOrGetProxy(c *gin.Context) {
	p := proxy.NewProxy()
	err := p.Server(api.Config, c.PostForm("url"), c.PostForm("ca"), c.PostForm("key"), c.PostForm("cert"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		signature := base64.RawStdEncoding.EncodeToString([]byte(c.PostForm("url")))
		p = api.Manager.Add(signature, p)

		go func() {
			err = api.Manager.Track(signature, p)

			if err != nil {
				logger.Log.Error("Proxy server failed", zap.Error(err))
				return
			}

			logger.Log.Info("proxy added", zap.String("signature", signature))
		}()

		c.JSON(http.StatusOK, p)
	}
}
