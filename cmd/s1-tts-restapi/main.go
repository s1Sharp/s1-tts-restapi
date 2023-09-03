package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/s1Sharp/s1-tts-restapi/controller"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/routes"
	"github.com/s1Sharp/s1-tts-restapi/service"
	"log"
	"os"
)

func main() {
	// read config from env
	cfg := config.ReadEnv()
	ctx := context.TODO()

	mongoStorage := storage.NewMongoStorage(cfg.MongoUrl, ctx)
	redisStorage := storage.NewRedisStorage(cfg.RedisUrl, cfg.HealthcheckMessage, ctx)
	taskService := service.NewTaskService(&mongoStorage, ctx)
	taskController := controller.NewTaskController(taskService)
	taskRouteController := routes.NewTaskControllerRoute(taskController)
	ginServer := gin.Default()
	s1TtsServer := ttsServer{
		ctx:                 ctx,
		mongoStorage:        &mongoStorage,
		redisStorage:        &redisStorage,
		taskService:         &taskService,
		TaskController:      &taskController,
		TaskRouteController: &taskRouteController,
		cfg:                 &cfg,
		server:              ginServer,
	}
	if err := s1TtsServer.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
