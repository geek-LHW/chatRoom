package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //The current pattern of the Client
}

func NewClient(serverIp string, serverPort int) *Client {
	// Create the client object
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	//Connected to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	//Returns the object
	return client
}

//Messages that the server responds to are displayed directly to standard output
func (client *Client) DealResponse() {
	//Once client. Conn has data, it copies it directly to stdout standard output, permanently blocking listening
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. The public chat mode")
	fmt.Println("2. Private conversation")
	fmt.Println("3. Update the user name")
	fmt.Println("0. Quit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>Please enter the number within the legal range<<<<")
		return false
	}

}

//Query Online Users
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

//Private chat mode
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>Please enter the chat object [user name] and exit :")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>Please enter the message content and exit to exit:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//The message is sent if it is not empty
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>Please enter the message content and exit to exit:")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>Please enter the chat object [user name] and exit:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	// Prompt the user for a message
	var chatMsg string

	fmt.Println(">>>>Please enter the chat content and exit.")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// Send to the server

		// Send the message if it is not empty
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>Please enter the chat content and exit.")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) UpdateName() bool {

	fmt.Println(">>>>Please enter user name:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		// Handle different business according to different patterns
		switch client.flag {
		case 1:
			// Public chat mode
			client.PublicChat()
			break
		case 2:
			// Private chat mode
			client.PrivateChat()
			break
		case 3:
			// Update the user name
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

//./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set server IP address (default is 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "Set the server port (8888 by default)")
}

func main() {
	//Command line parsing
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> Link server failed...")
		return
	}

	//Open a separate goroutine to handle the server's return receipt messages
	go client.DealResponse()

	fmt.Println(">>>>>Link to server successfully...")

	//Start the client's business
	client.Run()
}
