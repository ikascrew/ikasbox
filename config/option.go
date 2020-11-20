package config

import (
	"fmt"
	"strings"
)

type Option func(*Config) error

func LocalOnly() Option {
	return func(conf *Config) error {
		conf.Host = "localhost"
		return nil
	}
}

func Path(p string) Option {
	return func(conf *Config) error {
		conf.DatabasePath = p
		return nil
	}
}

func Port(p int) Option {
	return func(conf *Config) error {
		conf.Port = p
		return nil
	}
}

func Argument(args []string) Option {
	return func(conf *Config) error {
		if len(args) < 1 {
			return fmt.Errorf("ikasbox subcommand argument required.")
		}

		conf.SubCommand = args[0]

		if conf.SubCommand != "start" && conf.SubCommand != "init" {
			if len(args) < 2 {
				return fmt.Errorf("ikasbox [%s] command argument required.")
			}
			conf.Function = args[1]
			if len(args) > 2 {
				conf.Arguments = args[2:]
			}
		}
		return nil
	}
}

func Extension(exts string) Option {
	return func(conf *Config) error {
		conf.Extensions = strings.Split(exts, ",")
		return nil
	}
}
