package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Conn struct {
	Output chan []byte
	socket *websocket.Conn
	mux    *Mux
}

func NewConn(m *Mux, s *websocket.Conn) *Conn {

	conn := Conn{Output: make(chan []byte), socket: s, mux: m}

	go conn.read()
	go conn.write()
	return &conn
}

func (c *Conn) write() {
	for {
		msg, ok := <-c.Output

		if !ok { //the channel is closed
			return
		}

		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error writing to %p: %s", c, err.Error())
			return //is this enough ?
		}
	}
}

func (c *Conn) read() {
	defer func() {
		c.mux.unregister <- c
		c.socket.Close()
	}()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %p: %s", c, err.Error())
			break
		}
		log.Printf("Message from %p: %s", c, msg)
	}
}
