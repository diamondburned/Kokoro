package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"

	"github.com/Gigamons/Kokoro/constants"
	"github.com/Gigamons/Kokoro/handler"
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

	handler.CLIENT = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%v", conf.Redis.Hostname, conf.Redis.Port),
	})

	if _, err := handler.CLIENT.Ping().Result(); err != nil {
		logger.Error("Could not connect to Redis Server!")
		return
	}
	defer handler.CLIENT.Close()

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", os.ModePerm)
	}

	if _, err := os.Stat("data/screenshots"); os.IsNotExist(err) {
		os.Mkdir("data/screenshots", os.ModePerm)
	}

	os.Setenv("DEBUG", strconv.FormatBool(conf.Server.Debug))
	os.Setenv("CHEESEGULL", conf.CheeseGull.APIUrl)

	server.StartServer(conf.Server.Hostname, conf.Server.Port)
}
