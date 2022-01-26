package sip

import (
	"time"
)

type Transaction struct {
	ID        string
	CreatedAt time.Time
	c         chan *Response
}

func (t *Transaction) notify(res *Response) {
	if t.c != nil {
		select {
		case t.c <- res:
		default:
		}
	}
}

//Chan 获取事件
func (t *Transaction) Chan() chan *Response {
	return t.c
}

//Close 关闭事物
func (t *Transaction) Close() (err error) {
	close(t.c)
	return
}

func newTransaction(id string) *Transaction {
	return &Transaction{
		ID:        id,
		CreatedAt: time.Now(),
		c:         make(chan *Response, 1),
	}
}
