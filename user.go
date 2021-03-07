package main

import "net"

// User is a structure containing name 、Addr 、C and conn
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser is an API to create a user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//Start a goroutine that listens for messages on the current User Channel
	go user.ListenMessage()
	return user
}

//OnLine is an API of User's online business
func (u *User) OnLine() {
	// Add the user to the onlineMap
	u.server.maplock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.maplock.Unlock()

	// Broadcast the current user online message
	u.server.Broadcast(u, "The user is online")
}

//OffLine is an API of User's offline business
func (u *User) OffLine() {
	// Remove user from onlineMap
	u.server.maplock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.maplock.Unlock()

	// Broadcast the current user offline message
	u.server.Broadcast(u, "The user has logged off")
}

//DoMessage is an API for users to process messages for business
func (u *User) DoMessage(msg string) {
	u.server.Broadcast(u, msg)
}

// ListenMessage is a method that listens for the current user channel and sends a message directly to the opposite client
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
