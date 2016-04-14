package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/passthrough"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	var logger = lager.NewLogger("passthrough-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	driver := passthrough.NewPassthroughDriver(logger)
	if err := p.RegisterName("passthrough", driver); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)

}
