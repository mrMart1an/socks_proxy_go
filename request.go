package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)


const (
	ip4Conn byte = 1
	ip6Conn byte = 4
	DnsConn byte = 3
)

type Request struct {
	IP net.IP
	PORT []byte
	IPtype byte

	cmd byte
	serverVersion byte
}


// decode the request
// request buffer structure:
//[0: VER] - [1: CMD] - [2: RSV] - [3: ADDR TYPE]
func (r *Request) DecodeRequest(bufConn *bufio.ReadWriter) error {
	// read the request buffer
	bufCmd := make([]byte, 4)
	if _, e := bufConn.Read(bufCmd); e != nil {
		return fmt.Errorf("cmd read failed: %v", e)
	}


	// get request command and version
	r.serverVersion = bufCmd[0]
	r.cmd = bufCmd[1]
	

	// get dest addr length
	var addrLen byte
	switch bufCmd[3] {
		//ip type ipV4
		case ip4Conn:
			addrLen = 4

		//ip type ipV6
		case ip6Conn:
			addrLen = 16

		//if the client is requesting a dns name 
		case DnsConn:
			// read the domain name length
			len, e := bufConn.ReadByte()
			if e != nil {
				return fmt.Errorf("domain len failed")
			}

			addrLen = len


		//if the addess type is not recognize return an error
		default:
			return fmt.Errorf("unsupported address type")
	}



	//get dest addr ip
	//create a buffer to contain the address and port bytes
	bufAddr := make([]byte, addrLen + 2)
	if _, e := bufConn.Read(bufAddr); e != nil {
		return fmt.Errorf("failed reading address: %v", e)
	}


	//get the ip for the connection
	if bufCmd[3] == DnsConn {
		if e := r.domainLookup(bufAddr[:addrLen]); e != nil { return e }

	} else {
		r.IP = bufAddr[:addrLen]
		r.IPtype = bufCmd[3]
	}


	//get port for the connection
	r.PORT = bufAddr[addrLen:]
	if len(r.PORT) != 2 {
		return fmt.Errorf("invalid port")
	}


	//if everything is ok return nil
	return nil
}




// lookup the dns request domain name
func (r *Request) domainLookup(domain []byte) error {
	// lookup the domain and get a list of ip addresses
	ip_addresses, err := net.LookupIP(string(domain))
	if err != nil { return err }

	if len(ip_addresses) > 0 {
		r.IP = ip_addresses[0]
	} else {
		return fmt.Errorf("no ip addresses found")
	}

	//determine the ip type
	switch len(r.IP) {
		case 4:
			r.IPtype = 1

		case 16:
			r.IPtype = 4

		default:
			return fmt.Errorf("invalid ip found")
	}

	return nil
}




// generate a string for the socket.connect function
func (r *Request) generateAddressPort() (string, error) {
	// transform the port byte in a string
	port := fmt.Sprint(binary.BigEndian.Uint16(r.PORT));

	
	// transform the ip byte in a valid string
	if len(r.IP) == 4 || len(r.IP) == 16 {
		return fmt.Sprintf("[%v]:%v", r.IP.String(), port), nil
	} else {
		return "", fmt.Errorf("invalid ip address")
	}
}