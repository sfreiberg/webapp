package webapp

import (
	"os"

	"github.com/urfave/cli/v2"
)

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

func generateCmd(w *webapp) *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Usage:   "Generators",
		Subcommands: []*cli.Command{
			generateDirectoriesCmd(w),
		},
	}
}

func generateDirectoriesCmd(w *webapp) *cli.Command {
	return &cli.Command{
		Name:    "directories",
		Aliases: []string{"d"},
		Usage:   "Generate standard directory structure",
		Action: func(c *cli.Context) error {
			for _, dir := range stdDirs {
				if err := os.Mkdir(dir, 0744); err != nil {
					if !os.IsExist(err) {
						return err
					}
				}
			}
			return nil
		},
	}
}
