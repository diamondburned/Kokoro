package constants

import (
	"github.com/Gigamons/common/consts"
)

// Config is for the config.yml file
type Config struct {
	Server struct {
		Hostname string
		Port     int16
		Debug    bool
	}
	CheeseGull struct {
		APIUrl string
	}
	MySQL consts.MySQLConf
	Redis struct {
		Hostname string
		Port     int16
	}
}
