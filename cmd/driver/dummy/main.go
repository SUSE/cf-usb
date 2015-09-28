package main

import (
	"log"
	"net/rpc/jsonrpc"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/natefinch/pie"
)

func main() {
	log.SetPrefix("[dummydriver log] ")

	p := pie.NewProvider()
	if err := p.RegisterName("dummy", driver.NewDummyDriver()); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
