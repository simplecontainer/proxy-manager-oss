package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/simplecontainer/proxy-manager-oss/pkg/configuration"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{
		"Upstream",
	},
}

func NewProxy() *Proxy {
	return &Proxy{
		Timeout:       nil,
		ActiveSockets: &sync.Map{},
		Port:          0,
	}
}

func (proxy *Proxy) Server(config *configuration.Configuration, target string, ca string, key string, certificate string) error {
	URL, err := url.Parse(target)

	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	if certificate != "" && key != "" {
		clientCert, err := tls.X509KeyPair([]byte(certificate), []byte(key))

		if err != nil {
			return err
		}

		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	proxy.TLSConfig = tlsConfig
	proxy.TLSTransport = transport
	proxy.URL = URL

	return nil
}
func (proxy *Proxy) WebSocket(w http.ResponseWriter, r *http.Request, URL *url.URL) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	defer conn.Close()

	wsURL := fmt.Sprintf("wss://%s", URL)

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	serverConn, _, err := dialer.Dial(wsURL, r.Header)
	if err != nil {
		log.Printf("Failed to connect to upstream WebSocket: %v", err)
		return
	}
	defer serverConn.Close()

	errorChan := make(chan error, 2)

	go func() {
		proxy.ActiveSockets.Store(conn.NetConn().LocalAddr().String(), conn)

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				errorChan <- err
				return
			}

			if err := serverConn.WriteMessage(messageType, p); err != nil {
				errorChan <- err
				return
			}

			proxy.KeepAlive <- true
		}
	}()

	go func() {
		for {
			messageType, p, err := serverConn.ReadMessage()
			if err != nil {
				errorChan <- err
				return
			}

			if err := conn.WriteMessage(messageType, p); err != nil {
				errorChan <- err
				return
			}

			proxy.KeepAlive <- true
		}
	}()

	proxy.ActiveSockets.Delete(conn.NetConn().LocalAddr().String())
	<-errorChan
}
func GetUpstream(r *http.Request) string {
	upstream := r.Header.Get("Upstream")
	if upstream != "" {
		return upstream
	}

	wsProtocol := r.Header.Get("Sec-WebSocket-Protocol")
	if wsProtocol != "" && strings.HasPrefix(wsProtocol, "Upstream, ") {
		return strings.TrimPrefix(wsProtocol, "Upstream, ")
	}

	return ""
}
