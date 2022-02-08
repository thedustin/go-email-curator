package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/alecthomas/kong"
)

var (
	Name         = "go-email-currator"
	ReadableName = "Go E-Mail Currator"
	Version      = "0.0.0"
)

type MailServerCmd struct {
	MailServer   *url.URL `required arg:"" env:"CURRATOR_MAIL_SERVER"`
	MailUsername string   `help:"" env:"CURRATOR_MAIL_USERNAME"`
	MailPassword string   `help:"" env:"CURRATOR_MAIL_PASSWORD"`
}

var cli struct {
	Debug   bool `help:"Enables debug mode and increases logging"`
	Version bool `help:"Display the application version" short:"V"`
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)

	kong.Parse(&cli)

	if cli.Debug {
		logger.Printf("Cli parameters: %+v", cli)
	}

	if cli.Version {
		fmt.Printf("%s v%s\n", ReadableName, Version)
		return
	}
}
