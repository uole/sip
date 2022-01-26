package sip

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestReadResponse(t *testing.T) {
	s := []byte(`SIP/2.0 100 Trying
Via: SIP/2.0/UDP 192.168.4.169:5060;branch=z9hG4bK-524287-1---c34e1fe4153b4900
From: <sip:1001@192.168.9.185:5060;transport=UDP>;tag=dd669b44
To: <sip:15625229038@192.168.9.185:5060;transport=UDP>
Call-ID: 42VFkMGXZKZJ9Bz5Jfs3GQ..
CSeq: 2 INVITE
User-Agent: FreeSWITCH-mod_sofia/1.8.6~64bit
Content-Length: 0
`)
	if r, err := ReadResponse(bufio.NewReader(bytes.NewReader(s))); err != nil {
		t.Error(err)
	} else {
		fmt.Println(r)
	}
}
