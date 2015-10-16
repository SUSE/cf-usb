package main

import(
	"log"
	"net/rpc/jsonrpc"
	"os"
	
	"github.com/hpcloud/cf-usb/driver/redis"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
)

func main() {
	log.SetPrefix("[redis log] ")
	
	var logger = lager.NewLogger("redis-driver")
	
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))
	
	p := pie.NewProvider()
	if err := p.RegisterName("redis", redis.NewRedisDriver(logger)); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
