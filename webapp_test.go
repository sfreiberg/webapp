package webapp

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

func ok(c *fiber.Ctx) error { return c.SendString("ok") }

func TestAddr(t *testing.T) {
	w := &webapp{port: 8080}
	if got := w.Addr(); got != ":8080" {
		t.Errorf("Addr() = %q; want :8080", got)
	}
}

func TestTasks(t *testing.T) {
	ts := &tasks{}
	ts.Add(&cli.Command{Name: "a"}, &cli.Command{Name: "b"})
	ts.Add(&cli.Command{Name: "c"})

	all := ts.All()
	var names []string
	for _, c := range all {
		names = append(names, c.Name)
	}
	if !slices.Equal(names, []string{"a", "b", "c"}) {
		t.Fatalf("All() names = %v", names)
	}

	// All must return a copy that callers can modify freely.
	all[0] = &cli.Command{Name: "changed"}
	if ts.All()[0].Name != "a" {
		t.Error("modifying the slice returned by All() changed the stored tasks")
	}
}

func TestTasksConcurrent(t *testing.T) {
	ts := &tasks{}
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			ts.Add(&cli.Command{Name: "t"})
			_ = ts.All()
		})
	}
	wg.Wait()
	if n := len(ts.All()); n != 50 {
		t.Errorf("got %d tasks; want 50", n)
	}
}

func TestSortedRoutes(t *testing.T) {
	t.Chdir(t.TempDir()) // no views/public dirs
	w := New().(*webapp)
	w.Post("/b", ok)
	w.Get("/b", ok)
	w.Get("/a", ok)
	w.Delete("/a", ok)

	var got []string
	for _, r := range w.sortedRoutes() {
		got = append(got, r.Path+" "+r.Method)
	}
	// Each path/method appears in order of path, then method. Fiber also
	// registers HEAD for every GET route.
	want := []string{"/a DELETE", "/a GET", "/a HEAD", "/b GET", "/b HEAD", "/b POST"}
	if !slices.Equal(got, want) {
		t.Errorf("sortedRoutes() =\n%v\nwant\n%v", got, want)
	}
}

func TestNewServesStaticAndViews(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	for _, d := range []string{staticDir, viewsDir} {
		if err := os.Mkdir(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(staticDir, "hello.txt"), []byte("static file"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(viewsDir, "index.html"), []byte("Hello, {{ name }}!"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := New().(*webapp)
	w.Get("/page", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{"name": "<World>"})
	})

	get := func(path string) (int, string) {
		t.Helper()
		resp, err := w.Test(httptest.NewRequest("GET", path, nil))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		return resp.StatusCode, string(b)
	}

	if code, body := get("/hello.txt"); code != 200 || body != "static file" {
		t.Errorf("static: got %d %q", code, body)
	}
	// The django engine autoescapes variables.
	if code, body := get("/page"); code != 200 || body != "Hello, &lt;World&gt;!" {
		t.Errorf("view: got %d %q", code, body)
	}
	if code, _ := get("/missing"); code != 404 {
		t.Errorf("missing: got %d; want 404", code)
	}
}

func TestGenerateDirectories(t *testing.T) {
	t.Chdir(t.TempDir())
	w := New().(*webapp)
	app := &cli.App{Name: "test", Commands: []*cli.Command{generateCmd(w)}}

	// Running twice must succeed: existing directories are not an error.
	for range 2 {
		if err := app.Run([]string{"test", "generate", "directories"}); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range stdDirs {
		fi, err := os.Stat(d)
		if err != nil {
			t.Fatalf("%s not created: %v", d, err)
		}
		if !fi.IsDir() {
			t.Errorf("%s is not a directory", d)
		}
		// A directory that is readable by group/other must also be
		// searchable (x) by them, whatever the umask is.
		perm := fi.Mode().Perm()
		if (perm&0o040 != 0 && perm&0o010 == 0) || (perm&0o004 != 0 && perm&0o001 == 0) {
			t.Errorf("%s has mode %v: readable but not searchable", d, perm)
		}
	}
}

func TestTasksCommandRunsTask(t *testing.T) {
	t.Chdir(t.TempDir())
	w := New().(*webapp)
	ran := false
	w.AddTasks(&cli.Command{
		Name:   "seed",
		Action: func(*cli.Context) error { ran = true; return nil },
	})
	app := &cli.App{Name: "test", Commands: []*cli.Command{tasksCmd(w)}}
	if err := app.Run([]string{"test", "tasks", "seed"}); err != nil {
		t.Fatal(err)
	}
	if !ran {
		t.Error("task action did not run")
	}
}
