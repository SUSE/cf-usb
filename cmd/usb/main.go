package main

import (
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func runMain(writer io.Writer) {
	log.SetPrefix("[dummydriver log] ")
	usb := NewUsbApp()

	app := cli.NewApp()
	app.Name = "usb"

	for _, command := range usb.GetCommands() {
		for _, cliCommand := range command.GetCLICommands(usb) {
			app.Commands = append(app.Commands, cliCommand)
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
