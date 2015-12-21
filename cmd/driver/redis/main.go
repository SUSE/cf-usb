package main

import (
	"net/rpc/jsonrpc"
	"os"

	"github.com/hpcloud/cf-usb/driver/redis"
	"github.com/hpcloud/cf-usb/driver/redis/redisprovisioner"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	var logger = lager.NewLogger("redis-driver")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	provisioner := redisprovisioner.NewRedisProvisioner(logger)
	p := pie.NewProvider()
	if err := p.RegisterName("redis", redis.NewRedisDriver(logger, provisioner)); err != nil {
		logger.Fatal("register-plugin", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
