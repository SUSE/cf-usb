package lib

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"runtime"

	"github.com/natefinch/pie"
)

type Provider struct {
	driverType string
	client     *rpc.Client
}

func NewDriverProvider(driverType string) {
	log.SetPrefix("[master log] ")

	if runtime.GOOS == "windows" {
		driverType = driverType + ".exe"
	}

	driverPath := fmt.Sprintf("drivers/%s", driverType)

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)
	if err != nil {
		log.Fatalf("Error running driver: %s", err)
	}

	defer client.Close()

	p := Provider{driverType, client}
	res, err := p.Provision("master")
	if err != nil {
		log.Fatalf("error calling Provision: %s", err)
	}
	log.Printf("Response from plugin: %q", res)
}

func (p *Provider) Provision(provisonRequest string) (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.Provision", p.driverType), provisonRequest, &result)
	return result, err
}

func (p *Provider) GetCatalog() (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.GetCatalog", p.driverType), "", &result)
	return result, err
}
