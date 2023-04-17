package test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"

	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/db"
	"github.com/snehil-sinha/goBookStore/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var BaseUrl string

const testConfigPath = "/Users/snehil.sinha/Documents/bookstore/config.yaml"

func LoadTestConfig() (cfg *common.Config, err error) {
	// load config from yaml file
	cfg, err = common.LoadConfig(testConfigPath)
	if err != nil {
		err = fmt.Errorf("failed to load the test config file: %s", err)
		return
	}
	if cfg.Env != "test" {
		err = fmt.Errorf("error: Unable to run tests. Please ensure that the APP_ENV environment variable is set to test")
		return
	}
	return
}

func GetMockTestLogger(cfg *common.Config) (*common.Logger, error) {
	log, err := common.NewLogger(cfg.Env, cfg.GoBookStore.LOGPATH)
	return log, err
}

func GetMockAppInstance(cfg *common.Config) (s *common.App, err error) {
	// instantiate the logger
	log, err := GetMockTestLogger(cfg)
	if err != nil {
		err = fmt.Errorf("failed to instantiate the test logger: %s", err)
		return
	}

	// instantiate the App struct
	s = &common.App{
		Cfg: cfg,
		Log: log,
	}
	return
}

func StartTestSever(s *common.App) (srv *http.Server, err error) {
	return service.Start(s), err
}

func SetupServerShutdown(s *common.App, server *http.Server) {
	service.WaitForShutdown()
	service.GracefullyShutDownServer(s.Log, server)
}

func SignalShutDown(t *testing.T) {
	// Get the process ID of the target process
	pid := os.Getpid()

	// Send a SIGTERM signal to the target process
	err := syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		t.Fatalf("server shutdown failed: %s", err)
	}
}

func ClearDB(ctx context.Context) error {
	filter := bson.D{}
	_, err := db.GoBookStore.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to clear the DB, %s", err)
	}
	return nil
}

func CloseDBConnection(c *mongo.Client, ctx context.Context) error {
	err := c.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to close the DB connection, %s", err)
	}
	return nil
}
