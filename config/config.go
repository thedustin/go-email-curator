package config

import (
	"net/url"
	"sync"
)

type config struct {
	Server   *url.URL
	Username string
	Password string

	Filters []Filter
}

type Filter struct {
	Query   string
	Actions []string
}

var instance config
var once sync.Once

func Get() *config {
	once.Do(func() {
		instance = config{} // <-- thread safe
	})

	return &instance
}
