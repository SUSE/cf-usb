package main

import (
	"github.com/codegangsta/cli"
)

//CLICommandProvider defines methods to be used in GetCommands
type CLICommandProvider interface {
	GetCLICommands(Usb) []cli.Command
}
