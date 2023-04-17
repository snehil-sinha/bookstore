package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/snehil-sinha/goBookStore/common"
	"github.com/snehil-sinha/goBookStore/db"
	"github.com/snehil-sinha/goBookStore/models/book"
	"github.com/snehil-sinha/goBookStore/service/handlers"
)

func setGinMode(env string) {
	switch env {
	case "development":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "production":
		gin.SetMode(gin.ReleaseMode)
	}
}

// Used to start the service
func Start(s *common.App) *http.Server {

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

	setGinMode(s.Cfg.Env)

	r := gin.New()

	if gin.Mode() != gin.TestMode {
		logger := s.Log.Logger
		r.Use(LoggerWithConfig(logger, &HTTPLogCfg{
			TimeFormat: time.RFC3339,
			UTC:        true,
			SkipPaths:  []string{},
		}))
		r.Use(gin.Recovery())
	}

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return regexp.MustCompile(common.GetAllowedOriginsRegex()).MatchString(origin)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", handlers.PingHandler()) // health check

	bs := book.NewBookService()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", handlers.PingHandler())
		v1.GET("/books", handlers.FindBooksHandler(bs, s))
		v1.GET("/books/:id", handlers.FindBookHandler(bs, s))
		v1.POST("/books", handlers.CreateBookHandler(s))
		v1.PUT("/books/:id", handlers.UpdateBookHandler(s))
		v1.DELETE("/books/:id", handlers.DeleteBookHandler(s))

	}

	server := &http.Server{
		Addr:    s.Cfg.Bind + ":" + s.Cfg.Port,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Log.Sugar().Fatalf("Could not listen on %s:%s %v", s.Cfg.Bind, s.Cfg.Port, err)
		}
	}()
	return server
}

func WaitForShutdown() {
	// Wait for a signal to shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func GracefullyShutDownServer(log *common.Logger, server *http.Server) {

	log.Sugar().Info("Server shutting down...")

	// Gracefully shutdown the server, waiting for all active connections to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Sugar().Errorf("Error shutting down the server: %s", err)
	} else {
		log.Sugar().Info("Server stopped")
	}
}
