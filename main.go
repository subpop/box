package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	var err error
	var app *cli.App

	app = cli.NewApp()
	app.Name = "box"
	app.Commands = []cli.Command{
		{
			Name: "image",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Action: imageList,
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
