package sip

import (
	"bytes"
	"context"
	"github.com/uole/sip/pool"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	responseFeature = []byte("SIP")
)

type UDPTransport struct {
	conn         *net.UDPConn
	transMutex   sync.RWMutex
	reqChan      chan *Request
	transactions []*Transaction
}

func (tp *UDPTransport) Protocol() string {
	return ProtoUDP
}

func (tp *UDPTransport) Conn() net.Conn {
	return tp.conn
}

func (tp *UDPTransport) Request() chan *Request {
	return tp.reqChan
}

//traceTransaction 提交一个事物
func (tp *UDPTransport) traceTransaction(t *Transaction) {
	tp.transMutex.RLock()
	defer tp.transMutex.RUnlock()
	if tp.transactions == nil {
		tp.transactions = make([]*Transaction, 0)
	}
	tp.transactions = append(tp.transactions, t)
}

//notifyTransaction 通知一个事物完成
func (tp *UDPTransport) notifyTransaction(res *Response) (err error) {
	tp.transMutex.Lock()
	defer tp.transMutex.Unlock()
	callID := res.CallID()
	for _, trans := range tp.transactions {
		if trans.ID == callID {
			trans.notify(res)
			break
		}
	}
	return
}

//releaseTransaction 释放一个指定的事物
func (tp *UDPTransport) releaseTransaction(trans *Transaction) {
	tp.transMutex.Lock()
	defer tp.transMutex.Unlock()
	for i, v := range tp.transactions {
		if trans.ID == v.ID {
			tp.transactions = append(tp.transactions[:i], tp.transactions[i+1:]...)
			return
		}
	}
}

//Dial 新建立一个连接
func (tp *UDPTransport) Dial(addr string) (err error) {
	var udpAddr *net.UDPAddr
	if udpAddr, err = net.ResolveUDPAddr("udp", addr); err != nil {
		return
	}
	if tp.conn, err = net.DialUDP("udp", nil, udpAddr); err != nil {
		return
	}
	go tp.exchange()
	return
}

//exchange
func (tp *UDPTransport) exchange() {
	var (
		n          int
		err        error
		buf        []byte
		p          []byte
		res        *Response
		req        *Request
		addr       *net.UDPAddr
		isResponse bool
	)
	buf = make([]byte, 1024*10)
	for {
		if n, err = tp.conn.Read(buf); err != nil {
			continue
		}
		if n < 3 {
			continue
		}
		p = buf[:n]
		//parse the body
		bytesReader := pool.GetBytesReader(buf[:n])
		bufioReader := pool.GetBufioReader(bytesReader)
		if bytes.Compare(p[:3], responseFeature) == 0 {
			res, err = ReadResponse(bufioReader)
			isResponse = true
		} else {
			req, err = ReadRequest(bufioReader)
			isResponse = false
		}
		pool.PutBytesReader(bytesReader)
		pool.PutBufioReader(bufioReader)
		//parse failed
		if err != nil {
			log.Printf("parse buffer from %s: %s error: %s", addr.String(), string(p), err.Error())
			continue
		}
		if isResponse {
			err = tp.notifyTransaction(res)
		} else {
			select {
			case tp.reqChan <- req:
			case <-time.After(time.Millisecond * 200):
				log.Printf("put %s request timeout", req.Method)
			}
		}
	}
}

func (tp *UDPTransport) Write(p []byte) (n int, err error) {
	if tp.conn != nil {
		return tp.conn.Write(p)
	} else {
		err = io.ErrClosedPipe
	}
	return
}

func (tp *UDPTransport) Do(ctx context.Context, req *Request, callback ProcessFunc) (err error) {
	var (
		ok    bool
		res   *Response
		trans *Transaction
	)
	if _, err = tp.conn.Write(req.Bytes()); err != nil {
		return
	}
	trans = newTransaction(req.CallID())
	tp.traceTransaction(trans)
	defer tp.releaseTransaction(trans)
	for {
		select {
		case res = <-trans.Chan():
			ok, err = callback(res)
			if ok || err != nil {
				return
			}
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}
}

func (tp *UDPTransport) Close() (err error) {
	if tp.conn != nil {
		err = tp.conn.Close()
	}
	return
}

func NewUDPTransport() Transport {
	return &UDPTransport{reqChan: make(chan *Request, 100)}
}
