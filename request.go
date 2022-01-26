package sip

import (
	"bufio"
	"context"
	"github.com/google/uuid"
	"io"
	"strconv"
	"strings"
	"unsafe"
)

type Request struct {
	Method   Method
	Username string
	Address  string
	Proto    string
	Header   *Header
	Body     []byte
	Params   Map
	Context  context.Context
}

func (r *Request) WithContext(ctx context.Context) *Request {
	r.Context = ctx
	return r
}

func (r *Request) Clone() *Request {
	req := &Request{
		Method:   r.Method,
		Username: r.Username,
		Address:  r.Address,
		Proto:    r.Proto,
		Header:   r.Header.Clone(),
		Context:  r.Context,
	}
	if r.Params != nil {
		req.Params = r.Params.Clone()
	}
	if r.Body != nil {
		req.Body = make([]byte, len(r.Body))
		copy(req.Body[:], r.Body[:])
	}
	return req
}

func (r *Request) CallID() string {
	var callId string
	if head := r.Header.Get(HeaderCallID); head == nil {
		callId = uuid.New().String()
		r.Header.Set(HeaderCallID, &PlainHeader{Content: callId})
	} else {
		callId = head.(*PlainHeader).Content
	}
	return callId
}

func (r *Request) Bytes() []byte {
	str := r.String()
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{str, len(str)},
	))
}

func (r *Request) String() string {
	var sb strings.Builder
	sb.WriteString(string(r.Method) + " ")
	sb.WriteString("sip:")
	if r.Username != "" {
		sb.WriteString(r.Username + "@")
	}
	sb.WriteString(r.Address)
	if r.Params != nil {
		sb.WriteString(";" + r.Params.String())
	}
	sb.WriteString(" ")
	sb.WriteString(r.Proto)
	sb.WriteString("\r\n")
	if len(r.Body) == 0 {
		r.Header.Set(HeaderContentLength, &PlainHeader{Content: "0"})
	} else {
		r.Header.Set(HeaderContentLength, &PlainHeader{Content: strconv.Itoa(len(r.Body))})
	}
	sb.WriteString(r.Header.String())
	if r.Body != nil {
		sb.Write(r.Body)
	}
	return sb.String()
}

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return strings.TrimSpace(line[:s1]), strings.TrimSpace(line[s1+1 : s2]), strings.TrimSpace(line[s2+1:]), true
}

func ReadRequest(b *bufio.Reader) (req *Request, err error) {
	var (
		ok            bool
		method        string
		str           string
		buf           []byte
		pos           int
		contentLength int
	)
	req = &Request{}
	if buf, _, err = b.ReadLine(); err != nil {
		return
	}
	if method, str, req.Proto, ok = parseRequestLine(string(buf)); !ok {
		return
	}
	if pos = strings.Index(str, ":"); pos > -1 {
		str = str[pos+1:]
	}
	if pos = strings.Index(str, "@"); pos > -1 {
		req.Username = str[:pos]
		str = str[pos+1:]
	}
	if pos = strings.Index(str, ";"); pos > -1 {
		req.Address = str[:pos]
		if req.Params, err = parseMap(str[pos+1:]); err != nil {
			return
		}
	} else {
		req.Address = str
	}
	req.Method = Method(method)
	if req.Header, err = readHeader(b); err != nil {
		return
	}
	if req.Header.Has(HeaderContentLength) {
		contentLength, _ = strconv.Atoi(req.Header.Get(HeaderContentLength).String())
		if contentLength > 0 {
			req.Body = make([]byte, contentLength)
			contentLength, err = io.ReadFull(b, req.Body)
		}
	}
	return
}

func NewRequest(method Method, domain string) *Request {
	req := &Request{
		Method:  method,
		Address: domain,
		Proto:   "SIP/2.0",
		Header:  &Header{},
		Body:    nil,
	}
	return req
}

func NewDefaultRequest(method Method, domain string) *Request {
	req := NewRequest(method, domain)
	req.Header.Set(HeaderAllow, NewArrayHeader("INVITE", "ACK", "CANCEL", "BYE", "NOTIFY", "REFER", "MESSAGE", "OPTIONS", "INFO", "SUBSCRIBE"))
	req.Header.Set(HeaderSupported, NewArrayHeader("replaces", "norefersub", "extended-refer", "timer", "outbound", "path", "X-cisco-serviceuri"))
	req.Header.Set(HeaderAllowEvents, NewArrayHeader("presence", "kpml"))
	req.Header.Set(HeaderMaxForwards, NewMaxForwardHeader(70))
	return req
}
