package api

import (
	"github.com/simplecontainer/proxy-manager-oss/pkg/configuration"
	"github.com/simplecontainer/proxy-manager-oss/pkg/proxy"
)

type Api struct {
	Manager *proxy.Manager
	Config  *configuration.Configuration
}
