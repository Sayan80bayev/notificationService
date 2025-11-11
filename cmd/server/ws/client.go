package ws

import (
	"encoding/json"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	SendCh chan interface{}
	log    *logrus.Logger
}

func NewClient(userID uuid.UUID, conn *websocket.Conn) *Client {
	client := &Client{
		UserID: userID,
		Conn:   conn,
		SendCh: make(chan interface{}),
		log:    logging.GetLogger(),
	}
	return client
}

func (c *Client) ReadPump() {
	defer func() {
		Unregister(c.UserID)
		err := c.Conn.Close()
		if err != nil {
			c.log.Errorf("user %s disconnected: %v", c.UserID, err)
			return
		}
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			c.log.Errorf("user %s disconnected: %v", c.UserID, err)
			break
		}
	}
}

func (c *Client) WritePump() {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			c.log.Errorf("user %s disconnected: %v", c.UserID, err)
		}
	}(c.Conn)
	for msg := range c.SendCh {
		err := c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			c.log.Errorf("user %s disconnected: %v", c.UserID, err)
			return
		}
		data, err := json.Marshal(msg)
		if err != nil {
			c.log.Errorf("error marshalling message: %v", err)
			continue
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			c.log.Errorf("error writing websocket message: %v", err)
			break
		}
	}
}

func (c *Client) Send(msg interface{}) {
	select {
	case c.SendCh <- msg:
	default:
		c.log.Warnf("client %s send channel full, dropping message", c.UserID)
	}
}
