package main

import (
	"log"
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/mysql"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	log.SetPrefix("[mysql log] ")

	var logger = lager.NewLogger("mysql-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	p := pie.NewProvider()

	mysqldriver := driver.NewMysqlDriver(logger)

	if err := p.RegisterName("mysql", mysqldriver); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
