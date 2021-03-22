package main

import (
	"net"
	"strings"
)

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

//SendMsg is an API for user to send a message to himself
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

//DoMessage is an API for users to process messages for business
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//Query the current online users
		u.server.maplock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "online...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.maplock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//The message format: rename|tom
		newName := strings.Split(msg, "|")[1]
		//Check if name exists
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("The current user name is used\n")
		} else {
			u.server.maplock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.maplock.Unlock()

			u.Name = newName
			u.SendMsg("You have updated the user name:" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//The message format:  to|Tom|The message content

		//1. Get the user name of the other party
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("The message format is not correct, please use the \" to|Tom|hello \" format.\n")
			return
		}

		//2. Get the User object according to the User name
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("The user name does not exist\n")
			return
		}

		//3. Get the message content and send it to the User object of the other party
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("No message content, please resend\n")
			return
		}
		remoteUser.SendMsg(u.Name + "Said to you:" + content)

	} else {
		u.server.Broadcast(u, msg)
	}

}

// ListenMessage is a method that listens for the current user channel and sends a message directly to the opposite client
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
