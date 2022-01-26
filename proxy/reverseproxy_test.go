package proxy

import (
	"testing"
)

func TestNewReverseProxy(t *testing.T) {
	serve := NewReverse(nil)
	_ = serve.Serve("192.168.4.169:5060")
}
