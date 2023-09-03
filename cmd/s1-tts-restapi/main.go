package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"s1-tts-restapi/controller"
	"s1-tts-restapi/internal/config"
	"s1-tts-restapi/internal/storage"
	"s1-tts-restapi/routes"
	"s1-tts-restapi/service"
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
