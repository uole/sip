package proxy

import (
	"context"
	"github.com/uole/sip"
)

const (
	DirectionRequest  = 0x01
	DirectionResponse = 0x02
)

type Message struct {
	context   context.Context
	direction int
	request   *sip.Request
	response  *sip.Response
}

func (ctx *Message) CallID() string {
	if ctx.direction == DirectionRequest {
		return ctx.Request().CallID()
	} else {
		return ctx.Response().CallID()
	}
}

func (ctx *Message) Context() context.Context {
	return ctx.context
}

func (ctx *Message) Direction() int {
	return ctx.direction
}

func (ctx *Message) Request() *sip.Request {
	return ctx.request
}

func (ctx *Message) Response() *sip.Response {
	return ctx.response
}
