package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/hpcloud/cf-usb/lib/config"
)

type FileConfigProvider struct {
}

func NewFileConfigProvider() (*FileConfigProvider, error) {
	return nil, nil
}

func (k *FileConfigProvider) GetCLICommands(app Usb) []cli.Command {
	return []cli.Command{
		{
			Name: "fileConfigProvider",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Usage: "Path to the configuration file",
				},
			},
			Action: fileConfigProviderCommand(app),
			Usage:  `Provides a file for USB configuration`,
		},
	}
}

func fileConfigProviderCommand(app Usb) func(c *cli.Context) {
	return func(c *cli.Context) {
		configFilePath := c.String("path")

		if configFilePath == "" {
			cli.ShowCommandHelp(c, "fileConfigProvider")
			os.Exit(0)
		}

		_, err := os.Stat(configFilePath)
		if os.IsNotExist(err) {
			log.Fatal(fmt.Sprintf("Configuration file %s does not exist", configFilePath))
			os.Exit(1)
		}

		configuraiton := config.NewFileConfig(configFilePath)
		app.Run(configuraiton)
	}

}
