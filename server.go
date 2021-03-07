package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Server is a structure containing IP and ports
type Server struct {
	IP   string
	Port int
	//A list of online users
	OnlineMap map[string]*User
	maplock   sync.RWMutex

	//News broadcast
	Message chan string
}

// NewServer is An API to create a server
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessage is a method that listens on the Message broadcast message channel and sends a message to all online users as soon as it is available
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.maplock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.maplock.Unlock()
	}
}

// Broadcast is a method to broadcast a message
func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

// Handler is a method of dealing with business
func (s *Server) Handler(conn net.Conn) {
	// Currently linked business …
	// fmt.Println("链接建立成功")
	user := NewUser(conn, s)
	user.OnLine()

	// Receive the message sent by the client
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			// Extract user information and remove "\n"
			msg := string(buf[:n-1])
			// Broadcast the information received
			user.DoMessage(msg)
		}
	}()

	// HANDER is currently blocked
	select {}
}

//Start is an interface to start the server
func (s *Server) Start() {
	//socket listen
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net Listen err:", err)
		return
	}
	//close socket listen
	defer Listener.Close()

	//Start the goroutine that listens for messages

	go s.ListenMessage()
	for {
		//accept
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}

}
