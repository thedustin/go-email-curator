package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/thedustin/go-email-curator/cmd/currator/cmd"
	"github.com/thedustin/go-email-curator/cmd/currator/currate"
)

var (
	Name         = "go-email-currator"
	ReadableName = "Go E-Mail Currator"
	Version      = "0.0.0"
)

var cli struct {
	Currate currate.CurrateCmd `cmd:""`

	Debug   bool `help:"Enables debug mode and increases logging"`
	Version bool `help:"Display the application version" short:"V"`
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)

	ctx := kong.Parse(&cli)

	if cli.Debug {
		logger.Printf("Cli parameters: %+v", cli)
	}

	if cli.Version {
		fmt.Printf("%s v%s\n", ReadableName, Version)
		return
	}

	err := ctx.Run(&cmd.Context{
		Debug:  cli.Debug,
		Logger: logger,
	})

	ctx.FatalIfErrorf(err)
}
