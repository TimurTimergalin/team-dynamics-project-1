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

func Write[T ueJson.ResponsePayload](c *Connection, payload T) error {
	data, err := ueJson.SerializeResponse(payload)
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
	go c.readLoop()
	return c
}
