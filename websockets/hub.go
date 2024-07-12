package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
  id string
  clients     []*Client
  unregister  chan *Client
  register    chan *Client
  mutex       *sync.Mutex
}

func NewHub() *Hub {
  hub := Hub{
    clients: make([]*Client, 0),
    register: make(chan *Client),
    unregister: make(chan *Client),
    mutex: &sync.Mutex{},
  }
  return &hub
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
  log.Println("Handling web socket!")
  socket, err := upgrader.Upgrade(w, r, nil)  
  if err != nil {
    log.Println(err.Error())
    http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
    return
  }

  client := NewClient(hub, socket) 
  hub.register <- client

  go client.Write()
  go client.ReadFromSocket()
}

func (hub *Hub) Run () {
  for {
    select {
    case  client := <-hub.register:
      hub.onConnect(client)
    case client := <- hub.unregister:
      hub.onDisconnect(client)
    }
  } 
}

func (hub *Hub) onConnect(client *Client) {
  remoteAddr := client.socket.RemoteAddr()
  log.Println("Client Connected", remoteAddr)  
  hub.mutex.Lock()
  defer hub.mutex.Unlock()

  client.id = remoteAddr.String()
  hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
  remoteAddr := client.socket.RemoteAddr()
  log.Println("Client Disconnected", remoteAddr)  
  client.socket.Close()

  hub.mutex.Lock() 
  defer hub.mutex.Unlock()
  
  indexToRemove := -1 
  for i, cli := range hub.clients {
    if cli.id == client.id {
      indexToRemove = i
    }
  }
    
  copy(hub.clients[indexToRemove:], hub.clients[indexToRemove+1:])
  hub.clients[len(hub.clients)-1] = nil
  hub.clients = hub.clients[:len(hub.clients)-1] 
}

func (hub *Hub) Broadcast(message interface{}, ignore *Client) { 
  data,_ := json.Marshal(message)
  
  for _, client := range hub.clients {
    if client != ignore {
      client.outbound<- data 
    }
  }
}
