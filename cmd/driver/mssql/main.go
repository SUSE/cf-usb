package main

import (
	"log"
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/mssql"
	"github.com/hpcloud/cf-usb/driver/mssql/mssqlprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	log.SetPrefix("[mssql log] ")

	var logger = lager.NewLogger("mssql-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	provisioner := mssqlprovisioner.NewMssqlProvisioner(logger)
	mssqldriver := driver.NewMssqlDriver(logger, provisioner)

	if err := p.RegisterName("mssql", mssqldriver); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
