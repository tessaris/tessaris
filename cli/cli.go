package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/tesseris-go/tesseris"
	"github.com/tesseris-go/tesseris/router"
)

type Cli struct {
	t *tesseris.Tesseris
}

func New(t *tesseris.Tesseris) *Cli {
	return &Cli{t}
}

func (c *Cli) Run() {
	code := c.parseCmd(os.Stdin, os.Stdout, os.Stderr, os.Args)

	if code != 0 {
		os.Exit(code)
	}
}

func (c *Cli) serveCmd(prod bool) int {
	r := router.New(prod)
	r.Serve(c.t.Routes)

	return 0
}

const hintText = `Usage: ./tesseris <command> [<args>...]

Tesseris - batteries included framework for Go

See docs at https://tesseris.iskandervdh.nl/docs

Commands:
  help       Displays information about the tesseris environment
  version    Prints the version
  serve      Starts the development server
  prod       Starts the production server

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
	case "version", "--version":
		fmt.Fprintln(stdout, tesseris.Version())
		return 0
	case "help", "-help", "--help", "-h":
		fmt.Fprint(stdout, hintText)
		return 0
	default:
		fmt.Fprintf(stderr, "Unknown command: %s\n", args[1])
		return 1
	}
}
