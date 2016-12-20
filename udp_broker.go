package glutton

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/hectane/go-nonblockingchan"
)

// Connection struct for UDP connection
type Connection struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	ch     *nbc.NonBlockingChan
	f      *os.File
	buffer [1500]byte
	n      int
}

// UDPBroker is handling and UDP connection
func UDPBroker(c *Connection) {
	log.SetOutput(c.f)
	tmp := c.addr.String()
	if tmp == "<nil>" {
		log.Println("Error. Address:port == nil udp_broker.go addr.String()")
		return
	}
	str := strings.Split(tmp, ":")
	dp := GetUDPDesPort(str, c.ch)
	if dp == -1 {
		log.Println("Warning. Packet dropped! [UDP] udp_broker.go desPort == -1")
		return
	}
	host := GetHandler(dp)
	if len(host) < 2 {
		log.Println("Error. [UDP] No host found. Packet dropped!")
		return
	}
	udpAddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		log.Println("Error. udp_broker() net.ResolveUDPAddr Could not resolve host address!")
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Println("Error. udp_broker() net.DialUDP Failed to connect to host!")
		return
	}

	_, err = conn.Write(c.buffer[0:c.n])
	if err != nil {
		log.Println("Error. udp_broker() conn.Write() Failed to write on connection!")
		return
	}

	log.Printf("[%v -> %v] Payload: %v", c.addr, udpAddr, string(c.buffer[0:c.n]))

	var buf [1500]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		log.Println("Warning. udp_broker() conn.Read() Failed to read from connection!")
		return
	}

	log.Printf("[%v <- %v] Payload: %v", c.addr, udpAddr, string(buf[0:n]))

	num, err := c.conn.WriteToUDP(buf[0:n], c.addr)
	if err != nil {
		log.Printf("Error. [%v] %v\n", num, err)
	}

}
