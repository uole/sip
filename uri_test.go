package sip

import (
	"fmt"
	"testing"
)

func TestMap_Set(t *testing.T) {
	var m Map
	m.Set("a", "b")
	m.Set("aaa", "sdas\tdsds")
	t.Log(m.Get("a"))
	t.Log(m.String())
}

func Test_parseUri(t *testing.T) {
	if uri, err := parseUri("sip:1000:15625229038@192.168.4.169:40828;rinstance=e7be6d7faa64ed3f;transport=tcp?a=b&c=\"dsa\t !#$$ d\"&d=f"); err != nil {
		t.Error(err)
	} else {
		fmt.Println(uri)
	}
}
