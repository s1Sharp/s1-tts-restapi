package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/s1Sharp/s1-tts-restapi/controller"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/routes"
	"github.com/s1Sharp/s1-tts-restapi/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ttsServer struct {
	taskService         *service.TaskService
	TaskController      *controller.TaskController
	TaskRouteController *routes.TaskRouteController

	userService         *service.UserService
	UserController      *controller.UserController
	UserRouteController *routes.UserRouteController

	authService         *service.AuthService
	AuthController      *controller.AuthController
	AuthRouteController *routes.AuthRouteController

	mongoStorage *storage.MongoStorage
	redisStorage *storage.RedisStorage

	ctx    context.Context
	server *gin.Engine

	cfg *config.Config
}

func (s ttsServer) addHealthcheck() error {
	healthcheckMessage, err := s.redisStorage.Client.Get("healthcheck").Result()
	if err != nil {
		log.Printf("Failed to fetch healthchek message form redis: %s", err)
		healthcheckMessage = s.cfg.HealthcheckMessage
	}
	s.server.Group("").GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": healthcheckMessage})
	})
	return err
}

func (s ttsServer) initCors() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{s.cfg.ClientOrigin}
	corsConfig.AllowCredentials = true
	s.server.Use(cors.New(corsConfig))
}

func (s ttsServer) Run() error {
	err := s.addHealthcheck()
	if err != nil {
		return err
	}

	s.initCors()

	router := s.server.Group("/api/v1")
	s.TaskRouteController.TaskRoute(router)
	s.UserRouteController.UserRoute(router, *s.userService, *s.cfg)
	s.AuthRouteController.AuthRoute(router, *s.userService, *s.cfg)

	defer func() {
		// db disconnected
		if err := s.mongoStorage.Client.Disconnect(s.ctx); err != nil {
			log.Fatal(err)
		}
	}()

	srv := &http.Server{
		Addr:    ":" + s.cfg.HttpAddr,
		Handler: s.server,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
		return err
	}
	return nil
}
