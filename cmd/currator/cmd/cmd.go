package cmd

import "log"

type Context struct {
	Debug  bool
	Logger *log.Logger
}
