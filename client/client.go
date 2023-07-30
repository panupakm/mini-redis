// client to connect to mini redis server
package client

import (
	"fmt"
	"net"
	"os"
)

func Start(network, host, port string) {
	//establish connection
	connection, err := net.Dial(network, host+":"+port)
	if err != nil {
		os.Exit(1)
	}
	defer connection.Close()

	///send some data
	_, err = connection.Write([]byte("Hello Server! Greetings."))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
}
