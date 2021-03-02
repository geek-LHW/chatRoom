package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Hander(conn net.Conn) {
	// Currently linked business …
	fmt.Println("链接建立成功")

}

//An interface to start the server
func (this *Server) start() {
	//socket listen
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net Listen err:", err)
		return
	}
	//close socket listen
	defer Listener.Close()
	for {
		//accept
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		//do hander
		go this.Hander(conn)
	}

}
