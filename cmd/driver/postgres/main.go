package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/postgres"
	"github.com/hpcloud/cf-usb/driver/postgres/postgresprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	var logger = lager.NewLogger("postgres-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	provisioner := postgresprovisioner.NewPostgresProvisioner(logger)
	postgresDriver := postgres.NewPostgresDriver(logger, provisioner)

	if err := p.RegisterName("postgres", postgresDriver); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
