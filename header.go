package sip

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/rs/xid"
	"io"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
)

var (
	funcMap = make(map[string]ParserHeaderFunc)
)

func init() {
	funcMap["Via"] = parseViaHeaderFunc
	funcMap["Contact"] = parseAddressHeaderFunc
	funcMap["From"] = parseAddressHeaderFunc
	funcMap["To"] = parseAddressHeaderFunc
	funcMap["CSeq"] = parseSequenceHeaderFunc
	funcMap["Allow"] = parseArrayHeaderFunc
	funcMap["Supported"] = parseArrayHeaderFunc
	funcMap["Allow-Events"] = parseArrayHeaderFunc
	funcMap["Max-Forwards"] = parseMaxForwardHeaderFunc
	funcMap["Authorization"] = parseAuthorizationHeaderFunc
	funcMap["WWW-Authenticate"] = parseAuthorizationHeaderFunc
}

const (
	HeaderVia                = "Via"
	HeaderMaxForwards        = "Max-Forwards"
	HeaderContact            = "Contact"
	HeaderFrom               = "From"
	HeaderTo                 = "To"
	HeaderCallID             = "Call-ID"
	HeaderCSeq               = "CSeq"
	HeaderExpires            = "Expires"
	HeaderAllow              = "Allow"
	HeaderSupported          = "Supported"
	HeaderUserAgent          = "User-Agent"
	HeaderAllowEvents        = "Allow-Events"
	HeaderContentLength      = "Content-Length"
	HeaderContentType        = "Content-Type"
	HeaderAuthorization      = "Authorization"
	HeaderWWWAuthenticate    = "WWW-Authenticate"
	HeaderProxyAuthorization = "Proxy-Authorization"
	HeaderDate               = "Date"
	HeaderReason             = "Reason"
	HeaderRequire            = "Require"
	HeaderSessionExpires     = "Session-Expires"
	HeaderMinSE              = "Min-SE"
)

type (
	//ParserHeaderFunc 头部解析函数
	ParserHeaderFunc func(s string) (Value, error)

	Value interface {
		String() string
		Clone() Value
	}

	Header struct {
		mu     sync.RWMutex
		Keys   []string
		Values map[string]Value
	}

	PlainHeader struct {
		Content string
	}

	MaxForwardsHeader struct {
		Forward int
	}

	ViaHeader struct {
		Protocol        string
		ProtocolVersion string
		Transport       string
		Uri             *Uri
	}

	AuthorizationHeader struct {
		Method    string
		Realm     string
		Nonce     string
		Algorithm string
		QOP       string
		Username  string
		Response  string
		CNonce    string
		NC        string
		Uri       *Uri
	}

	SequenceHeader struct {
		Method   Method
		Sequence int
	}

	ArrayHeader struct {
		Values []string
	}

	AddressHeader struct {
		DisplayName string
		Uri         *Uri
		Params      Map
	}
)

func (h *MaxForwardsHeader) String() string {
	return strconv.Itoa(h.Forward)
}

func (h *MaxForwardsHeader) Clone() Value {
	return &MaxForwardsHeader{Forward: h.Forward}
}

func AttachParseFunc(s string, f ParserHeaderFunc) {
	funcMap[s] = f
}

func (h *AuthorizationHeader) String() string {
	var sb strings.Builder
	sb.WriteString(h.Method + " ")
	if h.Username != "" {
		sb.WriteString("username=\"" + h.Username + "\", ")
	}
	if h.Realm != "" {
		sb.WriteString("realm=\"" + h.Realm + "\", ")
	}
	if h.Nonce != "" {
		sb.WriteString("nonce=\"" + h.Nonce + "\", ")
	}
	if h.Uri != nil {
		sb.WriteString("uri=\"" + h.Uri.String() + "\", ")
	}
	if h.Response != "" {
		sb.WriteString("response=\"" + h.Response + "\", ")
	}
	if h.CNonce != "" {
		sb.WriteString("cnonce=\"" + h.CNonce + "\", ")
	}
	if h.NC != "" {
		sb.WriteString("nc=" + h.NC + ", ")
	}
	if h.QOP != "" {
		sb.WriteString("qop=" + h.QOP + ", ")
	}
	if h.Algorithm != "" {
		sb.WriteString("algorithm=" + h.Algorithm + ", ")
	}
	return strings.TrimRight(sb.String(), ", ")
}

func (h *AuthorizationHeader) Clone() Value {
	head := &AuthorizationHeader{
		Method:    h.Method,
		Realm:     h.Realm,
		Nonce:     h.Nonce,
		Algorithm: h.Algorithm,
		QOP:       h.QOP,
		Username:  h.Username,
		Response:  h.Response,
		CNonce:    h.CNonce,
		NC:        h.NC,
	}
	if h.Uri != nil {
		head.Uri = h.Uri.Clone()
	}
	return head
}

func NewAuthorizationResponseHeader(username, password string, req *AuthorizationHeader) *AuthorizationHeader {
	head := &AuthorizationHeader{
		Method:    req.Method,
		Realm:     req.Realm,
		Nonce:     req.Nonce,
		Algorithm: req.Algorithm,
		QOP:       req.QOP,
		Username:  username,
		Response:  "",
		CNonce:    xid.New().String(),
		NC:        "00000001",
		Uri:       &Uri{Host: req.Realm},
	}
	//username:relam:password
	h1 := hex.EncodeToString(MD5([]byte(username + ":" + req.Realm + ":" + password)))
	//method:uri
	h2 := hex.EncodeToString(MD5([]byte("REGISTER:" + head.Uri.String())))
	//HA1:nonce:nc:cnonce:qop:HA2
	head.Response = hex.EncodeToString(MD5([]byte(h1 + ":" + req.Nonce + ":" + head.NC + ":" + head.CNonce + ":" + req.QOP + ":" + h2)))
	return head
}

func (h *ViaHeader) String() string {
	var sb strings.Builder
	if h.Protocol == "" {
		h.Protocol = "SIP"
	}
	if h.ProtocolVersion == "" {
		h.ProtocolVersion = "2.0"
	}
	if h.Transport == "" {
		h.Transport = "UDP"
	}
	sb.WriteString(h.Protocol + "/" + h.ProtocolVersion + "/" + h.Transport)
	sb.WriteString(" ")
	sb.WriteString(h.Uri.String())
	return sb.String()
}

func (h *ViaHeader) Clone() Value {
	return &ViaHeader{
		Protocol:        h.Protocol,
		ProtocolVersion: h.ProtocolVersion,
		Transport:       h.Transport,
		Uri:             h.Uri.Clone(),
	}
}

func (h *ArrayHeader) String() string {
	return strings.Join(h.Values, ", ")
}

func (h *ArrayHeader) Clone() Value {
	hc := &ArrayHeader{Values: make([]string, len(h.Values))}
	for i, s := range h.Values {
		hc.Values[i] = s
	}
	return hc
}

func NewArrayHeader(s ...string) *ArrayHeader {
	return &ArrayHeader{Values: s}
}

func (h *SequenceHeader) String() string {
	return strconv.Itoa(h.Sequence) + " " + h.Method.String()
}

func (h *SequenceHeader) Clone() Value {
	return &SequenceHeader{
		Method:   h.Method,
		Sequence: h.Sequence,
	}
}

func NewSequenceHeader(method Method, seq int) *SequenceHeader {
	return &SequenceHeader{Sequence: seq, Method: method}
}

func (h *AddressHeader) String() string {
	var s string
	if h.DisplayName != "" {
		s = "\"" + h.DisplayName + "\" "
	}
	s += "<" + h.Uri.String() + ">"
	if h.Params != nil && len(h.Params) > 0 {
		s += ";" + h.Params.String()
	}
	return s
}

func (h *AddressHeader) Clone() Value {
	hc := &AddressHeader{
		DisplayName: h.DisplayName,
		Uri:         h.Uri.Clone(),
		Params:      h.Params.Clone(),
	}
	return hc
}

func (h *PlainHeader) String() string {
	return h.Content
}

func (h *PlainHeader) Clone() Value {
	return &PlainHeader{Content: h.Content}
}

func NewPlainHeader(v interface{}) *PlainHeader {
	return &PlainHeader{Content: fmt.Sprint(v)}
}

func NewMaxForwardHeader(b int) *MaxForwardsHeader {
	return &MaxForwardsHeader{Forward: b}
}

func (h *Header) Clone() *Header {
	vv := &Header{}
	h.mu.Lock()
	defer h.mu.Unlock()
	vv.Keys = make([]string, len(h.Keys))
	for i, k := range h.Keys {
		vv.Keys[i] = k
	}
	vv.Values = make(map[string]Value)
	for k, v := range h.Values {
		vv.Values[k] = v.Clone()
	}
	return vv
}

func (h *Header) Set(name string, value Value) {
	h.mu.Lock()
	defer h.mu.Unlock()
	name = textproto.CanonicalMIMEHeaderKey(name)
	if h.Keys == nil {
		h.Keys = make([]string, 0)
	}
	if h.Values == nil {
		h.Values = make(map[string]Value)
	}
	if _, ok := h.Values[name]; !ok {
		h.Keys = append(h.Keys, name)
	}
	h.Values[name] = value
}

//Get 获取指定的头信息
func (h *Header) Get(name string) Value {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if v, ok := h.Values[textproto.CanonicalMIMEHeaderKey(name)]; ok {
		return v
	}
	return nil
}

func (h *Header) Has(name string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if _, ok := h.Values[textproto.CanonicalMIMEHeaderKey(name)]; ok {
		return true
	}
	return false
}

//String 返回字符串数据
func (h *Header) String() string {
	var (
		ok  bool
		val Value
		sb  strings.Builder
	)
	for _, k := range h.Keys {
		if val, ok = h.Values[k]; !ok {
			continue
		}
		sb.WriteString(textproto.CanonicalMIMEHeaderKey(k))
		sb.WriteString(": ")
		sb.WriteString(val.String())
		sb.WriteString("\r\n")
	}
	sb.WriteString("\r\n")
	return sb.String()
}

//parseAddressHeaderFunc 解析地址信息
func parseAddressHeaderFunc(s string) (header Value, err error) {
	var (
		pos      int
		tagBegin int
		tagEnd   int
		length   int
	)
	tagBegin, tagEnd = -1, -1
	hv := &AddressHeader{Uri: &Uri{}}
	length = len(s)
	for pos = 0; pos < length; pos++ {
		if s[pos] == '<' {
			hv.DisplayName = strings.Trim(strings.TrimSpace(s[:pos]), "\"")
			tagBegin = pos
		}
		if s[pos] == '>' {
			tagEnd = pos
		}
	}
	if tagBegin == -1 || tagEnd == -1 {
		err = fmt.Errorf("missing '<>' %s", s)
		return
	}
	ss := s[tagBegin+1 : tagEnd]
	if hv.Uri, err = parseUri(ss); err != nil {
		return
	}
	if tagEnd != length-1 {
		hv.Params, err = parseMap(s[tagEnd+1:])
	}
	header = hv
	return
}

func parseAuthorizationHeaderFunc(s string) (header Value, err error) {
	var (
		pos int
		key string
		val string
	)
	hv := &AuthorizationHeader{}
	if pos = strings.Index(s, " "); pos == -1 {
		return
	}
	hv.Method = s[:pos]
	ss := strings.Split(s[pos+1:], ",")

	for _, sp := range ss {
		if pos = strings.Index(sp, "="); pos != -1 {
			key = strings.TrimSpace(sp[:pos])
			val = strings.Trim(strings.TrimSpace(sp[pos+1:]), "\"")
			switch strings.ToLower(key) {
			case "username":
				hv.Username = val
			case "realm":
				hv.Realm = val
			case "nonce":
				hv.Nonce = val
			case "response":
				hv.Response = val
			case "cnonce":
				hv.CNonce = val
			case "nc":
				hv.NC = val
			case "qop":
				hv.QOP = val
			case "algorithm":
				hv.Algorithm = val
			case "uri":
				hv.Uri, err = parseUri(val)
			}
		}
	}
	header = hv
	return
}

//parseArrayHeaderFunc 解析数组信息
func parseArrayHeaderFunc(s string) (header Value, err error) {
	hv := &ArrayHeader{}
	ss := strings.Split(s, ",")
	hv.Values = make([]string, len(ss))
	for i, vs := range ss {
		hv.Values[i] = strings.TrimSpace(vs)
	}
	header = hv
	return
}

//parseAddressHeaderFunc 解析纯文本头
func parseMaxForwardHeaderFunc(s string) (header Value, err error) {
	h := &MaxForwardsHeader{}
	h.Forward, err = strconv.Atoi(strings.TrimSpace(s))
	header = h
	return
}

//parseAddressHeaderFunc 解析纯文本头
func parsePlainsHeaderFunc(s string) (header Value, err error) {
	header = &PlainHeader{Content: s}
	return
}

//parseSequenceHeaderFunc 解析seq头信息
func parseSequenceHeaderFunc(s string) (header Value, err error) {
	hv := &SequenceHeader{}
	ss := strings.Split(s, " ")
	if len(ss) == 2 {
		hv.Method = Method(strings.TrimSpace(ss[1]))
		hv.Sequence, err = strconv.Atoi(ss[0])
	} else {
		err = fmt.Errorf("unknown string %s", s)
	}
	header = hv
	return
}

func parseViaHeaderFunc(s string) (header Value, err error) {
	var (
		ps string
	)
	hv := &ViaHeader{}
	ss := strings.Split(s, " ")
	if len(ss) < 2 {
		err = fmt.Errorf("unknown string '%s'", s)
		return
	}
	ps = ss[1]
	ss = strings.Split(ss[0], "/")
	if len(ss) == 3 {
		hv.Protocol, hv.ProtocolVersion, hv.Transport = ss[0], ss[1], ss[2]
	} else if len(ss) == 2 {
		hv.Protocol, hv.ProtocolVersion = ss[0], ss[1]
		hv.Transport = ProtoUDP
	} else {
		err = fmt.Errorf("invalid protocol string %s", ss[0])
		return
	}
	hv.Uri, err = parseUri(ps)
	header = hv
	return
}

//parseHeader 解析头部
func parseHeader(s string) (key string, value Value, err error) {
	var (
		pos int
		ok  bool
		str string
		fun ParserHeaderFunc
	)
	pos = strings.Index(s, ":")
	if pos == -1 {
		err = fmt.Errorf("unexpected multi-line response: %s", s)
		return
	}
	key = s[:pos]
	str = strings.TrimSpace(s[pos+1:])
	if fun, ok = funcMap[key]; ok {
		value, err = fun(str)
	} else {
		value, err = parsePlainsHeaderFunc(str)
	}
	return
}

//readHeader read head message
func readHeader(b *bufio.Reader) (header *Header, err error) {
	var (
		line  string
		key   string
		value Value
	)
	header = &Header{}
	tr := textproto.NewReader(b)
	for {
		if line, err = tr.ReadLine(); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		//读取完毕
		if len(line) == 0 {
			break
		}
		if key, value, err = parseHeader(strings.TrimSpace(line)); err != nil {
			continue
		}
		header.Set(key, value)
	}
	return
}
