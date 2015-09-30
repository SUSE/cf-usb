package main

import (
	"fmt"
	"log"

	"github.com/hpcloud/gocfbroker"
	"github.com/hpcloud/gocfbroker/boltdb"
)

var echoService EchoService

type EchoService struct {
}

func (e *EchoService) Provision(instanceID string, req gocfbroker.ProvisionRequest) (gocfbroker.ProvisionResponse, error) {
	log.Printf("\n\nProvision called with:\ninstanceID: %s\nprovisionRequest: %+v", instanceID, req)
	res := gocfbroker.ProvisionResponse{
		DashboardURL: fmt.Sprintf("http://localhost/instance/%s/service/%s/dashboard", instanceID, req.ServiceID),
	}
	return res, nil
}

func (e *EchoService) Deprovision(instanceID, serviceID, planID string) error {
	log.Printf("\n\nDeprovision called with:\ninstanceID: %s\nserviceID: %s\nplanID: %s", instanceID, serviceID, planID)
	return nil
}

func (e *EchoService) Update(instanceID string, req gocfbroker.UpdateProvisionRequest) error {
	log.Printf("\n\nUpdate called with:\ninstanceID: %s\nupdateProvisionRequest: %+v", instanceID, req)
	return nil
}

func (e *EchoService) Bind(instanceID, bindingID string, req gocfbroker.BindingRequest) (gocfbroker.BindingResponse, error) {
	log.Printf("\n\nBind called with:\ninstanceID: %s\nbindingID: %s\nbindRequest: %+v", instanceID, bindingID, req)
	res := gocfbroker.BindingResponse{
		Credentials: gocfbroker.MakeJSONRawMessage(`{"username": "user", "password": "password"}`),
	}
	return res, nil
}

func (e *EchoService) Unbind(instanceID, bindingID, serviceID, planID string) error {
	log.Printf("\n\nUnbind called with:\ninstanceID: %s\nbindingID: %s\nserviceID: %s\nplanID: %s", instanceID, bindingID, serviceID, planID)
	return nil
}

type Config struct {
	BoltFilename string `json:"bolt_filename"`
	BoltBucket   string `json:"bolt_bucket"`

	gocfbroker.Options
}

func main() {
	echoService := &EchoService{}
	var config Config

	if err := gocfbroker.LoadConfig("config.json", &config); err != nil {
		log.Fatalln("failed to load config:", err)
	}

	db, err := boltdb.New(config.BoltFilename, config.BoltBucket)
	if err != nil {
		log.Fatalln("failed to open database:", err)
	}

	broker, err := gocfbroker.New(echoService, db, config.Options)
	if err != nil {
		log.Fatalln(err)
	}

	broker.Start()
}
