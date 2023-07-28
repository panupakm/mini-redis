// mini-redis project server main.go
package main

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost'"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

func main() {
	server, err := net.Listen(SERVER_TYPE, fmt.Sprintf("%s:%s", SERVER_HOST, SERVER_PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Printf(("Listening on %s:%s\n"), SERVER_HOST, SERVER_PORT)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		if err != nil {
		}
		fmt.Println("client connected")
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	fmt.Printf("Received: ", string(buffer[:mLen]))
	_, err = connection.Write([]byte("Thanks! got your message: " + string(buffer[:mLen])))
	connection.Close()
}
