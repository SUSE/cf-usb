package main

import (
	"log"
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/mongo"
	"github.com/hpcloud/cf-usb/driver/mongo/mongoprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	log.SetPrefix("[mongodb log] ")

	var logger = lager.NewLogger("mongodb-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	provisioner := mongoprovisioner.New(logger)
	mongodriver := mongo.NewMongoDriver(logger, provisioner)

	if err := p.RegisterName("mongo", mongodriver); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
