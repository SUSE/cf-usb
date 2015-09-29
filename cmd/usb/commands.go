package main

import (
	"github.com/codegangsta/cli"
)

type CLICommandProvider interface {
	GetCLICommands(Usb) []cli.Command
}
