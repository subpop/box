package main

import (
	"log"
	"os"

	"github.com/subpop/vm"

	"github.com/urfave/cli"
)

func main() {
	var err error
	var app *cli.App

	app = cli.NewApp()
	app.Name = "vm"
	app.Commands = []cli.Command{
		{
			Name:        "create",
			Usage:       "Creates a new domain from the specified image",
			UsageText:   "vm create [OPTION]... [IMAGE]",
			Description: "The create command defines new domains using the given image as a backing disk. If no --name option is specified, the domain is given a random name.",
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				return vm.Create(name, c.String("image"), c.StringSlice("disk"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name,n",
					Usage:    "Assign `NAME` to the domain",
					Required: true,
				},
				cli.StringSliceFlag{
					Name:  "disk,d",
					Usage: "Attach `FILE` to the domain as a secondary disk",
				},
			},
		},
		{
			Name:        "list",
			Usage:       "List defined domains",
			UsageText:   "vm list [OPTION]...",
			Description: "The list command prints a table of defined domains. By default, only active (running) domains are listed. Specify --all to print inactive domains as well.",
			Action: func(c *cli.Context) error {
				active := true
				inactive := false

				if c.Bool("inactive") {
					active = false
					inactive = true
				}

				if c.Bool("all") {
					active = true
					inactive = true
				}

				return vm.List(active, inactive)
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all",
					Usage: "Include inactive domains",
				},
				cli.BoolFlag{
					Name:  "inactive",
					Usage: "List only inactive domains",
				},
			},
		},
		{
			Name: "destroy",
			Action: func(c *cli.Context) error {
				return vm.Destroy(c.String("name"), c.Bool("force"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name,n",
					Required: true,
				},
				cli.BoolFlag{
					Name: "force,f",
				},
			},
		},
		{
			Name: "up",
			Action: func(c *cli.Context) error {
				return vm.Up(c.String("name"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name,n",
					Required: true,
				},
			},
		},
		{
			Name: "down",
			Action: func(c *cli.Context) error {
				return vm.Down(c.String("name"), c.Int("id"), c.Bool("force"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name,n",
					Required: true,
				},
				cli.BoolFlag{
					Name: "force",
				},
			},
		},
		{
			Name: "restart",
			Action: func(c *cli.Context) error {
				return vm.Restart(c.String("name"), c.Bool("force"), c.BoolT("graceful"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name,n",
					Required: true,
				},
				cli.BoolFlag{
					Name:     "force,f",
					Required: false,
				},
				cli.BoolTFlag{
					Name:     "graceful,g",
					Required: false,
				},
			},
		},
		{
			Name: "connect",
			Action: func(c *cli.Context) error {
				return vm.Connect(c.String("name"), c.String("mode"), c.String("user"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "name,n",
					Required: true,
				},
				cli.StringFlag{
					Name:  "mode,m",
					Value: "ssh",
				},
				cli.StringFlag{
					Name:  "user,u",
					Value: "root",
				},
			},
		},
		{
			Name: "image",
			Subcommands: []cli.Command{
				{
					Name: "list",
					Action: func(c *cli.Context) error {
						return vm.ImageList()
					},
				},
				{
					Name: "get",
					Action: func(c *cli.Context) error {
						return vm.ImageGet(c.String("url"))
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "url",
							Required: true,
						},
					},
				},
				{
					Name: "remove",
					Action: func(c *cli.Context) error {
						return vm.ImageRemove(c.String("name"), c.Bool("force"))
					},
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
		{
			Name: "template",
			Subcommands: []cli.Command{
				{
					Name: "list",
					Action: func(c *cli.Context) error {
						return vm.TemplateList(c.String("sort"))
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "sort,s",
							Value: "name",
						},
					},
				},
				{
					Name: "sync",
					Action: func(c *cli.Context) error {
						return vm.TemplateSync()
					},
				},
				{
					Name: "info",
					Action: func(c *cli.Context) error {
						return vm.TemplateInfo(c.String("name"), c.String("arch"))
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "name,n",
							Required: true,
						},
						cli.StringFlag{
							Name:  "arch,a",
							Value: "x86_64",
						},
					},
				},
				{
					Name: "get",
					Action: func(c *cli.Context) error {
						return vm.TemplateGet(c.String("name"), c.String("arch"))
					},
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
					Name: "remove",
					Action: func(c *cli.Context) error {
						return vm.ImageRemove(c.String("name"), c.Bool("force"))
					},
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
