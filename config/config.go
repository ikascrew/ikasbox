package config

type Config struct {
	ContentPath  string
	DatabasePath string

	Host string
	Port string

	Debug bool
}

var gConfig *Config = nil

func setConfig() {
	gConfig = &Config{}
	gConfig.ContentPath = "/tmp"
	gConfig.DatabasePath = "ikasbox.db"
	gConfig.Host = ""
	gConfig.Port = "5555"
	gConfig.Debug = false
}

func Get() *Config {
	if gConfig == nil {
		setConfig()
	}
	return gConfig
}
