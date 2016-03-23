package commands

import "github.com/codegangsta/cli"

type driverInfo struct {
	DriverName         string
	GeneratedPath      string
	GenerateParameters bool
	BaseImport         string
}

var Commands []cli.Command = []cli.Command{
	{
		Name:        "generate",
		Description: "Generate a new drivers",
		Usage:       "driver-generator generate my-driver --path generated-driver --accept-user-parameters",
		Action:      GenerateCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path",
				Usage: "the path of the new driver",
			},
			cli.BoolFlag{
				Name:  "accept-user-parameters",
				Usage: "Specify if the driver support user provided parameters",
			}},
	},
}
