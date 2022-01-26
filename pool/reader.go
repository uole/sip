package pool

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

var (
	bytesReaderPool sync.Pool
	bufioReaderPool sync.Pool
)

func GetBytesReader(buf []byte) *bytes.Reader {
	if v := bytesReaderPool.Get(); v == nil {
		return bytes.NewReader(buf)
	} else {
		r := v.(*bytes.Reader)
		r.Reset(buf)
		return r
	}
}

func PutBytesReader(r *bytes.Reader) {
	bytesReaderPool.Put(r)
}

func GetBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v == nil {
		return bufio.NewReader(r)
	} else {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
}

func PutBufioReader(r *bufio.Reader) {
	bufioReaderPool.Put(r)
}
