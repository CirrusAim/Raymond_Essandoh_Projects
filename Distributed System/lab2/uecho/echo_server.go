// Leave an empty line above this comment.
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const BUFFER_SIZE = 1024

type UDPServer struct {
	conn *net.UDPConn
}

func (u *UDPServer) Close() {
	u.conn.Close()
}

func (u *UDPServer) processRequest(addr *net.UDPAddr, req []byte) {
	// function writes data back to client
	writeData := func(cmd, output string) {
		_, err := u.conn.WriteToUDP([]byte(output), addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "CMD : %s .Error : %s\n", cmd, err)
		}
	}

	reqStr := string(req)

	if strings.Count(reqStr, "|:|") > 1 {
		writeData("", "Unknown command")
		return
	}
	reqArr := strings.Split(reqStr, "|:|")
	cmd := reqArr[0]
	data := reqArr[1]

	var output string
	switch cmd {
	case "UPPER":
		output = strings.ToUpper(data)
	case "LOWER":
		output = strings.ToLower(data)
	case "CAMEL":
		output = strings.Title(strings.ToLower(data))
	case "SWAP":
		output = strings.Map(swapCase, data)
	case "ROT13":
		output = rot13(data)
	default:
		output = "Unknown command"
	}

	writeData(cmd, output)

}

// NewUDPServer returns a new UDPServer listening on addr. It should return an
// error if there was any problem resolving or listening on the provided addr.
func NewUDPServer(addr string) (*UDPServer, error) {
	udpAddress, err := net.ResolveUDPAddr("udp", addr)
	if isError(err) {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", udpAddress)
	if isError(err) {
		return nil, err
	}

	return &UDPServer{
		conn: udpConn,
	}, nil
}

// ServeUDP starts the UDP server's read loop. The server should read from its
// listening socket and handle incoming client requests as according to the
// the specification.
func (u *UDPServer) ServeUDP() {
	defer u.Close()

	for {
		var tempBuffer [BUFFER_SIZE]byte
		n, addr, err := u.conn.ReadFromUDP(tempBuffer[0:])
		if isError(err) {
			continue
		}
		go u.processRequest(addr, tempBuffer[:n])
	}
}

func isError(err error) bool {
	if err != nil {
		if socketIsClosed(err) {
			fmt.Fprintf(os.Stderr, "Socket closed. Reinitiate connection")
		} else {
			fmt.Fprintf(os.Stderr, "Error : %s\n", err)
		}
		return true
	}
	return false
}

// socketIsClosed is a helper method to check if a listening socket has been
// closed.
func socketIsClosed(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

func swapCase(ch rune) rune {
	switch {
	case 'a' <= ch && ch <= 'z':
		return ch - 'a' + 'A'
	case 'A' <= ch && ch <= 'Z':
		return ch - 'A' + 'a'
	default:
		return ch
	}
}

// function retun rot13 encrypted data
func rot13(data string) string {
	p := []byte(data)
	for i := 0; i < len(p); i++ {
		if (p[i] >= 'A' && p[i] <= 'M') || (p[i] >= 'a' && p[i] <= 'm') {
			p[i] += 13
		} else if (p[i] >= 'N' && p[i] <= 'Z') || (p[i] >= 'n' && p[i] <= 'z') {
			p[i] -= 13
		}
	}
	return string(p[:])
}
