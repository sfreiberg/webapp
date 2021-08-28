package webapp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

type Webapp interface {
	fiber.Router
	Run() error
}

type webapp struct {
	*fiber.App
	cli  *cli.App
	port int
}

func (w *webapp) Run() error {
	return w.cli.Run(os.Args)
}

func (w *webapp) Addr() string {
	return fmt.Sprintf(":%v", w.port)
}

func New() Webapp {
	app := &webapp{App: fiber.New()}
	app.cli = &cli.App{
		Name:  filepath.Base(os.Args[0]),
		Usage: "Built with github.com/sfreiberg/webapp",
		Commands: []*cli.Command{
			serverCmd(app),
		},
	}

	return app
}
