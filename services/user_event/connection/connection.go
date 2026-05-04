package connection

import (
	"github.com/gorilla/websocket"
	ueJson "team_dynamics/user_event/json"
	"time"
)

type ReceivedMessage struct {
	Request *ueJson.Request
	Err     error
}

type Connection struct {
	conn    *websocket.Conn
	msgChan chan ReceivedMessage
}

func (c *Connection) Messages() <-chan ReceivedMessage {
	return c.msgChan
}

func (c *Connection) Write(response *ueJson.Response) error {
	data, err := ueJson.SerializeResponse(response)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Connection) readLoop() {
	defer close(c.msgChan)
	for {
		_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.msgChan <- ReceivedMessage{Err: err}
			break
		}

		req, err := ueJson.ParseRequest(data)
		if err != nil {
			_ = c.conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error()),
				time.Now().Add(5*time.Second),
			)
			c.msgChan <- ReceivedMessage{Err: err}
			break
		}

		c.msgChan <- ReceivedMessage{Request: req}
	}
}

func (c *Connection) pingLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := c.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(3*time.Second))
		if err != nil {
			return
		}
	}
}

func (c *Connection) Close() error {
	_ = c.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(5*time.Second),
	)
	return c.conn.Close()
}

func WrapConnection(conn *websocket.Conn) *Connection {
	c := &Connection{
		conn:    conn,
		msgChan: make(chan ReceivedMessage),
	}
	c.conn.SetPongHandler(func(appData string) error {
		return c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
	c.conn.SetPingHandler(func(appData string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return c.conn.WriteControl(websocket.PongMessage, nil, time.Time{})
	})
	go c.readLoop()
	go c.pingLoop()
	return c
}
