package sip

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unsafe"
)

type Response struct {
	Proto         string
	Status        string
	StatusCode    int
	Header        *Header
	Body          []byte
	ContentLength int
	Request       *Request
}

func parseResponseLine(line string) (proto string, statusCode int, status string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	var err error
	s2 += s1 + 1
	proto = strings.TrimSpace(line[:s1])
	if statusCode, err = strconv.Atoi(strings.TrimSpace(line[s1:s2])); err == nil {
		ok = true
	}
	status = strings.TrimSpace(line[s2:])
	return
}

func ReadResponse(b *bufio.Reader) (res *Response, err error) {
	var (
		ok  bool
		buf []byte
	)
	res = &Response{}
	if buf, _, err = b.ReadLine(); err != nil {
		return
	}
	if res.Proto, res.StatusCode, res.Status, ok = parseResponseLine(string(buf)); !ok {
		return
	}
	if res.Header, err = readHeader(b); err != nil {
		return
	}
	res.ContentLength, _ = strconv.Atoi(res.Header.Get(HeaderContentLength).String())
	if res.ContentLength > 0 {
		res.Body = make([]byte, res.ContentLength)
		res.ContentLength, err = io.ReadFull(b, res.Body)
	}
	return
}
func (r *Response) Clone() *Response {
	res := &Response{
		Proto:         r.Proto,
		Status:        r.Status,
		StatusCode:    r.StatusCode,
		Header:        r.Header.Clone(),
		Body:          nil,
		ContentLength: r.ContentLength,
		Request:       r.Request,
	}
	if r.Body != nil {
		res.Body = make([]byte, len(r.Body))
		copy(res.Body, r.Body)
	}
	return res
}

func (r *Response) CallID() string {
	if head := r.Header.Get(HeaderCallID); head != nil {
		return head.(*PlainHeader).Content
	}
	return ""
}

func (r *Response) Bytes() []byte {
	str := r.String()
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{str, len(str)},
	))
}

func (r *Response) String() string {
	var sb strings.Builder
	if r.StatusCode == 0 {
		r.StatusCode = StatusOK
	}
	if r.Status == "" {
		r.Status = StatusText(r.StatusCode)
	}
	if r.Proto == "" {
		r.Proto = "SIP/2.0"
	}
	sb.WriteString(r.Proto + " " + strconv.Itoa(r.StatusCode) + " " + r.Status + "\r\n")
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

func NewResponse(code int, req *Request) *Response {
	res := &Response{
		Proto:      "SIP/2.0",
		Status:     StatusText(code),
		Header:     &Header{},
		StatusCode: code,
	}
	if req == nil {
		return res
	}
	if req.Header.Has(HeaderCSeq) {
		res.Header.Set(HeaderCSeq, req.Header.Get(HeaderCSeq).Clone())
	}
	if req.Header.Has(HeaderCallID) {
		res.Header.Set(HeaderCallID, req.Header.Get(HeaderCallID).Clone())
	}
	if req.Header.Has(HeaderFrom) {
		res.Header.Set(HeaderFrom, req.Header.Get(HeaderFrom).Clone())
	}
	if req.Header.Has(HeaderTo) {
		res.Header.Set(HeaderTo, req.Header.Get(HeaderTo).Clone())
	}
	return res
}
