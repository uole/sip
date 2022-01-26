package proxy

import (
	"github.com/uole/sip"
	"net"
)

type (
	Conn interface {
		Addr() net.Addr
		Request(req *sip.Request) (err error)
		Response(res *sip.Response) (err error)
	}

	UdpConn struct {
		addr *net.UDPAddr
		conn *net.UDPConn
	}
)

func (conn *UdpConn) Addr() net.Addr {
	return conn.addr
}

func (conn *UdpConn) Request(req *sip.Request) (err error) {
	_, err = conn.conn.WriteToUDP(req.Bytes(), conn.addr)
	return
}

func (conn *UdpConn) Response(res *sip.Response) (err error) {
	_, err = conn.conn.WriteToUDP(res.Bytes(), conn.addr)
	return
}

func newUDPConn(addr string, conn *net.UDPConn) *UdpConn {
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	return &UdpConn{
		addr: udpAddr,
		conn: conn,
	}
}
