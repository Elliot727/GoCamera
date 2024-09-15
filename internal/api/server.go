package api

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
)

type WebSocket struct {
	conn net.Conn
}

func NewWebSocket(conn net.Conn) *WebSocket {
	return &WebSocket{conn: conn}
}

func computeAcceptKey(key string) string {
	guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(key + guid))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (ws *WebSocket) HandleHandshake(req *http.Request) error {
	if req.Header.Get("Upgrade") != "websocket" || req.Header.Get("Connection") != "Upgrade" {
		return fmt.Errorf("not a WebSocket request")
	}

	key := req.Header.Get("Sec-WebSocket-Key")
	acceptKey := computeAcceptKey(key)

	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n\r\n", acceptKey)

	_, err := ws.conn.Write([]byte(response))
	if err != nil {
		return fmt.Errorf("failed to write handshake response: %v", err)
	}

	return nil
}
