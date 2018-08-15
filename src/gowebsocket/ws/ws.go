package ws

import (
	"github.com/gorilla/websocket"
	"sync"
	"errors"
)

type Connection struct {
	Conn       *websocket.Conn
	MessageIn  chan []byte
	MessageOut chan []byte
	CloseChan  chan byte
	mutex      sync.Mutex
	isClosed   bool
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection) {
	conn = &Connection{
		Conn:       wsConn,
		MessageIn:  make(chan []byte, 1000),
		MessageOut: make(chan []byte, 1000),
		CloseChan:  make(chan byte, 1),
	}

	go conn.ReadLoop()

	go conn.WriteLoop()

	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.MessageIn:
	case <-conn.CloseChan:
		err = errors.New("connection is closed")
	}

	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.MessageOut <- data:
	case <-conn.CloseChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) Close() {
	conn.Conn.Close()

	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.CloseChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

func (conn *Connection) ReadLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.Conn.ReadMessage(); err != nil {
			goto ERR
		}
		select {
		case conn.MessageIn <- data:
		case <-conn.CloseChan:
			goto ERR
		}

	}

ERR:
	conn.Close()
}

func (conn *Connection) WriteLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.MessageOut:
		case <-conn.CloseChan:
			goto ERR
		}

		if err = conn.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}
