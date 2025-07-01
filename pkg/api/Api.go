package api

import (
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/configuration"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/proxy"
)

func New(config *configuration.Configuration) *Api {
	return &Api{
		Manager: proxy.New(),
		Config:  config,
	}
}
