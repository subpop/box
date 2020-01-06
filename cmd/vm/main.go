package main

import (
	"fmt"
	"io"
	"os"

	"github.com/subpop/vm"

	"github.com/urfave/cli/v2"
)

func main() {
	var err error
	var app *cli.App

	app = cli.NewApp()
	app.Name = "vm"
	app.Version = "0.1"
	app.Usage = "control virtual machines"
	app.Description = "`vm` is a program to manage and interact with virtual machines."
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:        "create",
			Usage:       "Creates a new domain from the specified image",
			UsageText:   "vm create [command options] [image name]",
			Description: "The create command defines new domains using the given image as a backing disk. If no --name option is specified, the domain is given a random name.",
			Action: func(c *cli.Context) error {
				image := c.Args().First()
				opts := vm.CreateOptions{
					ConnectAfterCreate:    !c.Bool("detach"),
					IsTransient:           c.Bool("transient"),
					CreateInitialSnapshot: !c.Bool("no-snapshot"),
				}
				cfg := vm.CreateConfig{
					UEFI: c.Bool("uefi"),
				}
				return vm.Create(c.String("name"), image, c.StringSlice("disk"), opts, cfg)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "name",
					Usage:   "Assign `NAME` to the domain",
					Aliases: []string{"n"},
				},
				&cli.StringSliceFlag{
					Name:    "disk",
					Usage:   "Attach `FILE` to the domain as a secondary disk",
					Aliases: []string{"d"},
				},
				&cli.BoolFlag{
					Name:  "detach",
					Usage: "Detach from the newly created domain",
				},
				&cli.BoolFlag{
					Name:    "transient",
					Usage:   "Create a non-persistent domain",
					Aliases: []string{"t"},
				},
				&cli.BoolFlag{
					Name:  "no-snapshot",
					Usage: "Disable taking an initial snapshot upon creation",
				},
				&cli.BoolFlag{
					Name:  "uefi",
					Usage: "Use UEFI boot loader",
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
				&cli.BoolFlag{
					Name:  "all",
					Usage: "Include inactive domains",
				},
				&cli.BoolFlag{
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
				name := c.Args().First()
				if name == "" {
					return vm.ErrDomainNameRequired
				}
				return vm.Destroy(name, c.Bool("force"))
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "force",
					Usage:   "Immediately destroy the domain, without prompting",
					Aliases: []string{"f"},
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
				return vm.Up(name, c.Bool("connect"))
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "connect",
					Usage:   "Immediately connect to the started domain",
					Aliases: []string{"c"},
				},
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
				&cli.BoolFlag{
					Name:    "force",
					Usage:   "Immediately stop the domain, without prompting",
					Aliases: []string{"f"},
				},
				&cli.BoolFlag{
					Name:    "graceful",
					Usage:   "Power off the domain gracefully",
					Aliases: []string{"g"},
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
				&cli.BoolFlag{
					Name:    "force",
					Usage:   "Immediately restart the domain, without prompting",
					Aliases: []string{"f"},
				},
				&cli.BoolFlag{
					Name:    "graceful",
					Usage:   "Restart the domain gracefully",
					Aliases: []string{"g"},
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
				return vm.Connect(name, c.String("mode"), c.String("user"), c.String("identity"))
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "mode",
					Usage:   "Connection mode: serial, console, or ssh",
					Aliases: []string{"m"},
					Value:   "serial",
				},
				&cli.StringFlag{
					Name:    "user",
					Usage:   "User to connect as over SSH",
					Aliases: []string{"u"},
					Value:   "root",
				},
				&cli.StringFlag{
					Name:    "identity",
					Usage:   "Attempt SSH authentication using `IDENTITY`",
					Aliases: []string{"i"},
				},
			},
		},
		{
			Name:        "inspect",
			Usage:       "Show details about a domain",
			UsageText:   "vm inspect [command options] [domain name]",
			Description: "Show details about a domain. Pass the --format option with 'json' or 'xml' as the argument to output in JSON or XML respectively.",
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				if name == "" {
					return vm.ErrDomainNameRequired
				}
				return vm.Inspect(name, c.String("format"))
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "format",
					Usage:   "Specify output format",
					Aliases: []string{"f"},
				},
			},
		},
		{
			Name:  "image",
			Usage: "Manage backing disk images",
			Subcommands: []*cli.Command{
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
						return vm.ImageGet(path, c.String("rename"), c.Bool("quiet"))
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "rename",
							Usage:   "Rename backing disk image to `NAME`",
							Aliases: []string{"r"},
						},
						&cli.BoolFlag{
							Name:    "quiet",
							Usage:   "No progress output",
							Aliases: []string{"q"},
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
						&cli.BoolFlag{
							Name:    "force",
							Usage:   "Force removal of a backing disk image without prompting",
							Aliases: []string{"f"},
						},
					},
				},
			},
		},
		{
			Name:  "template",
			Usage: "Manage backing disk templates from libguestfs",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List templates available for import",
					Action: func(c *cli.Context) error {
						return vm.TemplateList(c.String("sort"))
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "sort",
							Usage:   "Sort list by `VALUE`",
							Aliases: []string{"s"},
							Value:   "name",
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
						&cli.StringFlag{
							Name:    "arch",
							Usage:   "Specify alternate architecture",
							Aliases: []string{"a"},
							Value:   "x86_64",
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
						return vm.TemplateGet(name, c.String("arch"), c.Bool("quiet"))
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "arch",
							Usage:   "Specify alternative architecture",
							Aliases: []string{"a"},
							Value:   "x86_64",
						},
						&cli.BoolFlag{
							Name:    "quiet",
							Usage:   "No progress output",
							Aliases: []string{"q"},
						},
					},
				},
			},
		},
		{
			Name:  "snapshot",
			Usage: "Manage domain snapshots",
			Subcommands: []*cli.Command{
				{
					Name:      "list",
					Usage:     "List snapshots for a domain",
					UsageText: "vm snapshot list [command options] [domain name]",
					Action: func(c *cli.Context) error {
						name := c.Args().First()
						if name == "" {
							return vm.ErrDomainNameRequired
						}
						return vm.SnapshotList(name)
					},
				},
				{
					Name:      "create",
					Usage:     "Take a new snapshot for a domain",
					UsageText: "vm snapshot create [command options] [domain name]",
					Action: func(c *cli.Context) error {
						domain := c.Args().First()
						if domain == "" {
							return vm.ErrDomainNameRequired
						}
						return vm.SnapshotCreate(domain, c.String("name"))
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "name",
							Usage:   "Create a snapshot with `NAME`",
							Aliases: []string{"n"},
						},
					},
				},
				{
					Name:      "remove",
					Usage:     "Remove a snapshot for a domain",
					UsageText: "vm snapshot remove [command options] [domain name]",
					Action: func(c *cli.Context) error {
						domain := c.Args().First()
						if domain == "" {
							return vm.ErrDomainNameRequired
						}
						return vm.SnapshotRemove(domain, c.String("snapshot"))
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "snapshot",
							Usage:    "Remove snapshot named `NAME`",
							Aliases:  []string{"s"},
							Required: true,
						},
					},
				},
				{
					Name:      "revert",
					Usage:     "Revert a domain to snapshot",
					UsageText: "vm snapshot revert [command options] [domain name]",
					Action: func(c *cli.Context) error {
						domain := c.Args().First()
						if domain == "" {
							return vm.ErrDomainNameRequired
						}
						return vm.SnapshotRevert(domain, c.String("snapshot"))
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "snapshot",
							Usage:    "Revert to `SNAPSHOT`",
							Aliases:  []string{"s"},
							Required: true,
						},
					},
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:   "generate-man-page",
			Hidden: true,
		},
		&cli.BoolFlag{
			Name:   "generate-fish-completion",
			Hidden: true,
		},
		&cli.BoolFlag{
			Name:   "generate-markdown",
			Hidden: true,
		},
	}
	app.BashComplete = bashComplete
	app.Action = func(c *cli.Context) error {
		type GenerationFunc func() (string, error)
		var generationFunc GenerationFunc
		if c.Bool("generate-fish-completion") {
			generationFunc = c.App.ToFishCompletion
		} else if c.Bool("generate-man-page") {
			generationFunc = c.App.ToMan
		} else if c.Bool("generate-markdown") {
			generationFunc = c.App.ToMarkdown
		}
		data, err := generationFunc()
		if err != nil {
			return err
		}
		fmt.Println(data)
		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		vm.LogErrorAndExit(err)
	}
}

func bashCompleteCommand(cmd *cli.Command, w io.Writer) {
	for _, name := range cmd.Names() {
		fmt.Fprintf(w, "%v\n", name)
	}

	for _, flag := range cmd.VisibleFlags() {
		for _, name := range flag.Names() {
			if len(name) > 1 {
				fmt.Fprintf(w, "--%v\n", name)
			} else {
				fmt.Fprintf(w, "-%v\n", name)
			}
		}
	}

	for _, command := range cmd.Subcommands {
		bashCompleteCommand(command, w)
	}
}

func bashComplete(c *cli.Context) {
	for _, command := range c.App.VisibleCommands() {
		bashCompleteCommand(command, c.App.Writer)
	}

	// global flags
	for _, flag := range c.App.VisibleFlags() {
		for _, name := range flag.Names() {
			if len(name) > 1 {
				fmt.Fprintf(c.App.Writer, "--%v\n", name)
			} else {
				fmt.Fprintf(c.App.Writer, "-%v\n", name)
			}
		}
	}
}
