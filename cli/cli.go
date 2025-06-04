package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/tessaris/tessaris/app"
	"github.com/tessaris/tessaris/config"
	"github.com/tessaris/tessaris/router"
	"github.com/tessaris/tessaris/version"
)

type Cli struct {
	app *app.App
}

func New(routes router.Routes) *Cli {
	return &Cli{app.New(config.New(), routes)}
}

func (c *Cli) Run() {
	code := c.parseCmd(os.Stdin, os.Stdout, os.Stderr, os.Args)

	if code != 0 {
		os.Exit(code)
	}
}

func (c *Cli) serveCmd(prod bool) int {
	r := router.New(prod, c.app.Config)
	r.Serve(c.app.Routes)

	return 0
}

func (c *Cli) routesCmd() int {
	r := router.New(false, c.app.Config)
	r.ListRoutes(c.app.Routes)

	return 0
}

const hintText = `Usage: ./tessaris <command> [<args>...]

Tessaris - batteries included framework for Go

See docs at https://tessaris.dev/daocs

Commands:
  help       Displays information about the tessaris environment
  version    Prints the version

  serve      Starts the development server
  prod       Starts the production server

  routes     List the routes
`

func (c *Cli) parseCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(stderr, hintText)
		return 0
	}

	switch args[1] {
	case "serve":
		return c.serveCmd(false)
	case "prod":
		return c.serveCmd(true)
	case "make":
		return c.makeCmd(stdin, stdout, stderr, args[2:])
	case "routes":
		return c.routesCmd()
	case "version", "--version":
		fmt.Fprintln(stdout, version.Version())
		return 0
	case "help", "-help", "--help", "-h":
		fmt.Fprint(stdout, hintText)
		return 0
	default:
		fmt.Fprintf(stderr, "Unknown command: %s\n", args[1])
		return 1
	}
}
