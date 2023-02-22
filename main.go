package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/service"
)

// Driver function
func main() {
	fmt.Println("GoLang Book Store")

	var (
		err        error
		configFile string = "config.yaml"
	)

	// load config from yaml file
	cfg, err := common.LoadConfig(configFile)
	if err != nil {
		fmt.Println("error loading config: ", err)
		os.Exit(-1)
	}

	// loading go .env file
	er := godotenv.Load()
	if er != nil {
		fmt.Println("error loading .env file")
		os.Exit(-1)
	}

	// loading env specific variable
	err = common.LoadEnvSpecificConfigVariables(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// instantiating the logger
	log, err := common.NewLogger(cfg.Env, cfg.GoBookStore.LOGPATH)
	if err != nil {
		fmt.Println("error instantiating the logger instance: ", err)
		os.Exit(-1)
	}

	s := &common.App{
		Cfg: cfg,
		Log: log,
	}

	// start the service
	service.Start(s)
}
