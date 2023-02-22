package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/db"
	"github.com/snehil-sinha/goBookStore/service/handlers"
)

// Used to start the service
func Start(s *common.App) {

	var err error

	// Flush the buffered logs (if any) after successfully starting the service
	defer s.Log.Core().Sync()

	err = db.New(s.Log, s.Cfg.GoBookStore.DB, s.Cfg.GoBookStore.URI)
	if err != nil {
		s.Log.Fatal(err.Error())
	} else {
		s.Log.Info("successfully initialized the GoBookStore")
	}

	s.Log.Sugar().Infof("starting HTTP listeners [%s:%s]", s.Cfg.Bind, s.Cfg.Port)

	r := gin.New()

	logger := s.Log.Logger
	r.Use(LoggerWithConfig(logger, &HTTPLogCfg{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{},
	}))
	r.Use(gin.Recovery())
	r.GET("/health", handlers.PingHandler()) // health check

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", handlers.PingHandler())
		v1.GET("/books", handlers.FindBooksHandler(s))
		v1.GET("/books/:id", handlers.FindBookHandler(s))
		v1.POST("/books", handlers.CreateBookHandler(s))
		v1.PUT("/books/:id", handlers.UpdateBookHandler(s))
		v1.DELETE("/books/:id", handlers.DeleteBookHandler(s))

	}

	if err := r.Run(s.Cfg.Bind + ":" + s.Cfg.Port); err != nil {
		s.Log.Fatal(err.Error())
	}

}
