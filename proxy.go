package main

import (
	"fmt"
	"net"
)

func main() {
	// create the credential struct
	// TODO read json for credential
	var ac = AuthCredential{}

	// Get the server
	socksServer := GetServer(&ac)


	// listen of the specified interface and port
	// TODO read config file for port and interface
	l, err := net.Listen("tcp", "[127.0.0.1]:8000")
	if err != nil {
		fmt.Println("error during socket setup: ", err)
	}


	// handle connection
	fmt.Printf("Server started, waiting for connections...")
	for {
		conn, _ := l.Accept()

		// start a goroutine for each connection
		go socksServer.serveConn(conn)
	}
}
