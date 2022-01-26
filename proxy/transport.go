package proxy

import "net"

type Transport interface {
	Network() string
	Addr() net.Addr
}

type udpTransport struct {
	conn *net.UDPConn
}

func (t *udpTransport) Network() string {
	return "UDP"
}

func (t *udpTransport) Addr() net.Addr {
	return t.conn.LocalAddr()
}

func newUDPTransport(conn *net.UDPConn) *udpTransport {
	return &udpTransport{conn: conn}
}
