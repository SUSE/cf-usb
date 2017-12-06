package main

import (
	"net"

	"strconv"

	"github.com/SUSE/cf-usb/lib/config"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/localip"

	"github.com/cloudfoundry/gibson"
	"github.com/cloudfoundry/yagnats"
)

//StartRouteRegistration starts regitration of routes based on the config received as parameter
func (usb *UsbApp) StartRouteRegistration(config *config.Config, log lager.Logger) {
	logger := log.Session("cf-router-registrar")

	natsMembers := usb.config.RoutesRegister.NatsMembers

	routesToRegister := map[int]string{}
	ip := ""

	_, brokerAPIPortString, err := net.SplitHostPort(usb.config.BrokerAPI.Listen)
	if err != nil {
		logger.Fatal("invalid-brokerapi-listening-address", err)
	}

	brokerAPIPort, err := strconv.Atoi(brokerAPIPortString)
	if err != nil {
		logger.Fatal("invalid-type-brokerapi-listening-address", err)
	}

	routesToRegister[brokerAPIPort] = usb.config.RoutesRegister.BrokerAPIHost

	if usb.config.ManagementAPI != nil && usb.config.RoutesRegister.ManagmentAPIHost != "" {
		mgmtaddr := usb.config.ManagementAPI.Listen
		_, mgmAPIPortString, err := net.SplitHostPort(mgmtaddr)
		if err != nil {
			logger.Fatal("invalid-management-api-listening-address", err)
		}

		mgmAPIPort, err := strconv.Atoi(mgmAPIPortString)
		if err != nil {
			logger.Fatal("invalid-type-brokerapi-listening-address", err)
		}

		routesToRegister[mgmAPIPort] = usb.config.RoutesRegister.ManagmentAPIHost
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
