package cli

import (
	"fmt"
	"io"
)

const makeHintText = `Nothing to make specified. Please specify a something to make.

Things to make:
  controller
  model
  migration
  middleware
`

func makeControllerCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, "Please specify a controller name\n")

		return 1
	}

	return 0
}

func makeModelCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, "Please specify a model name\n")

		return 1
	}

	return 0
}

func makeMigrationCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, "Please specify a migration name\n")

		return 1
	}

	return 0
}

func makeMiddlewareCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, "Please specify a middleware name\n")

		return 1
	}

	return 0
}

func makeCmd(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	if len(args) < 1 {
		fmt.Fprint(stderr, makeHintText)

		return 1
	}

	switch args[0] {
	case "controller":
		return makeControllerCmd(stdin, stdout, stderr, args[1:])
	case "model":
		return makeModelCmd(stdin, stdout, stderr, args[1:])
	case "migration":
		return makeMigrationCmd(stdin, stdout, stderr, args[1:])
	case "middleware":
		return makeMiddlewareCmd(stdin, stdout, stderr, args[1:])
	default:
		fmt.Fprintf(stderr, "Unknown thing to make: %s\n", args[0])

		return 1
	}
}
