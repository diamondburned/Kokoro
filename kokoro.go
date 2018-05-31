package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Gigamons/Kokoro/constants"
	"github.com/Gigamons/Kokoro/server"
	"github.com/Gigamons/common/consts"
	"github.com/Gigamons/common/helpers"
	"github.com/Gigamons/common/logger"
)

func init() {
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		helpers.CreateConfig("config", constants.Config{MySQL: consts.MySQLConf{Database: "gigamons", Hostname: "localhost", Port: 3306, Username: "root"}})
		fmt.Println("I've just created a config.yml! please edit!")
		os.Exit(0)
	}
}

func main() {
	var err error
	var conf constants.Config

	helpers.GetConfig("config", &conf)
	helpers.Connect(conf.MySQL)
	if err = helpers.DB.Ping(); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer helpers.DB.Close()

	os.Setenv("DEBUG", strconv.FormatBool(conf.Server.Debug))
	os.Setenv("CHEESEGULL", conf.CheeseGull.APIUrl)

	server.StartServer(conf.Server.Hostname, conf.Server.Port)
}
