package proxy

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/uole/sip"
	"github.com/uole/sip/pool"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	responseFeature = []byte("SIP")

	ErrorMissingToHead = errors.New("head to missing")
)

type (
	ReverseProxy struct {
		ctx                context.Context
		udpConn            *net.UDPConn
		processLocker      sync.RWMutex
		processes          map[string]*Process //处理铜须
		transChan          chan *Transaction   //事物处理
		routes             []*Route            //路由表
		relationshipLocker sync.RWMutex
		relationships      map[string]*Relationship //关系表
	}
)

//rewriteRequest 重写sip请求
func (rp *ReverseProxy) rewriteRequest(trans *Transaction) *sip.Request {
	originalRequest := trans.Request()
	rewriteRequest := originalRequest.Clone()
	//match address
	if rewriteRequest.Address == trans.transport.Addr().String() {
		if trans.Address() == trans.Caller().Addr().String() {
			rewriteRequest.Address = trans.Callee().Addr().String()
		} else {
			rewriteRequest.Address = trans.Caller().Addr().String()
		}
	}
	if originalRequest.Params != nil {
		rewriteRequest.Params = originalRequest.Params.Clone()
	}
	if originalRequest.Header.Has(sip.HeaderContact) {
		originalContactHeader := originalRequest.Header.Get(sip.HeaderContact).(*sip.AddressHeader)
		rewriteContactHeader := &sip.AddressHeader{
			Uri:    sip.NewUri(originalContactHeader.Uri.User, trans.transport.Addr().String(), sip.Map{}).EnableProtocol(),
			Params: originalContactHeader.Params.Clone(),
		}
		rewriteContactHeader.Uri.Params.Set("transport", trans.transport.Network())
		rewriteRequest.Header.Set(sip.HeaderContact, rewriteContactHeader)
	}
	if originalRequest.Header.Has(sip.HeaderVia) {
		originalViaHeader := originalRequest.Header.Get(sip.HeaderVia).(*sip.ViaHeader)
		rewriteViaHeader := &sip.ViaHeader{
			Protocol:        "SIP",
			ProtocolVersion: "2.0",
			Transport:       "UDP",
			Uri:             sip.NewUri("", trans.transport.Addr().String(), originalViaHeader.Uri.Params.Clone()),
		}
		rewriteRequest.Header.Set(sip.HeaderVia, rewriteViaHeader)
	}

	if rewriteRequest.Header.Has(sip.HeaderFrom) {
		fromHeader := rewriteRequest.Header.Get(sip.HeaderFrom).(*sip.AddressHeader)
		fromHeader.Uri.Params.Set("transport", trans.transport.Network())
		if rewrite, ok := trans.Rewrite(); ok {
			if fromHeader.Uri.Host == rewrite.From {
				fromHeader.Uri.Host = rewrite.To
				fromHeader.Uri.Port = 0
			} else if fromHeader.Uri.Host == rewrite.To {
				fromHeader.Uri.Host = rewrite.From
			}
		}
	}

	if rewriteRequest.Header.Has(sip.HeaderTo) {
		toHeader := rewriteRequest.Header.Get(sip.HeaderTo).(*sip.AddressHeader)
		toHeader.Uri.Params.Set("transport", trans.transport.Network())
		if rewrite, ok := trans.Rewrite(); ok {
			//呼出场景
			if toHeader.Uri.Host == rewrite.From {
				toHeader.Uri.Host = rewrite.To
				toHeader.Uri.Port = 0
			} else if toHeader.Uri.Host == rewrite.To {
				toHeader.Uri.Host = rewrite.From
			}
			//呼入场景
			if toHeader.Uri.Address() == trans.transport.Addr().String() {
				if trans.Address() == trans.Caller().Addr().String() {
					toHeader.Uri.SetAddress(trans.Callee().Addr().String())
				} else {
					toHeader.Uri.SetAddress(trans.Caller().Addr().String())
				}
			}
		}
	}
	return rewriteRequest
}

//rewriteResponse 重写sip响应消息
func (rp *ReverseProxy) rewriteResponse(trans *Transaction) *sip.Response {
	originalResponse := trans.Response()
	rewriteResponse := originalResponse.Clone()
	if originalResponse.Header.Has(sip.HeaderVia) {
		originalViaHeader := originalResponse.Header.Get(sip.HeaderVia).(*sip.ViaHeader)
		rewriteViaHeader := &sip.ViaHeader{
			Protocol:        "SIP",
			ProtocolVersion: "2.0",
			Transport:       "UDP",
			Uri:             sip.NewUri("", trans.Caller().Addr().String(), originalViaHeader.Uri.Params.Clone()),
		}
		if trans.Address() == trans.Caller().Addr().String() {
			rewriteViaHeader.Uri = sip.NewUri("", trans.Callee().Addr().String(), sip.Map{})
		} else {
			rewriteViaHeader.Uri = sip.NewUri("", trans.Caller().Addr().String(), sip.Map{})
		}
		if originalViaHeader.Uri.Params.Get("branch") != "" {
			rewriteViaHeader.Uri.Params.Set("branch", originalViaHeader.Uri.Params.Get("branch"))
		}
		rewriteViaHeader.Uri.Params.Set("rport", strconv.Itoa(rewriteViaHeader.Uri.Port))
		rewriteResponse.Header.Set(sip.HeaderVia, rewriteViaHeader)
	}
	if originalResponse.Header.Has(sip.HeaderContact) {
		originalContactHeader := originalResponse.Header.Get(sip.HeaderContact).(*sip.AddressHeader)
		rewriteContactHeader := &sip.AddressHeader{
			Uri:    sip.NewUri(originalContactHeader.Uri.User, trans.transport.Addr().String(), sip.Map{}).EnableProtocol(),
			Params: originalContactHeader.Params.Clone(),
		}
		rewriteContactHeader.Uri.Params.Set("transport", trans.transport.Network())
		rewriteResponse.Header.Set(sip.HeaderContact, rewriteContactHeader)
	}
	if rewriteResponse.Header.Has(sip.HeaderFrom) {
		fromHeader := rewriteResponse.Header.Get(sip.HeaderFrom).(*sip.AddressHeader)
		fromHeader.Uri.Params.Set("transport", trans.transport.Network())
		if rewrite, ok := trans.Rewrite(); ok {
			//呼出场景
			if fromHeader.Uri.Host == rewrite.To {
				fromHeader.Uri.Host = rewrite.From
			} else if fromHeader.Uri.Host == rewrite.From {
				fromHeader.Uri.Host = rewrite.To
			}
		}
	}
	if rewriteResponse.Header.Has(sip.HeaderTo) {
		toHeader := rewriteResponse.Header.Get(sip.HeaderTo).(*sip.AddressHeader)
		toHeader.Uri.Params.Set("transport", trans.transport.Network())
		if rewrite, ok := trans.Rewrite(); ok {
			//呼出场景
			if toHeader.Uri.Host == rewrite.To {
				toHeader.Uri.Host = rewrite.From
			} else if toHeader.Uri.Host == rewrite.From {
				toHeader.Uri.Host = rewrite.To
			}
			//呼入场景
			if toHeader.Uri.Address() == trans.transport.Addr().String() {
				if trans.Address() == trans.Caller().Addr().String() {
					toHeader.Uri.SetAddress(trans.Callee().Addr().String())
				} else {
					toHeader.Uri.SetAddress(trans.Caller().Addr().String())
				}
			}
			//呼入场景
			if toHeader.Uri.Address() == trans.Callee().Addr().String() {
				toHeader.Uri.SetAddress(trans.transport.Addr().String())
			}
		}
	}
	return rewriteResponse
}

//roundTripper 数据转发
func (rp *ReverseProxy) roundTripper(trans *Transaction) (err error) {
	if trans.message.Direction() == DirectionRequest {
		request := rp.rewriteRequest(trans)
		if request.Header.Has(sip.HeaderMaxForwards) {
			forwardHeader := request.Header.Get(sip.HeaderMaxForwards).(*sip.MaxForwardsHeader)
			forwardHeader.Forward = forwardHeader.Forward - 1
			if forwardHeader.Forward <= 0 {
				err = trans.Caller().Response(sip.NewResponse(sip.StatusLoopDetected, request))
				return
			}
		}
		if trans.Address() == trans.Caller().Addr().String() {
			err = trans.Callee().Request(request)
		} else {
			err = trans.Caller().Request(request)
		}
	} else {
		response := rp.rewriteResponse(trans)
		if trans.Address() == trans.Callee().Addr().String() {
			err = trans.Caller().Response(response)
		} else {
			err = trans.Callee().Response(response)
		}
	}
	return
}

//updateRelationship 更新绑定关系
func (rp *ReverseProxy) updateRelationship(conn Conn, msg *Message) *Relationship {
	var (
		domainName string
	)
	req := msg.Request()
	fromHead := req.Header.Get(sip.HeaderFrom).(*sip.AddressHeader)
	domainName = fromHead.Uri.Host
	//if domain rewrite rules exists
	for _, route := range rp.routes {
		if route.Domain == domainName {
			if route.RewriteTo != "" {
				domainName = route.RewriteTo
			}
			break
		}
	}
	username := fmt.Sprintf("%s@%s", fromHead.Uri.User, domainName)
	rp.relationshipLocker.Lock()
	defer rp.relationshipLocker.Unlock()
	relationship, ok := rp.relationships[username]
	if !ok {
		relationship = &Relationship{
			User:           fromHead.Uri.User,
			Domain:         domainName,
			OriginalDomain: fromHead.Uri.Host,
		}
		rp.relationships[username] = relationship
		log.Printf("bind user %s relationship %s", username, conn.Addr().String())
	}
	relationship.Conn = conn
	return relationship
}

//findRoute 查找请求的路由
func (rp *ReverseProxy) findRoute(req *sip.Request) (route *Route, err error) {
	fromHead := req.Header.Get(sip.HeaderFrom).(*sip.AddressHeader)
	domainName := fromHead.Uri.Host
	for _, r := range rp.routes {
		if r.Domain == domainName {
			route = r
			break
		}
	}
	if route == nil {
		err = fmt.Errorf("domain %s route not found", domainName)
	}
	return
}

//findRelationship 查找绑定关系
func (rp *ReverseProxy) findRelationship(req *sip.Request) (relationship *Relationship, err error) {
	var (
		ok bool
	)
	if !req.Header.Has(sip.HeaderTo) {
		err = ErrorMissingToHead
		return
	}
	if !req.Header.Has(sip.HeaderContact) {
		err = ErrorMissingToHead
		return
	}
	toHead := req.Header.Get(sip.HeaderTo).(*sip.AddressHeader)
	contactHead := req.Header.Get(sip.HeaderContact).(*sip.AddressHeader)
	rp.relationshipLocker.RLock()
	defer rp.relationshipLocker.RUnlock()
	username := fmt.Sprintf("%s@%s", toHead.Uri.User, toHead.Uri.Host)
	//如果直接找到对应的用户信息
	if relationship, ok = rp.relationships[username]; ok {
		return
	}
	//使用IP的方式进行查找数据
	username = fmt.Sprintf("%s@%s", toHead.Uri.User, contactHead.Uri.Host)
	if relationship, ok = rp.relationships[username]; ok {
		return
	}
	err = fmt.Errorf("%s relationship not found", username)
	return
}

//getProcess 获取一个处理器
func (rp *ReverseProxy) getProcess(conn Conn, msg *Message) (process *Process, err error) {
	var (
		ok           bool
		route        *Route
		relationship *Relationship
	)
	rp.processLocker.Lock()
	defer rp.processLocker.Unlock()
	if process, ok = rp.processes[msg.CallID()]; ok {
		return
	}
	if msg.Direction() == DirectionResponse {
		err = fmt.Errorf("not found")
		return
	}
	process = NewProcess(msg.CallID())
	process.caller = conn
	//bypass route
	if route, err = rp.findRoute(msg.Request()); err == nil {
		process.callee = newUDPConn(route.Address(), rp.udpConn)
		process.route = route
		rp.processes[msg.CallID()] = process
		if msg.Direction() == DirectionRequest && msg.Request().Method == sip.MethodRegister {
			rp.updateRelationship(conn, msg)
		}
		return
	}
	//find relationship
	if relationship, err = rp.findRelationship(msg.Request()); err == nil {
		process.callee = relationship.Conn
		process.relationship = relationship
		rp.processes[msg.CallID()] = process
	}
	return
}

func (rp *ReverseProxy) udpServe(addr string) (err error) {
	var (
		n          int
		proc       *Process
		remoteAddr *net.UDPAddr
		localAddr  *net.UDPAddr
	)
	if localAddr, err = net.ResolveUDPAddr("udp", addr); err != nil {
		return
	}
	if rp.udpConn, err = net.ListenUDP("udp", localAddr); err != nil {
		return
	}
	buf := make([]byte, 1024*32)
	for {
		if n, remoteAddr, err = rp.udpConn.ReadFromUDP(buf); err != nil {
			break
		}
		if n < 3 {
			continue
		}
		msg := &Message{}
		bytesReader := pool.GetBytesReader(buf[:n])
		bufioReader := pool.GetBufioReader(bytesReader)
		if bytes.Compare(buf[:3], responseFeature) == 0 {
			msg.direction = DirectionResponse
			msg.response, err = sip.ReadResponse(bufioReader)
		} else {
			msg.direction = DirectionRequest
			msg.request, err = sip.ReadRequest(bufioReader)
		}
		pool.PutBytesReader(bytesReader)
		pool.PutBufioReader(bufioReader)
		if err != nil {
			log.Printf("parse sip message error: %s", err.Error())
			continue
		}
		//获取处理程序
		if proc, err = rp.getProcess(&UdpConn{conn: rp.udpConn, addr: remoteAddr}, msg); err != nil {
			if msg.Direction() == DirectionRequest {
				res := sip.NewResponse(sip.StatusTemporarilyUnavailable, msg.Request())
				_, _ = rp.udpConn.WriteToUDP(res.Bytes(), remoteAddr)
			}
			log.Printf("get sip message %s process error: %s", msg.CallID(), err.Error())
			continue
		}
		trans := newTransaction(msg, proc, remoteAddr, newUDPTransport(rp.udpConn))
		trans.process.Push(msg)
		select {
		case rp.transChan <- trans:
		case <-rp.ctx.Done():
		case <-time.After(time.Millisecond * 100):
		}
	}
	return
}

func (rp *ReverseProxy) eventLoop() {
	for {
		select {
		case trans := <-rp.transChan:
			if err := rp.roundTripper(trans); err != nil {
				log.Printf(err.Error())
			}
		case <-rp.ctx.Done():
			return
		}
	}
}

//Serve 开启服务
func (rp *ReverseProxy) Serve(addr string) (err error) {
	go func() {
		err = rp.udpServe(addr)
	}()
	rp.eventLoop()
	return
}

//NewReverse 穿件一个代理服务
func NewReverse(routes []*Route) *ReverseProxy {
	proxy := &ReverseProxy{
		transChan:     make(chan *Transaction, 1024),
		ctx:           context.Background(),
		processes:     make(map[string]*Process),
		relationships: make(map[string]*Relationship),
		routes:        routes,
	}
	if proxy.routes == nil {
		proxy.routes = make([]*Route, 0)
	}
	return proxy
}
