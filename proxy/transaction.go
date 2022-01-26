package proxy

import (
	"github.com/uole/sip"
	"net"
)

type (
	rewriter struct {
		From string
		To   string
	}

	Transaction struct {
		address   string
		process   *Process
		message   *Message
		transport Transport
	}
)

func (t *Transaction) ID() string {
	return t.message.CallID()
}

func (t *Transaction) Address() string {
	return t.address
}

func (t *Transaction) Caller() Conn {
	return t.process.Caller()
}

func (t *Transaction) Callee() Conn {
	return t.process.Callee()
}

func (t *Transaction) Request() *sip.Request {
	return t.message.Request()
}

func (t *Transaction) Response() *sip.Response {
	return t.message.Response()
}

func (t *Transaction) Rewrite() (*rewriter, bool) {
	if t.process.route != nil {
		if t.process.route.RewriteTo != "" {
			return &rewriter{
				From: t.process.route.Domain,
				To:   t.process.route.RewriteTo,
			}, true
		}
	}
	if t.process.relationship != nil {
		if t.process.relationship.OriginalDomain != t.process.relationship.Domain {
			return &rewriter{
				From: t.process.relationship.OriginalDomain,
				To:   t.process.relationship.Domain,
			}, true
		}
	}
	return nil, false
}

func newTransaction(msg *Message, proc *Process, source net.Addr, transport Transport) *Transaction {
	return &Transaction{
		process:   proc,
		message:   msg,
		address:   source.String(),
		transport: transport,
	}
}
