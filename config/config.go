package config

type Config struct {
	ContentPath  string
	DatabasePath string

	Host string
	Port string
}

var gConfig *Config = nil

func setConfig() {
	gConfig = &Config{}
	gConfig.ContentPath = "/tmp"
	gConfig.DatabasePath = "ikasbox.db"
	gConfig.Host = "localhost"
	gConfig.Port = "5555"

}

func Get() *Config {
	if gConfig == nil {
		setConfig()
	}
	return gConfig
}
