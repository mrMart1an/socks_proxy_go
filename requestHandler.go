package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// error constant
const (
	generalFailure byte = 1
)




func (r *Request) handleRequest(conn *bufio.ReadWriter) error {
	//get the address in string version
	destAddr, err := r.generateAddressPort()
	if err != nil {
		return fmt.Errorf("addr error: %v", err)
	}
	

	// start an handler based on the requested command
	switch r.cmd {
		// tcp/ip stream 
		case 1:
			return r.handleConnectTCP(conn, destAddr)

		// tcp/ip port binding
		case 2:
			return fmt.Errorf("not implemented")

		// udp port association
		case 3:
			return fmt.Errorf("not implemented")

		default:
			return fmt.Errorf("unrecognizd command form client")
	}
}


// handle connections -------------------------------------------------------

// connect to the target
func (r* Request) handleConnectTCP(clientConn *bufio.ReadWriter, destAddr string) error {
	// create the conncetion to the destination
	destConn, err := net.Dial("tcp", destAddr)
	if err == nil {
		defer destConn.Close()
	}

	// send client answer
	if err = r.sendRequestAnswer(clientConn, err); err != nil {
		return err
	}


	// create the tunnel chanel
	var ch = make(chan error)

	// strat the tunnels
	go r.directionalTunnel(clientConn, destConn, ch)
	go r.directionalTunnel(destConn, clientConn, ch)

	//wait for goroutine to return
	for i := 0; i < 2; i++ {
		// TODO error handling
		<- ch 
	}

	return nil
}




// handle connections -------------------------------------------------------

// answer to the client request
// answer structure:
// [0: version], [1: status], [2: resv], [3: ipType], [IP, PORT]
func (r *Request) sendRequestAnswer(clientConn *bufio.ReadWriter, connErr error) error {
	//create the buffer
	msg := make([]byte, 4, 6 + len(r.IP))
	var returnErr error

	msg[0] = r.serverVersion
	msg[3] = r.IPtype

	// get connection status 
	if connErr != nil {
		msg[1] = generalFailure
	}


	// address and port
	msg = append(msg, r.IP..., )
	msg = append(msg, r.PORT...)


	// send the message
	if _, returnErr = clientConn.Write(msg[:]); returnErr != nil {
		return fmt.Errorf("failed to anwser client: %v", returnErr)
	}
	if returnErr = clientConn.Writer.Flush(); returnErr != nil {
		return fmt.Errorf("failed to anwser client: %v", returnErr)
	} 
	

	// return the inputed conn error if present
	if connErr != nil {
		return connErr
	} 

	return returnErr
}




func (r *Request) directionalTunnel(conn1 io.Writer, conn2 io.Reader, ch chan error) {
	_, err := io.Copy(conn1, conn2)
	
	//return error to stop the tunneling 
	ch <- err
}