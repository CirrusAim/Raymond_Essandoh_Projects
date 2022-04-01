package network

import (
	"context"
	"encoding/json"
	"lab3/failuredetector"
	"log"
	"net"
)

const BufferSize = 1024

type UDPServer struct {
	conn   *net.UDPConn
	fd     *failuredetector.EvtFailureDetector
	hbSend chan failuredetector.Heartbeat
}

type udpClientInfo struct {
	conn    *net.UDPConn
	address string
}
type UDPClient struct {
	idConnInfo map[int]*udpClientInfo
}

func (u *UDPServer) Close() {
	u.conn.Close()
}

// processRequest decode the incoming request by unmarshalling it and then forwarding it to failure detector
func (u *UDPServer) processRequest(addr *net.UDPAddr, req []byte) {
	var hb failuredetector.Heartbeat
	err := json.Unmarshal(req, &hb)
	if err != nil {
		log.Printf("Error unmarshalling data. Skipping. Err : %v", err)
	}
	u.fd.DeliverHeartbeat(hb)
}

func NewUDPServer(addr string, fd *failuredetector.EvtFailureDetector,
	hbSend chan failuredetector.Heartbeat) *UDPServer {

	conn, err := initServer(addr)
	if err != nil {
		log.Fatalf("Unable to start UDP server. Error: %v", err)
	}

	return &UDPServer{
		conn:   conn,
		fd:     fd,
		hbSend: hbSend,
	}

}

func initServer(addr string) (*net.UDPConn, error) {
	udpAddress, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		return nil, err
	}
	return udpConn, nil
}

// ServeUDP start listening to incoming requests
func (u *UDPServer) ServeUDP(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			u.Close()
			return
		default:
			var tempBuffer [BufferSize]byte
			n, addr, err := u.conn.ReadFromUDP(tempBuffer[0:])
			if err != nil {
				log.Printf("Error : %s", err)
				continue
			}
			go u.processRequest(addr, tempBuffer[:n])
		}
	}
}

func NewUDPClient(addersMap map[int]string) *UDPClient {

	client := &UDPClient{
		idConnInfo: make(map[int]*udpClientInfo),
	}

	for k, v := range addersMap {
		addr, err := net.ResolveUDPAddr("udp", v)
		if err != nil {
			log.Fatalf("Error  : %s", err)
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Fatalf("Error  : %s", err)
		}
		connInfo := &udpClientInfo{
			address: v,
			conn:    conn,
		}
		client.idConnInfo[k] = connInfo
	}

	return client
}

func (c *UDPClient) SendPayload(to int, payload []byte) error {
	udpInfo := c.idConnInfo[to]
	_, err := udpInfo.conn.Write(payload)
	if err != nil {
		log.Printf("Error  : %s", err)
	}

	return err
}

func (c *UDPClient) Close() {
	for _, v := range c.idConnInfo {
		v.conn.Close()
	}
}
