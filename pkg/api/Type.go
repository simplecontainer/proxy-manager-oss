package api

import (
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/configuration"
	"gitlab.com/simplecontainer/proxy-manager-oss/pkg/proxy"
)

type Api struct {
	Manager *proxy.Manager
	Config  *configuration.Configuration
}
