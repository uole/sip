package sip

import (
	"strconv"
)

type Context struct {
	req  *Request
	sess *Session
}

//CallID 返回当前的ID
func (ctx *Context) CallID() string {
	return ctx.req.Header.Get(HeaderCallID).String()
}

//Request request 对象
func (ctx *Context) Request() *Request {
	return ctx.req
}

//SipFrom 获取
func (ctx *Context) SipFrom() *AddressHeader {
	if !ctx.Request().Header.Has(HeaderFrom) {
		return nil
	}
	return ctx.Request().Header.Get(HeaderFrom).(*AddressHeader)
}

//SipTo 获取
func (ctx *Context) SipTo() *AddressHeader {
	if !ctx.Request().Header.Has(HeaderTo) {
		return nil
	}
	return ctx.Request().Header.Get(HeaderTo).(*AddressHeader)
}

//write 返回一个Response对象
func (ctx *Context) Write(res *Response) (err error) {
	if !res.Header.Has(HeaderVia) {
		viaHead := ctx.req.Header.Get(HeaderVia).(*ViaHeader)
		if viaHead.Uri.Port == 0 {
			viaHead.Uri.Port = 5060
		}
		resViaHead := &ViaHeader{
			Protocol:        viaHead.Protocol,
			ProtocolVersion: viaHead.ProtocolVersion,
			Transport:       ctx.sess.transport.Protocol(),
			Uri:             viaHead.Uri.Clone(),
		}
		if viaHead.Uri.Port > 0 {
			resViaHead.Uri.Params.Set("rport", strconv.Itoa(viaHead.Uri.Port))
		}
		if viaHead.Uri.Params.Get("branch") != "" {
			resViaHead.Uri.Params.Set("branch", viaHead.Uri.Params.Get("branch"))
		}
		res.Header.Set(HeaderVia, resViaHead)
	}
	if !res.Header.Has(HeaderUserAgent) {
		res.Header.Set(HeaderUserAgent, defaultUserAgentHead)
	}
	if !res.Header.Has(HeaderContact) {
		res.Header.Set(HeaderContact, &AddressHeader{Uri: &Uri{IsEncrypted: false, User: ctx.sess.Id, Host: ctx.sess.transport.Conn().LocalAddr().String(), Params: map[string]string{
			"transport": ctx.sess.transport.Protocol(),
		}}})
	}
	_, err = ctx.sess.transport.Write(res.Bytes())
	return
}
