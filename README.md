# Webapp

[![GoDoc](https://godoc.org/github.com/sfreiberg/webapp?status.png)](https://godoc.org/github.com/sfreiberg/webapp)

## Warning

This repository is experimental and opinionated. You should not rely on it for
production unless you are prepared to keep up with breaking changes.

## About

There is a lot of common code that goes with building web applications in Go.
While this is fine when that is your primary job it can make building hobby
sites or just experimenting much slower and more tedious. Webapp attempts to
improve that situation but this is an experiment. Also, the code will be
opinionated and my opinions may not match yours.

## Example

```
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sfreiberg/webapp"
)

func main() {
	app := webapp.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	if err := app.Run(); err != nil {
		log.Fatalf("Error starting web app: %s\n", err)
	}
}
```

## Usage

```
NAME:
   your-app - Built with github.com/sfreiberg/webapp

USAGE:
   your-app [global options] command [command options] [arguments...]

COMMANDS:
   server, s    Start web server
   generate, g  Generators
   routes       Display routes
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```