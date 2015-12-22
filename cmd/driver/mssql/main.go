package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/mssql"
	"github.com/hpcloud/cf-usb/driver/mssql/mssqlprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	var logger = lager.NewLogger("mssql-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	provisioner := mssqlprovisioner.NewMssqlProvisioner(logger)
	mssqldriver := driver.NewMssqlDriver(logger, provisioner)

	if err := p.RegisterName("mssql", mssqldriver); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
