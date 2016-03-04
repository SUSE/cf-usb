package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func validateArgsCount(c *cli.Context, expected int) {
	if len(c.Args()) < expected || len(c.Args()) > expected {
		fmt.Println("Invlid number of arguments. Expected:", expected, "got:", len(c.Args()))
		showHelpAndExit(c)
	}
}

func showHelpAndExit(c *cli.Context) {
	cli.ShowCommandHelp(c, c.Command.Name)
	os.Exit(1)
}
