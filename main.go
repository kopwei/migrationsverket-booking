package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/kopwei/migrationsverket-booking/cmd"
	"github.com/urfave/cli"
)

// VERSION is the global value for software version
var VERSION = "v0.1.0-dev"

func main() {
	if err := mainErr(); err != nil {
		logrus.Fatal(err)
	}
}

func mainErr() error {
	app := cli.NewApp()
	app.Name = "migrationsverketquerier"
	app.Usage = "Check free time slot on Migration Board's website"
	app.Version = VERSION
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Author = "kopwei"
	app.Commands = []cli.Command{
		cmd.Check(),
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Debug logging",
		},
	}
	return app.Run(os.Args)
}
