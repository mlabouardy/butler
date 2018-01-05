package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "butler"
	app.Usage = "Import/Export Jenkins Jobs"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Mohamed Labouardy",
			Email: "mohamed@labouardy.com",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "jobs",
			Usage: "Jenkins Jobs Management",
			Subcommands: []cli.Command{
				{
					Name:    "import",
					Usage:   "Import Jenkins Jobs",
					Aliases: []string{"i"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "server, s",
							Usage: "Jenkins server",
						},
						cli.StringFlag{
							Name:  "username, u",
							Usage: "Jenkins username",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "Jenkins password",
						},
					},
					Action: func(c *cli.Context) error {
						var server = c.String("server")
						var username = c.String("username")
						var password = c.String("password")

						if server == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}

						err := ImportJobs(server, username, password)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						return nil
					},
				},
				{
					Name:    "export",
					Usage:   "Export Jenkins Jobs",
					Aliases: []string{"e"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "server, s",
							Usage: "Jenkins server",
						},
						cli.StringFlag{
							Name:  "username, u",
							Usage: "Jenkins username",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "Jenkins password",
						},
					},
					Action: func(c *cli.Context) error {
						var server = c.String("server")
						var username = c.String("username")
						var password = c.String("password")

						if server == "" {
							cli.ShowSubcommandHelp(c)
						}

						err := ExportJobs(server, username, password)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						return nil
					},
				},
			},
		},
		{
			Name:  "plugins",
			Usage: "Jenkins Plugins Management",
			Subcommands: []cli.Command{
				{
					Name:    "import",
					Usage:   "Import Jenkins Plugins",
					Aliases: []string{"i"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "server, s",
							Usage: "Jenkins server",
						},
						cli.StringFlag{
							Name:  "username, u",
							Usage: "Jenkins username",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "Jenkins password",
						},
					},
					Action: func(c *cli.Context) error {
						var server = c.String("server")
						var username = c.String("username")
						var password = c.String("password")

						if server == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}

						err := ImportPlugins(server, username, password)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						return nil
					},
				},
				{
					Name:    "export",
					Usage:   "Export Jenkins Plugins",
					Aliases: []string{"e"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "server, s",
							Usage: "Jenkins server",
						},
						cli.StringFlag{
							Name:  "username, u",
							Usage: "Jenkins username",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "Jenkins password",
						},
					},
					Action: func(c *cli.Context) error {
						var server = c.String("server")
						var username = c.String("username")
						var password = c.String("password")

						if server == "" {
							cli.ShowSubcommandHelp(c)
						}

						err := ExportPlugins(server, username, password)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						return nil
					},
				},
			},
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Command not found %q !", command)
	}
	app.Run(os.Args)
}
