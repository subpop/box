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
				{
					Name:   "sync",
					Action: imageSync,
				},
				{
					Name:   "get",
					Action: imageGet,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "n,name",
							Required: true,
						},
						cli.StringFlag{
							Name:  "a,arch",
							Value: "x86_64",
						},
					},
				},
				{
					Name:   "remove",
					Action: imageRemove,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "n,name",
							Required: true,
						},
						cli.BoolFlag{
							Name: "force",
						},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
