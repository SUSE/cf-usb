package main

import (
	"log"
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/mongo"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	log.SetPrefix("[mongodb log] ")

	var logger = lager.NewLogger("mongodb-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()
	if err := p.RegisterName("mongo", driver.NewMongoDriver(logger)); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
