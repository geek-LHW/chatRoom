package main

import "net"

// User is a structure containing name 、Addr 、C and conn
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser is an API to create a user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	//Start a goroutine that listens for messages on the current User Channel
	go user.ListenMessage()
	return user
}

// ListenMessage is a method that listens for the current user channel and sends a message directly to the opposite client
func (s *User) ListenMessage() {
	for {
		msg := <-s.C
		s.conn.Write([]byte(msg + "\n"))
	}
}
