package sip

import (
	"context"
	"net"
)

const (
	ProtoUDP = "UDP"
	ProtoTCP = "TCP"
)

type (
	ProcessFunc func(res *Response) (handled bool, err error)

	Transport interface {
		Dial(addr string) (err error)
		Protocol() string
		Conn() net.Conn
		Request() chan *Request
		Do(ctx context.Context, req *Request, fun ProcessFunc) (err error)
		Write(p []byte) (n int, err error)
		Close() (err error)
	}
)
