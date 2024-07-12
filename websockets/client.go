package websockets

import (
	"io"
	"log"

	"github.com/gorilla/websocket"
)


type Client struct {
  id        string
  hub       *Hub
  socket    *websocket.Conn
  outbound  chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
  return &Client{
    hub: hub, 
    socket: socket,
    outbound: make(chan []byte),
  }
}

func (c *Client) Write() {
  for {
    select {
    case message, ok := <-c.outbound:
      if !ok {
        c.socket.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      c.socket.WriteMessage(websocket.TextMessage, message)
    }
  } 
}

func (c *Client) ReadFromSocket() {
    for  {
        if _, reader, err  := c.socket.NextReader(); err != nil {
            log.Println("Closing socket with", c.id)
            c.hub.unregister<-c
            return
        } else {
            data, err := io.ReadAll(reader)
            if (err != nil) {
                log.Println("error:", err.Error())
                continue
            }
            log.Println(">", string(data))
        }
    }
}

