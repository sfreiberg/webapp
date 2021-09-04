package webapp

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"text/tabwriter"

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

func routesCmd(w *webapp) *cli.Command {
	return &cli.Command{
		Name:  "routes",
		Usage: "Display routes",
		Action: func(c *cli.Context) error {
			tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			out := "%s\t%s\t%s\t\n"
			fmt.Fprintf(tw, out, "Path", "Method", "Controller")
			fmt.Fprintf(tw, out, "====", "======", "==========")

			for _, r := range w.sortedRoutes() {
				for _, handler := range r.Handlers {
					funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
					fmt.Fprintf(tw, out, r.Path, r.Method, funcName)
				}
			}
			return tw.Flush()
		},
	}
}

func tasksCmd(w *webapp) *cli.Command {
	return &cli.Command{
		Name:        "tasks",
		Usage:       "Display tasks",
		Aliases:     []string{"t"},
		Subcommands: w.tasks.All(),
	}
}
