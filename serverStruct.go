package main

import (
	"bufio"
	"fmt"
	"net"
)




type Server struct {
	//auth credential
	authCredential AuthCredential

	serverVersion byte
}


// return a instance of a sock server
func GetServer(authC *AuthCredential) *Server {
	server := Server{
		authCredential: *authC,
	}

	// set the server version
	server.serverVersion = 5

	return &server
}



func (s *Server) serveConn(conn net.Conn) {
	// close connection nicely in case of failure
	defer conn.Close()

	// get a bufio readWtriter
	var bufConn = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// read greeting message --------------------
	// check version number
	vers ,err := bufConn.ReadByte()
	if err != nil { 
		fmt.Print("conn read failed")
		return
	}

	if vers != 5 {
		fmt.Print("unsupported socks version")
		return
	}


	// client try to authenticate
	if status, err := s.tryAuth(bufConn); !status || err != nil {
		fmt.Println("authentication failed\nError:", err)
		return
	}

	
	//proccess the request
	var request = Request{}
	if e := request.DecodeRequest(bufConn); e != nil {
		fmt.Println(e)
		return
	}


	//handle request
	if e := request.handleRequest(bufConn); e != nil {
		fmt.Println(e)
	}
}