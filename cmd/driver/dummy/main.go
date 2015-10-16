package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/dummy"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	var logger = lager.NewLogger("dummy-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	driver := dummydriver.NewDummyDriver(logger)
	if err := p.RegisterName("dummy", driver); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)

}
