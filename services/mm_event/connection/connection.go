package connection

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	mmeJson "team_dynamics/mm_event/json"
	"time"
)

type ReceivedMessage struct {
	Request *mmeJson.Request
	Err     error
}

type Connection struct {
	conn    *websocket.Conn
	msgChan chan ReceivedMessage
}

func (c *Connection) Write(resp *mmeJson.Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Connection) Messages() <-chan ReceivedMessage {
	return c.msgChan
}

func ptr[T any](v T) *T {
	return &v
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

		var req mmeJson.Request
		if err := json.Unmarshal(data, &req); err != nil {
			errMsg := err.Error()
			resp := &mmeJson.Response{
				Type:         mmeJson.Error,
				ErrorMessage: &errMsg,
			}
			_ = c.Write(resp)
			continue
		}

		if req.Type != mmeJson.Register && req.Type != mmeJson.Unregister {
			_ = c.Write(&mmeJson.Response{
				Type:         mmeJson.Error,
				ErrorMessage: ptr(fmt.Sprintf("Unknown type: %s", req.Type)),
			})
		}

		c.msgChan <- ReceivedMessage{Request: &req}
	}
}

func (c *Connection) Close() error {
	// Send close frame with normal closure code
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
