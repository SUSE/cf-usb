package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/rabbitmq"
	"github.com/hpcloud/cf-usb/driver/rabbitmq/rabbitmqprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	var logger = lager.NewLogger("rabbitmq-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	provisioner := rabbitmqprovisioner.NewRabbitmqProvisioner(logger)
	rabbitmqDriver := rabbitmq.NewRabbitmqDriver(logger, provisioner)

	if err := p.RegisterName("rabbitmq", rabbitmqDriver); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
