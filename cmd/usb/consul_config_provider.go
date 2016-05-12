package main

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/hpcloud/cf-usb/lib/config"

	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"
)

//ConsulConfigProvider provides a consul config
type ConsulConfigProvider struct {
}

//NewConsulConfigProvider returns a new consul config provider
func NewConsulConfigProvider() (*ConsulConfigProvider, error) {
	return nil, nil
}

//GetCLICommands returns the CLI Commands details from ConsulConfigProvider
func (k *ConsulConfigProvider) GetCLICommands(app Usb) []cli.Command {
	return []cli.Command{
		{
			Name: "consulConfigProvider",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address, a",
					Usage: "Consul address and port",
				},
				cli.StringFlag{
					Name:  "datacenter, d",
					Usage: "Consul datacenter",
				},
				cli.StringFlag{
					Name:  "username, u",
					Usage: "Consul username",
				},
				cli.StringFlag{
					Name:  "password,p",
					Usage: "Consul password",
				},
				cli.StringFlag{
					Name:  "schema, s",
					Usage: "Consul schema",
				},
				cli.StringFlag{
					Name:  "token, t",
					Usage: "Consul token",
				},
			},
			Action: consulConfigProviderCommand(app),
			Usage:  `Set consul configuration address`,
		},
	}
}

func consulConfigProviderCommand(app Usb) func(c *cli.Context) {
	return func(c *cli.Context) {
		logger := NewLogger(strings.ToLower(c.GlobalString("loglevel")))

		consulAddress := c.String("address")

		if consulAddress == "" {
			cli.ShowCommandHelp(c, "consulConfigProvider")
			os.Exit(0)
		}

		consulDatacenter := c.String("datacenter")
		consulUser := c.String("username")
		consulPass := c.String("password")
		consulSchema := c.String("schema")
		consulToken := c.String("token")

		var consulConfig api.Config
		consulConfig.Address = consulAddress
		consulConfig.Datacenter = consulDatacenter

		var auth api.HttpBasicAuth
		auth.Username = consulUser
		auth.Password = consulPass

		consulConfig.HttpAuth = &auth
		consulConfig.Scheme = consulSchema

		consulConfig.Token = consulToken

		provisioner, err := consul.New(&consulConfig)
		if err != nil {
			logger.Fatal("consul config provider", err)
		}

		configuraiton := config.NewConsulConfig(provisioner)

		app.Run(configuraiton, logger)
	}

}
