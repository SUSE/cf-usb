package commands

import "github.com/codegangsta/cli"

type driverInfo struct {
	DriverName    string
	GeneratedPath string
	BaseImport    string
}

var Commands []cli.Command = []cli.Command{
	{
		Name:        "generate",
		Description: "Generate a new drivers",
		Usage:       "driver-generator generate my-driver --path generated-driver",
		Action:      GenerateCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path",
				Usage: "the path of the new driver",
			}},
	},
}
