package webapp

import (
	"sync"

	"github.com/urfave/cli/v2"
)

type tasks struct {
	tasks   []*cli.Command
	tasksMu sync.RWMutex
}

func (t *tasks) Add(task ...*cli.Command) {
	t.tasksMu.Lock()
	t.tasks = append(t.tasks, task...)
	t.tasksMu.Unlock()
}

func (t *tasks) All() []*cli.Command {
	t.tasksMu.RLock()
	tasks := make([]*cli.Command, len(t.tasks))
	copy(tasks, t.tasks)
	t.tasksMu.RUnlock()
	return tasks
}
