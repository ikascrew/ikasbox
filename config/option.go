package config

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
