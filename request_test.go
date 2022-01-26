package sip

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestReadRequest(t *testing.T) {
	s := []byte(`INVITE sip:6363@192.168.4.169:48273;rinstance=a73836e86ca6411f SIP/2.0
Via: SIP/2.0/UDP 192.168.9.186:5060;rport;branch=z9hG4bKPjc04b36f9-2b54-4620-9693-cb7674e6954c
From: "15625229038" <sip:15625229038@192.168.9.186>;tag=73ddb69f-1471-454c-876f-a732b88f96fb
To: <sip:6363@192.168.4.169;rinstance=a73836e86ca6411f>
Contact: <sip:asterisk@192.168.9.186:5060>
Call-ID: 91182449-1b4a-4488-a9ae-d150a2271cb8
CSeq: 7286 INVITE
Allow: OPTIONS, SUBSCRIBE, NOTIFY, PUBLISH, INVITE, ACK, BYE, CANCEL, UPDATE, PRACK, REGISTER, REFER, MESSAGE
Supported: 100rel, timer, replaces, norefersub
Session-Expires: 1800
Min-SE: 90
Max-Forwards: 70
User-Agent: FPBX-13.0.192.8(13.27.0)
Content-Type: application/sdp
Content-Length:   362

v=0
o=- 961825727 961825727 IN IP4 192.168.9.186
s=Asterisk
c=IN IP4 192.168.9.186
t=0 0
m=audio 10558 RTP/AVP 0 8 3 111 18 101
a=rtpmap:0 PCMU/8000
a=rtpmap:8 PCMA/8000
a=rtpmap:3 GSM/8000
a=rtpmap:111 G726-32/8000
a=rtpmap:18 G729/8000
a=fmtp:18 annexb=no
a=rtpmap:101 telephone-event/8000
a=fmtp:101 0-16
a=ptime:20
a=maxptime:150
a=sendrecv`)
	if r, err := ReadRequest(bufio.NewReader(bytes.NewReader(s))); err != nil {
		t.Error(err)
	} else {
		fmt.Println(r)
	}
}
