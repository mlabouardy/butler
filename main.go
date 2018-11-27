package main

import (
	"fmt"
	"os"
	"strings"
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
		cli.Author{
			Name:  "Dominik Schröter",
			Email: "dominik.schroeter@bmw.de",
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
							Name:   "username, u",
							Usage:  "Jenkins username",
							EnvVar: "JENKINS_USER",
						},
						cli.StringFlag{
							Name:   "password, p",
							Usage:  "Jenkins password",
							EnvVar: "JENKINS_PASSWORD",
						},
						cli.StringFlag{
							Name:  "folder, f",
							Usage: "Jenkins Folder",
						},
					},
					Action: func(c *cli.Context) error {
						var server = getSanitizedUrl(c.String("server"))
						var username = c.String("username")
						var password = c.String("password")
						var folder = c.String("folder")

						if server == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}

						err := ImportJobs(server, username, password, folder)
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
							Name:   "username, u",
							Usage:  "Jenkins username",
							EnvVar: "JENKINS_USER",
						},
						cli.StringFlag{
							Name:  "folder, f",
							Usage: "Jenkins Folder",
						},
						cli.StringFlag{
							Name:   "password, p",
							Usage:  "Jenkins password",
							EnvVar: "JENKINS_PASSWORD",
						},
						cli.BoolFlag{
							Name:   "skip-folder, sf",
							Usage:  "Skip folder",
							EnvVar: "JENKINS_SKIP_FOLDER",
						},
					},
					Action: func(c *cli.Context) error {
						var server = getSanitizedUrl(c.String("server"))
						var username = c.String("username")
						var password = c.String("password")
						var skipFolder = c.Bool("skip-folder")
						var folder = c.String("folder")

						if server == "" {
							cli.ShowSubcommandHelp(c)
						}

						err := ExportJobs(server, folder, username, password, skipFolder)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						return nil
					},
				},
			},
		},
		{
			Name:  "credentials",
			Usage: "Jenkins Credential Management",
			Subcommands: []cli.Command{
				{
					Name:    "decrypt",
					Usage:   "Decrypt credentials of Jenkins folder",
					Aliases: []string{"d"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "server",
							Usage: "Jenkins url",
						},
						cli.StringFlag{
							Name:  "folder, f",
							Usage: "Jenkins Folder",
						},
						cli.StringFlag{
							Name:   "username, u",
							Usage:  "Jenkins username",
							EnvVar: "JENKINS_USER",
						},
						cli.StringFlag{
							Name:   "password, p",
							Usage:  "Jenkins password",
							EnvVar: "JENKINS_PASSWORD",
						},
					},
					Action: func(c *cli.Context) error {
						var url = getSanitizedUrl(c.String("server"))
						var username = c.String("username")
						var password = c.String("password")
						var folder = c.String("folder")

						if url == "" || folder == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}

						err := DecryptFolderCredentials(url, folder, username, password)
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						return nil
					},
				},
				{
					Name:    "apply",
					Usage:   "Apply (from STDIN) credentials of Jenkins folder",
					Aliases: []string{"a"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "server",
							Usage: "Jenkins url",
						},
						cli.StringFlag{
							Name:  "folder, f",
							Usage: "Jenkins Folder",
						},
						cli.StringFlag{
							Name:   "username, u",
							Usage:  "Jenkins username",
							EnvVar: "JENKINS_USER",
						},
						cli.StringFlag{
							Name:   "password, p",
							Usage:  "Jenkins password",
							EnvVar: "JENKINS_PASSWORD",
						},
					},
					Action: func(c *cli.Context) error {
						var url = getSanitizedUrl(c.String("server"))
						var username = c.String("username")
						var password = c.String("password")
						var folder = c.String("folder")

						if url == "" || folder == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}

						err := ApplyFolderCredentials(url, folder, username, password)
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
							Name:   "username, u",
							Usage:  "Jenkins username",
							EnvVar: "JENKINS_USER",
						},
						cli.StringFlag{
							Name:   "password, p",
							Usage:  "Jenkins password",
							EnvVar: "JENKINS_PASSWORD",
						},
					},
					Action: func(c *cli.Context) error {
						var server = getSanitizedUrl(c.String("server"))
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
							Name:   "username, u",
							Usage:  "Jenkins username",
							EnvVar: "JENKINS_USER",
						},
						cli.StringFlag{
							Name:   "password, p",
							Usage:  "Jenkins password",
							EnvVar: "JENKINS_PASSWORD",
						},
					},
					Action: func(c *cli.Context) error {
						var server = getSanitizedUrl(c.String("server"))
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

func getSanitizedUrl(url string) string {
	if url != "" && !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	url = strings.TrimRight(url, "/")

	return url
}
