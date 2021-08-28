package webapp

import "github.com/urfave/cli/v2"

func serverCmd(w *webapp) *cli.Command {
	return &cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Start web server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   3000,
				Usage:   "port for the server to listen on",
				EnvVars: []string{"PORT"},
			},
		},
		Action: func(c *cli.Context) error {
			w.port = c.Int("port")
			return w.Listen(w.Addr())
		},
	}
}
