package main

import (
	"flag"
	"fmt"
	"github.com/uole/sip/proxy"
	yaml "gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Listen string         `json:"listen" yaml:"listen"`
	Routes []*proxy.Route `json:"routes" yaml:"routes"`
}

var (
	configFlag = flag.String("config", "", "")
)

func main() {
	var (
		fp  *os.File
		err error
	)
	flag.Parse()
	cfg := &Config{Listen: "0.0.0.0:5060"}

	if *configFlag != "" {
		if fp, err = os.Open(*configFlag); err == nil {
			if err = yaml.NewDecoder(fp).Decode(cfg); err != nil {
				fmt.Println(err)
				_ = fp.Close()
				os.Exit(1)
			}
			_ = fp.Close()
		}
	}
	serve := proxy.NewReverse(cfg.Routes)
	_ = serve.Serve(cfg.Listen)
}
