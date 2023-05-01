package main

import (
	"io"
	"os"

	"github.com/arpitchauhan/simple-database/cmd"
)

func main() {
	run(os.Stdout, os.Stderr)
}

func run(wout io.Writer, werr io.Writer) {
	cmd.Execute(wout, werr)
}
