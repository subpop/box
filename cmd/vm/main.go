package main

import (
	"fmt"
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
			UsageText:   "vm create [command options] [image name]",
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
			UsageText:   "vm list [command options]",
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
			Name:        "destroy",
			Usage:       "Destroy a domain",
			UsageText:   "vm destroy [command options] [domain name]",
			Description: "The destroy command destroys the specified domain, prompting the user for confirmation (unless --force is passed).",
			Action: func(c *cli.Context) error {
				return vm.Destroy(c.String("name"), c.Bool("force"))
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force,f",
					Usage: "Immediately destroy the domain, without prompting",
				},
			},
		},
		{
			Name:      "up",
			Usage:     "Start a domain",
			UsageText: "vm up [command options] [domain name]",
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				if name == "" {
					return vm.ErrDomainNameRequired
				}
				return vm.Up(name)
			},
		},
		{
			Name:      "down",
			Usage:     "Stop a domain",
			UsageText: "vm down [command options] [domain name]",
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				if name == "" {
					return vm.ErrDomainNameRequired
				}
				return vm.Down(name, c.Bool("force"), c.Bool("graceful"))
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force,f",
					Usage: "Immediately stop the domain, without prompting",
				},
				cli.BoolFlag{
					Name:  "graceful,g",
					Usage: "Power off the domain gracefully",
				},
			},
		},
		{
			Name:      "restart",
			Usage:     "Restart a domain",
			UsageText: "vm restart [command options] [domain name]",
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				if name == "" {
					return vm.ErrDomainNameRequired
				}
				return vm.Restart(name, c.Bool("force"), c.Bool("graceful"))
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force,f",
					Usage: "Immediately restart the domain, without prompting",
				},
				cli.BoolFlag{
					Name:  "graceful,g",
					Usage: "Restart the domain gracefully",
				},
			},
		},
		{
			Name:        "connect",
			Usage:       "Connect to a running domain",
			UsageText:   "vm connect [command options] [domain name]",
			Description: `Connect to a running domain. The --mode option changes the virtual device that is connected to. 'serial' connects to the domain's serial PTY. 'console' attempts to connect to a VirtIO PTY on the domain (if the domain supports VirtIO character devices). 'ssh' establishes an SSH session and attempts password authentication.`,
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				if name == "" {
					return vm.ErrDomainNameRequired
				}
				return vm.Connect(name, c.String("mode"), c.String("user"))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "mode,m",
					Usage: "Connection mode: serial, console, or ssh",
					Value: "serial",
				},
				cli.StringFlag{
					Name:  "user,u",
					Usage: "User to connect as over SSH",
					Value: "root",
				},
			},
		},
		{
			Name:  "image",
			Usage: "Manage backing disk images",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "List available backing disk images",
					Action: func(c *cli.Context) error {
						return vm.ImageList()
					},
				},
				{
					Name:      "get",
					Usage:     "Retrieve a new backing disk image",
					UsageText: "vm image get [command options] [URL or path]",
					Action: func(c *cli.Context) error {
						path := c.Args().First()
						if path == "" {
							return vm.ErrURLOrPathRequired
						}
						return vm.ImageGet(path)
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name,n",
							Usage: "Rename backing disk image to `NAME`",
						},
					},
				},
				{
					Name:      "remove",
					Usage:     "Remove a backing disk image",
					UsageText: "vm image remove [command options] [image name]",
					Action: func(c *cli.Context) error {
						name := c.Args().First()
						if name == "" {
							return vm.ErrImageNameRequired
						}
						return vm.ImageRemove(name, c.Bool("force"))
					},
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "force,f",
							Usage: "Force removal of a backing disk image without prompting",
						},
					},
				},
			},
		},
		{
			Name:  "template",
			Usage: "Manage backing disk templates from libguestfs",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "List templates available for import",
					Action: func(c *cli.Context) error {
						return vm.TemplateList(c.String("sort"))
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "sort,s",
							Usage: "Sort list by `VALUE`",
							Value: "name",
						},
					},
				},
				{
					Name:  "sync",
					Usage: "Refresh available templates from build service",
					Action: func(c *cli.Context) error {
						return vm.TemplateSync()
					},
				},
				{
					Name:      "info",
					Usage:     "Print details about a template",
					UsageText: "vm template info [command options] [template name]",
					Action: func(c *cli.Context) error {
						name := c.Args().First()
						if name == "" {
							return vm.ErrTemplateNameRequired
						}
						return vm.TemplateInfo(name, c.String("arch"))
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "arch,a",
							Usage: "Specify alternate architecture",
							Value: "x86_64",
						},
					},
				},
				{
					Name:      "get",
					Usage:     "Retrieve and prepare a template from build service",
					UsageText: "vm template get [command options] [template name]",
					Action: func(c *cli.Context) error {
						name := c.Args().First()
						if name == "" {
							return vm.ErrTemplateNameRequired
						}
						return vm.TemplateGet(name, c.String("arch"))
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "arch,a",
							Usage: "Specify alternative architecture",
							Value: "x86_64",
						},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
