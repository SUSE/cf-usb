package lib

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/natefinch/pie"
)

type DriverProvider struct {
	driverType string
	client     *rpc.Client

	DriverProperties config.DriverProperties
}

func NewDriverProvider(driverType string, driverProperties config.DriverProperties) (DriverProvider, error) {
	provider := DriverProvider{}

	driverPath := filepath.Join("drivers", driverType)
	if runtime.GOOS == "windows" {
		driverPath = driverPath + ".exe"
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)
	if err != nil {
		return provider, err
	}

	provider.client = client
	provider.DriverProperties = driverProperties
	provider.driverType = driverType
	response, err := provider.Init(driverProperties)
	if err != nil {
		return provider, err
	}

	log.Println("Init driver reponse:", response)

	return provider, nil
}

func (p *DriverProvider) Init(driverProperties config.DriverProperties) (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.Init", p.driverType), driverProperties, &result)
	return result, err
}

func (p *DriverProvider) Provision(provisonRequest string) (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.Provision", p.driverType), provisonRequest, &result)
	return result, err
}

func (p *DriverProvider) GetCatalog() (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.GetCatalog", p.driverType), "", &result)
	return result, err
}
