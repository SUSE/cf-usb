package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	usb := NewUsbApp()

	app := cli.NewApp()
	app.Name = "usb"

	for _, command := range usb.GetCommands() {
		for _, cliCommand := range command.GetCLICommands(usb) {
			app.Commands = append(app.Commands, cliCommand)
		}
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "loglevel",
			Value: "debug",
			Usage: "Set log level (info,debug,error,fatal)",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
