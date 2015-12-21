package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/mongo"
	"github.com/hpcloud/cf-usb/driver/mongo/mongoprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	var logger = lager.NewLogger("mongodb-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	provisioner := mongoprovisioner.New(logger)
	mongodriver := mongo.NewMongoDriver(logger, provisioner)

	if err := p.RegisterName("mongo", mongodriver); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
