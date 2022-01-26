package sip

import (
	"fmt"
	"testing"
)

func Test_parseViaHeaderFunc(t *testing.T) {
	if hv, err := parseViaHeaderFunc("SIP/2.0/TCP 192.168.4.169:40828;branch=z9hG4bK-524287-1---405263a6a9549471"); err != nil {
		t.Error(err)
	} else {
		fmt.Println(hv)
	}
}
