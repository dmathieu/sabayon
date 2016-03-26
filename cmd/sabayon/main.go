package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/dmathieu/sabayon/commands"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	app := cli.NewApp()
	app.Name = "sabayon"
	app.Usage = "Manage Letsencrypt certificates on heroku"
	app.Commands = []cli.Command{
		commands.SetupCmd,
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
