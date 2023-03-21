package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/db"
	"github.com/snehil-sinha/goBookStore/service"
)

const usage = `Usage:
%s  -c | -config /path/to/config -b | -bind <ip address> -p | -port <port number>
`

// Driver function
func main() {
	fmt.Println("GoLang Book Store")

	var (
		err        error
		configFile string
		bind       string
		port       string
	)

	// loading go .env file
	er := godotenv.Load()
	if er != nil {
		fmt.Println("error loading .env file")
		os.Exit(-1)
	}

	// Setup the flags
	flag.Usage = func() { // [1]
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&configFile, "c", "", "Configuration file path.")
	flag.StringVar(&configFile, "config", "", "Configuration file path.")
	flag.StringVar(&bind, "b", common.Bind, "IP address to bind")
	flag.StringVar(&bind, "bind", common.Bind, "IP address to bind")
	flag.StringVar(&port, "p", common.Port, "port number to listen")
	flag.StringVar(&port, "port", common.Port, "port number to listen")

	flag.Parse()

	if configFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	// load config from yaml file
	cfg, err := common.LoadConfig(configFile)
	if err != nil {
		fmt.Println("error loading config: ", err)
		os.Exit(-1)
	}

	//set APP_ENV in .env file when running in local
	if env := common.GetAppEnv(); env != "" {
		cfg.Env = env
	} else {
		fmt.Println("service environment not found")
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
		fmt.Println("error instantiating the logger: ", err)
		os.Exit(-1)
	}

	s := &common.App{
		Cfg: cfg,
		Log: log,
	}
	cfg.Bind = bind
	cfg.Port = port

	// start the service
	server := service.Start(s)
	// wait for a signal to shutdown server
	service.WaitForShutdown()
	// gracefully shutdown the server
	service.GracefullyShutDownServer(s.Log, server)
	// close the DB connection
	db.Client.Close(context.TODO(), log)
}
