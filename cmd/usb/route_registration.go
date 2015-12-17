package main

import (
	"net"

	"strconv"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/localip"

	"github.com/cloudfoundry/gibson"
	"github.com/cloudfoundry/yagnats"
)

func (usb *UsbApp) StartRouteRegistration(config *config.Config, log lager.Logger) {
	logger := log.Session("cf-router-registrar")

	natsMembers := usb.config.RoutesRegister.NatsMembers

	routesToRegister := map[int]string{}
	ip := ""

	_, brokerApiPortString, err := net.SplitHostPort(usb.config.BrokerAPI.Listen)
	if err != nil {
		logger.Fatal("invalid-brokerapi-listening-address", err)
	}

	brokerApiPort, err := strconv.Atoi(brokerApiPortString)
	if err != nil {
		logger.Fatal("invalid-type-brokerapi-listening-address", err)
	}

	routesToRegister[brokerApiPort] = usb.config.RoutesRegister.BrokerAPIHost

	if usb.config.ManagementAPI != nil && usb.config.RoutesRegister.ManagmentAPIHost != "" {
		mgmtaddr := usb.config.ManagementAPI.Listen
		_, mgmApiPortString, err := net.SplitHostPort(mgmtaddr)
		if err != nil {
			logger.Fatal("invalid-management-api-listening-address", err)
		}

		mgmApiPort, err := strconv.Atoi(mgmApiPortString)
		if err != nil {
			logger.Fatal("invalid-type-brokerapi-listening-address", err)
		}

		routesToRegister[mgmApiPort] = usb.config.RoutesRegister.ManagmentAPIHost
	}

	ip, err = localip.LocalIP()
	if err != nil {
		logger.Fatal("error-discovering-localip", err)
	}

	natsConn, err := yagnats.Connect(natsMembers)
	if err != nil {
		logger.Fatal("nats-connetion-failed", err)
	}

	client := gibson.NewCFRouterClient(ip, natsConn)

	logger.Debug("start-greeting")
	err = client.Greet()
	if err != nil {
		logger.Fatal("greet-failed", err)
	}

	for port, host := range routesToRegister {
		logger.Info("start-register", lager.Data{"port": port, "host": host})
		err = client.Register(port, host)
		if err != nil {
			logger.Fatal("register-failed", err)
		}
	}
}
