package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Mux struct {
	connections map[*Conn]bool

	Broadcast chan Muxable

	register chan *Conn

	unregister chan *Conn
}

type Muxable interface {
	Marshal() *[]byte
}

func NewMux() *Mux {

	mux := Mux{connections: make(map[*Conn]bool),
		Broadcast:  make(chan Muxable),
		register:   make(chan *Conn),
		unregister: make(chan *Conn),
	}
	go mux.loop()

	return &mux
}

func (m *Mux) loop() {
	var conn *Conn
	var muxable Muxable
	var msg *[]byte
	for {
		select {
		case conn = <-m.register:
			//register new connection
			m.connections[conn] = true
			log.Printf("Client registered: %p, %d total.", conn, len(m.connections))

		case conn = <-m.unregister:
			//remove connection
			delete(m.connections, conn)
			close(conn.Output)
			log.Printf("Client unregistered: %p, %d total.", conn, len(m.connections))

		case muxable = <-m.Broadcast:
			msg = muxable.Marshal()
			for conn, _ := range m.connections {
				conn.Output <- *msg
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func (m *Mux) Handle(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Could not upgrade http request: %s", err.Error())
		return
	}

	conn := NewConn(m, ws)
	m.register <- conn
}
