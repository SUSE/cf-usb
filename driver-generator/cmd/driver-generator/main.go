package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/hpcloud/cf-usb/driver-generator/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "driver-generator"
	app.Usage = "usb driver generator client"
	app.Version = Version

	app.Commands = commands.Commands

	app.Run(os.Args)
}
