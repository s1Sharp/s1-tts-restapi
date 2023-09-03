package main

import (
	"context"
	"github.com/s1Sharp/s1-tts-restapi/controller"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/routes"
	"github.com/s1Sharp/s1-tts-restapi/service"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ttsServer struct {
	taskService         *service.TaskService
	TaskController      *controller.TaskController
	TaskRouteController *routes.TaskRouteController

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

	defer func() {
		// db disconnected
		if err := s.mongoStorage.Client.Disconnect(s.ctx); err != nil {
			log.Fatal(err)
		}
	}()

	return s.server.Run(":" + s.cfg.HttpAddr)
}
