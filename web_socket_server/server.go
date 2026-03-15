package web_socket_server

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

package websocketfw

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	Path            string
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	SendBufferSize  int
}

func DefaultConfig() *Config {
	return &Config{
		Path:            "/ws",
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingPeriod:      (60 * time.Second * 9) / 10,
		MaxMessageSize:  4096,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		SendBufferSize:  256,
	}
}

type OnMessageFunc func(client *Client, msgType int, data []byte)

type UserInfo struct {
	UserID uint64
}

type Client struct {
	conn     *websocket.Conn
	send     chan outgoingMessage
	userInfo *UserInfo
}

type outgoingMessage struct {
	msgType int
	data    []byte
}

type Hub struct {
	identified   map[uint64]*Client
	unidentified map[*Client]bool
	register     chan *Client
	unregister   chan *Client
	broadcast    chan *broadcastMsg
	mu           sync.RWMutex
}

type broadcastMsg struct {
	msgType int
	data    []byte
}

func newHub() *Hub {
	return &Hub{
		identified:   make(map[uint64]*Client),
		unidentified: make(map[*Client]bool),
		register:     make(chan *Client, 100),
		unregister:   make(chan *Client, 100),
		broadcast:    make(chan *broadcastMsg, 100),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if client.userInfo != nil {
				h.identified[client.userInfo.UserID] = client
			} else {
				h.unidentified[client] = true
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if client.userInfo != nil {
				delete(h.identified, client.userInfo.UserID)
			} else {
				delete(h.unidentified, client)
			}
			h.mu.Unlock()
			close(client.send)

		case msg := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.identified {
				select {
				case client.send <- outgoingMessage{msgType: msg.msgType, data: msg.data}:
				default:
					go h.closeClient(client)
				}
			}
			for client := range h.unidentified {
				select {
				case client.send <- outgoingMessage{msgType: msg.msgType, data: msg.data}:
				default:
					go h.closeClient(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) closeClient(client *Client) {
	select {
	case h.unregister <- client:
	default:
		client.conn.Close()
	}
}

type WebSocketServer struct {
	hub        *Hub
	onMessage  OnMessageFunc
	config     *Config
	upgrader   websocket.Upgrader
	bufferPool *sync.Pool
}

func NewWebSocketServer(config *Config, onMessage OnMessageFunc) *WebSocketServer {
	if config == nil {
		config = DefaultConfig()
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  config.ReadBufferSize,
		WriteBufferSize: config.WriteBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	bufferPool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, config.MaxMessageSize)
		},
	}
	return &WebSocketServer{
		hub:        newHub(),
		onMessage:  onMessage,
		config:     config,
		upgrader:   upgrader,
		bufferPool: bufferPool,
	}
}

func (s *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{
		conn: conn,
		send: make(chan outgoingMessage, s.config.SendBufferSize),
	}
	s.hub.register <- client
	go s.readPump(client)
	go s.writePump(client)
}

func (s *WebSocketServer) readPump(client *Client) {
	defer func() {
		s.hub.unregister <- client
		client.conn.Close()
	}()
	client.conn.SetReadLimit(s.config.MaxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(s.config.PongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(s.config.PongWait))
		return nil
	})
	for {
		msgType, reader, err := client.conn.NextReader()
		if err != nil {
			break
		}
		buf := s.bufferPool.Get().([]byte)
		buf = buf[:0]
		for {
			if len(buf) == cap(buf) {
				break
			}
			remaining := buf[len(buf):cap(buf)]
			n, err := reader.Read(remaining)
			buf = buf[:len(buf)+n]
			if err == io.EOF {
				break
			}
			if err != nil {
				s.bufferPool.Put(buf[:0])
				return
			}
		}
		s.onMessage(client, msgType, buf)
		s.bufferPool.Put(buf[:0])
	}
}

func (s *WebSocketServer) writePump(client *Client) {
	ticker := time.NewTicker(s.config.PingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.conn.SetWriteDeadline(time.Now().Add(s.config.WriteWait))
			if err := client.conn.WriteMessage(msg.msgType, msg.data); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(s.config.WriteWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *WebSocketServer) IdentifyClient(client *Client, userID uint64, info *UserInfo) error {
	if client.userInfo != nil {
		return errors.New("client already identified")
	}
	s.hub.mu.Lock()
	defer s.hub.mu.Unlock()
	if _, ok := s.hub.unidentified[client]; !ok {
		return errors.New("client not in hub")
	}
	if oldClient, exists := s.hub.identified[userID]; exists {
		delete(s.hub.identified, userID)
		s.hub.closeClient(oldClient)
	}
	delete(s.hub.unidentified, client)
	info.UserID = userID
	client.userInfo = info
	s.hub.identified[userID] = client
	return nil
}

func (s *WebSocketServer) SendTo(userID uint64, msgType int, data []byte) {
	s.hub.mu.RLock()
	client, ok := s.hub.identified[userID]
	s.hub.mu.RUnlock()
	if !ok {
		return
	}
	select {
	case client.send <- outgoingMessage{msgType: msgType, data: data}:
	default:
		s.hub.closeClient(client)
	}
}

func (s *WebSocketServer) SendAll(msgType int, data []byte) {
	s.hub.broadcast <- &broadcastMsg{msgType: msgType, data: data}
}

func (s *WebSocketServer) Close(userID uint64) {
	s.hub.mu.RLock()
	client, ok := s.hub.identified[userID]
	s.hub.mu.RUnlock()
	if ok {
		s.hub.closeClient(client)
	}
}

func (s *WebSocketServer) CloseAll() {
	s.hub.mu.RLock()
	all := make([]*Client, 0, len(s.hub.identified)+len(s.hub.unidentified))
	for _, client := range s.hub.identified {
		all = append(all, client)
	}
	for client := range s.hub.unidentified {
		all = append(all, client)
	}
	s.hub.mu.RUnlock()
	for _, client := range all {
		s.hub.closeClient(client)
	}
}

func (s *WebSocketServer) Run() {
	go s.hub.Run()
	http.Handle(s.config.Path, s)
}
