package webapp

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django"
	"github.com/urfave/cli/v2"
)

var (
	// Standard directories
	controllersDir = "controllers"
	modelsDir      = "models"
	tasksDir       = "tasks"
	staticDir      = "public"
	viewsDir       = "views"

	stdDirs = []string{
		controllersDir,
		modelsDir,
		tasksDir,
		staticDir,
		viewsDir,
	}
)

type Webapp interface {
	fiber.Router
	Run() error
	AddCommand(*cli.Command)
	AddTasks(...*cli.Command)
}

type webapp struct {
	*fiber.App
	tasks *tasks
	cli   *cli.App
	port  int
}

func New() Webapp {
	config := fiber.Config{}
	if _, err := os.Stat(viewsDir); !os.IsNotExist(err) {
		config.Views = django.New(viewsDir, ".html")
	}
	app := &webapp{
		App: fiber.New(config),
		tasks: &tasks{
			tasks: []*cli.Command{},
		},
	}

	// serve static files if the static directory exists
	if _, err := os.Stat(staticDir); !os.IsNotExist(err) {
		app.Static("/", staticDir)
	}

	return app
}

func (w *webapp) Run() error {
	w.cli = &cli.App{
		Name:  filepath.Base(os.Args[0]),
		Usage: "Built with github.com/sfreiberg/webapp",
		Commands: []*cli.Command{
			generateCmd(w),
			routesCmd(w),
			serverCmd(w),
			tasksCmd(w),
		},
	}

	return w.cli.Run(os.Args)
}

func (w *webapp) AddTasks(task ...*cli.Command) {
	w.tasks.Add(task...)
}

func (w *webapp) AddCommand(cmd *cli.Command) {
	w.cli.Commands = append(w.cli.Commands, cmd)
}

func (w *webapp) Addr() string {
	return fmt.Sprintf(":%v", w.port)
}

func (w *webapp) sortedRoutes() []*fiber.Route {
	r := []*fiber.Route{}
	for i := range w.Stack() {
		for j := range w.Stack()[i] {
			r = append(r, w.Stack()[i][j])
		}
	}
	sort.Sort(routes(r))
	return r
}

type routes []*fiber.Route

func (r routes) Len() int {
	return len(r)
}

func (r routes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r routes) Less(i, j int) bool {
	if r[i].Path == r[j].Path {
		return r[i].Method < r[j].Method
	}
	return r[i].Path < r[j].Path
}
