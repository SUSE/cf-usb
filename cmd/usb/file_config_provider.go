package main

import (
	"github.com/codegangsta/cli"
	"github.com/hpcloud/cf-usb/lib/config"
	"os"
	"strings"
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
		logger := NewLogger(strings.ToLower(c.GlobalString("loglevel")))
		configFilePath := c.String("path")

		if configFilePath == "" {
			cli.ShowCommandHelp(c, "fileConfigProvider")
			os.Exit(0)
		}

		_, err := os.Stat(configFilePath)
		if os.IsNotExist(err) {
			logger.Fatal("configuration-file-not-found", err)
		}

		configuraiton := config.NewFileConfig(configFilePath)
		app.Run(configuraiton, logger)
	}
}
