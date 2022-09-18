package main

import (
	"bufio"
	"fmt"
)

const noSupportedMethod byte = 255


// attempt authentication
func (s* Server) tryAuth(bufConn *bufio.ReadWriter) (bool, error) {
	//check if the client support the require auth method
	if b, e := s.checkAuthMethod(bufConn); !b || e != nil {
		return false, e
	}

	serverMethod := s.authCredential.GetServerMethod()
	if serverMethod == 0 {
		//if no authentication is require return true
		return true, nil
	} else {
		fmt.Println("Not implemented yet")
		return false, nil
	}
	
}




// select an implemented auth method from the ones the client offers
func (s* Server) checkAuthMethod(bufConn *bufio.ReadWriter) (bool, error) {
	var serverMethod = s.authCredential.GetServerMethod()


	// read the first byte and get number of availabe methods
	nAuth ,err := bufConn.ReadByte()
	if err != nil {
		return false, fmt.Errorf("error auth read: %v", err)
	}

	// read availabe methods in the slice methods
	methods := make([]byte, nAuth)
	_, err = bufConn.Read(methods)
	if err != nil {
		return false, fmt.Errorf("error auth method read: %v", err)
	}



	// check if ont the auth method is availabe
	// if it is select it
	var selectMethod = noSupportedMethod
	for _, v := range methods {
		if v == serverMethod { selectMethod = serverMethod }
	}



	// answer to the client with the server auth method
	msg := [2]byte{s.serverVersion, selectMethod}
	if _, err = bufConn.Write(msg[:]); err != nil {
		return false, fmt.Errorf("failed to anwser client: %v", err)
	}
	if err = bufConn.Writer.Flush(); err != nil {
		return false, fmt.Errorf("failed to anwser client: %v", err)
	}

	

	// return status
	if selectMethod == 255 {
		return false, nil
	}

	return true, nil
}