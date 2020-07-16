package config

import (
	"golang.org/x/xerrors"
)

type Config struct {
	ContentPath  string
	DatabasePath string

	Host string
	Port int

	Debug bool
}

var gConfig *Config = nil

func defaultConfig() *Config {
	conf := Config{}
	conf.ContentPath = "/tmp"
	conf.DatabasePath = "ikasbox.db"
	conf.Host = ""
	conf.Port = 5555
	conf.Debug = false
	return &conf
}

func Set(opts ...Option) error {

	if opts == nil {
		return nil
	}

	gConfig = defaultConfig()

	for _, opt := range opts {
		err := opt(gConfig)
		if err != nil {
			return xerrors.Errorf("option error: %w", err)
		}
	}

	return nil
}

func Get() *Config {
	if gConfig == nil {
		gConfig = defaultConfig()
	}
	return gConfig
}
