package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/mysql"
	"github.com/hpcloud/cf-usb/driver/mysql/mysqlprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {

	var logger = lager.NewLogger("mysql-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	provisioner := mysqlprovisioner.New(logger)
	mysqldriver := driver.NewMysqlDriver(logger, provisioner)

	if err := p.RegisterName("mysql", mysqldriver); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
