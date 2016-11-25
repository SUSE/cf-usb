package main

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/hpcloud/cf-usb/lib/config"
)

//MysqlConfigProvider provides a mysql config
type MysqlConfigProvider struct {
}

//NewMysqlConfigProvider returns a new mysql config provider
func NewMysqlConfigProvider() (*MysqlConfigProvider, error) {
	return nil, nil
}

//GetCLICommands returns the CLI Commands details from MysqlConfigProvider
func (k *MysqlConfigProvider) GetCLICommands(app Usb) []cli.Command {
	return []cli.Command{
		{
			Name: "mysqlConfigProvider",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address, a",
					Usage: "Mysql address and port",
				},
				cli.StringFlag{
					Name:  "database, db",
					Usage: "Mysql database",
				},
				cli.StringFlag{
					Name:  "username, u",
					Usage: "Mysql username",
				},
				cli.StringFlag{
					Name:  "password,p",
					Usage: "Mysql password",
				},
			},
			Action: mysqlConfigProviderCommand(app),
			Usage:  `Set mysql configuration address`,
		},
	}
}

func mysqlConfigProviderCommand(app Usb) func(c *cli.Context) {
	return func(c *cli.Context) {
		logger := NewLogger(strings.ToLower(c.GlobalString("loglevel")))

		mysqlAddress := c.String("address")

		if mysqlAddress == "" {
			cli.ShowCommandHelp(c, "mysqlConfigProvider")
			os.Exit(0)
		}

		mysqlDatabase := c.String("database")
		mysqlUser := c.String("username")
		mysqlPass := c.String("password")

		configuration, err := config.NewMysqlConfig(mysqlAddress, mysqlUser, mysqlPass, mysqlDatabase)
		if err != nil {
			logger.Fatal("mysql-config-provider-init", err)
		}
		app.Run(configuration, logger)
	}

}
