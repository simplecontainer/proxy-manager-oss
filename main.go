package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/api"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/api/middlewares"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/configuration"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/logger"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logger.Log = logger.NewLogger(logLevel, []string{"stdout"}, []string{"stderr"})
	fmt.Println(fmt.Sprintf("logging level set to %s (override with LOG_LEVEL env variable)", logLevel))

	config := configuration.New()

	api := api.New(config)

	go func() {
		err := proxy.StartMasterProxy(api.Manager, api.Config)

		if err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	r.Use(middlewares.CORS(api.Config))

	r.OPTIONS("/proxy", api.CORS)
	r.POST("/proxy", api.AddOrGetProxy)

	log.Println(fmt.Sprintf("Listening on :%s", config.Port))

	if config.Environment == configuration.PRODUCTION_ENV {
		err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)

		if err != nil {
			panic(err)
		}
	} else {
		err := http.ListenAndServeTLS(fmt.Sprintf(":%s", config.Port), "./app.simplecontainer.io.cert.pem", "./app.simplecontainer.io.pem", r)

		if err != nil {
			panic(err)
		}
	}
}
