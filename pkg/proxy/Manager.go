package proxy

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/simplecontainer/proxy-manager-oss/pkg/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

func New() *Manager {
	return &Manager{
		Proxies: sync.Map{},
		mu:      sync.Mutex{},
	}
}

func (mgr *Manager) Add(URL string, p *Proxy) *Proxy {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	logger.Log.Info("added proxy", zap.Any("proxy", p))

	mgr.Proxies.Store(URL, p)
	return p
}
func (mgr *Manager) Find(URL string) (*Proxy, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	entry, ok := mgr.Proxies.Load(URL)

	if !ok {
		return nil, errors.New("not found")
	}

	return entry.(*Proxy), nil
}
func (mgr *Manager) Track(URL string, p *Proxy) error {
	logger.Log.Info("started tracking proxy", zap.Any("proxy", p))

	p.Timeout = time.AfterFunc(120*time.Second, func() {
		p.KeepAlive <- false
		err := mgr.Terminate(URL, p)

		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
	})

	go mgr.trackWebsocket(p)
	return nil
}
func (mgr *Manager) trackWebsocket(p *Proxy) {
	for {
		select {
		case keepAlive := <-p.KeepAlive:
			logger.Log.Info("keeping alive socket", zap.Any("proxy", p))

			if keepAlive {
				p.Timeout.Reset(120 * time.Second)
			} else {
				return
			}
			break
		}
	}
}

func (mgr *Manager) Terminate(URL string, p *Proxy) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	logger.Log.Info("terminating proxy", zap.String("UUID", URL))

	p.ActiveSockets.Range(func(key, value interface{}) bool {
		logger.Log.Info("killing user webocket", zap.Any("key", key))

		conn := value.(*websocket.Conn)
		conn.Close()

		return true
	})

	mgr.Proxies.Delete(URL)
	return nil
}
