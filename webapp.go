package webapp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django"
	"github.com/urfave/cli/v2"
)

var (
	// Standard directories
	staticDir      = "public"
	controllersDir = "controllers"
	modelsDir      = "models"
	viewsDir       = "views"

	stdDirs = []string{
		staticDir,
		controllersDir,
		modelsDir,
		viewsDir,
	}
)

type Webapp interface {
	fiber.Router
	Run() error
	AddCommand(*cli.Command)
}

type webapp struct {
	*fiber.App
	cli  *cli.App
	port int
}

func (w *webapp) Run() error {
	return w.cli.Run(os.Args)
}

func (w *webapp) AddCommand(cmd *cli.Command) {
	w.cli.Commands = append(w.cli.Commands, cmd)
}

func (w *webapp) Addr() string {
	return fmt.Sprintf(":%v", w.port)
}

func New() Webapp {
	config := fiber.Config{}
	if _, err := os.Stat(viewsDir); !os.IsNotExist(err) {
		config.Views = django.New(viewsDir, ".html")
	}
	app := &webapp{
		App: fiber.New(config),
	}
	app.cli = &cli.App{
		Name:  filepath.Base(os.Args[0]),
		Usage: "Built with github.com/sfreiberg/webapp",
		Commands: []*cli.Command{
			serverCmd(app),
			generateCmd(app),
		},
	}

	return app
}
