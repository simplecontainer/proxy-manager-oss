package proxy

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Proxy struct {
	Timeout       *time.Timer `json:"-"`
	ActiveSockets *sync.Map   `json:"-"`
	Port          int
	TLSConfig     *tls.Config     `json:"-"`
	TLSTransport  *http.Transport `json:"-"`
	KeepAlive     chan bool       `json:"-"`
	URL           *url.URL
}

type Manager struct {
	Proxies sync.Map
	mu      sync.Mutex
}
