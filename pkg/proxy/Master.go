package proxy

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/simplecontainer/proxy-manager-oss/pkg/api/middlewares"
	"github.com/simplecontainer/proxy-manager-oss/pkg/configuration"
	"github.com/simplecontainer/proxy-manager-oss/pkg/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func StartMasterProxy(mgr *Manager, config *configuration.Configuration) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			middlewares.SetCORSHeaders(w, config.AllowOrigin)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		upstream := GetUpstream(r)

		proxy, err := mgr.Find(upstream)
		if err != nil {
			logger.Log.Error("failed to find proxy", zap.Error(err), zap.String("upstream", upstream))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.ToLower(r.Header.Get("Connection")) == "upgrade" && strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
			WebSocket(w, r, proxy)
			return
		}

		p := httputil.NewSingleHostReverseProxy(proxy.URL)
		p.Transport = &http.Transport{
			TLSClientConfig: proxy.TLSConfig,
		}

		p.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Log.Error("proxy error", zap.Error(err))

			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Bad Gateway"))
		}

		p.ModifyResponse = func(r *http.Response) error {
			middlewares.SetCORSHeaders(&headerWrapper{r.Header}, config.AllowOrigin)
			return nil
		}

		p.ServeHTTP(w, r)
	})

	if config.Environment == configuration.PRODUCTION_ENV {
		server := &http.Server{
			Addr:    fmt.Sprintf(":%s", config.MasterPort),
			Handler: handler,
		}

		log.Println(fmt.Sprintf("Master proxy production mode - listening on :%s", config.MasterPort))
		return server.ListenAndServe()
	} else {
		server := &http.Server{
			Addr:    fmt.Sprintf(":%s", config.MasterPort),
			Handler: handler,
		}

		log.Println(fmt.Sprintf("Master proxy development mode - listening on :%s", config.MasterPort))
		return server.ListenAndServeTLS("./app.simplecontainer.io.cert.pem", "./app.simplecontainer.io.pem")
	}
}

func WebSocket(w http.ResponseWriter, r *http.Request, p *Proxy) {
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{"Upstream"},
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	defer conn.Close()

	wssURL := fmt.Sprintf("wss://%s%s", strings.Replace(p.URL.String(), "https://", "", -1), r.URL.Path)

	requestHeader := http.Header{}
	for k, vs := range r.Header {
		// Filter out websocket-specific headers
		switch strings.ToLower(k) {
		case "connection", "upgrade", "sec-websocket-key",
			"sec-websocket-version", "sec-websocket-extensions", "sec-websocket-protocol":
			continue
		default:
			for _, v := range vs {
				requestHeader.Add(k, v)
			}
		}
	}

	dialer := websocket.Dialer{
		TLSClientConfig: p.TLSConfig,
		Subprotocols:    []string{"Authorization"},
	}

	serverConn, resp, err := dialer.Dial(wssURL, requestHeader)
	if err != nil {
		logger.Log.Error("Failed to connect to upstream WebSocket", zap.Error(err), zap.Any("response", resp))
		return
	}
	defer serverConn.Close()

	errorChan := make(chan error, 2)

	p.ActiveSockets.Store(conn.NetConn().LocalAddr().String(), conn)

	go func() {
		for {
			messageType, bytes, err := conn.ReadMessage()
			if err != nil {
				logger.Log.Error("Failed to read to client WebSocket", zap.Error(err), zap.Any("response", resp))
				errorChan <- err
				return
			}

			if err := serverConn.WriteMessage(messageType, bytes); err != nil {
				logger.Log.Error("Failed to write to client WebSocket", zap.Error(err), zap.Any("response", resp))
				errorChan <- err
				return
			}
		}
	}()

	go func() {
		for {
			messageType, bytes, err := serverConn.ReadMessage()
			if err != nil {
				logger.Log.Error("Failed to read to upstream WebSocket", zap.Error(err), zap.Any("response", resp))
				errorChan <- err
				return
			}

			if err := conn.WriteMessage(messageType, bytes); err != nil {
				logger.Log.Error("Failed to write to upstream WebSocket", zap.Error(err), zap.Any("response", resp))
				errorChan <- err
				return
			}
		}
	}()

	<-errorChan
}

type headerWrapper struct {
	header http.Header
}

func (h *headerWrapper) Header() http.Header {
	return h.header
}

func (h *headerWrapper) Write([]byte) (int, error) {
	return 0, nil
}

func (h *headerWrapper) WriteHeader(statusCode int) {}
