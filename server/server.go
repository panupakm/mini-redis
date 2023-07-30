// mini-redis project server main.go
package server

import (
	"fmt"
	"net"
	"os"
)

func Start(network, host, port string) {
	server, err := net.Listen(network, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Printf(("Listening on %s:%s\n"), host, port)
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

	fmt.Printf("Received: %s", string(buffer[:mLen]))
	_, err = connection.Write([]byte("Thanks! got your message: " + string(buffer[:mLen])))
	connection.Close()
}
