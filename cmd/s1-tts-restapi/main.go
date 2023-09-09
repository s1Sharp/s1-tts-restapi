package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/s1Sharp/s1-tts-restapi/controller"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/logger"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/routes"
	"github.com/s1Sharp/s1-tts-restapi/service"
)

var (
	log = logger.GetLogger()
)

func main() {
	defer log.Close()

	// read config from env
	cfg := config.ReadEnv()
	ctx := context.TODO()

	mongoStorage := storage.NewMongoStorage(cfg.MongoUrl, ctx)
	redisStorage := storage.NewRedisStorage(cfg.RedisUrl, cfg.HealthcheckMessage, ctx)

	taskService := service.NewTaskService(&mongoStorage, ctx)
	taskController := controller.NewTaskController(taskService)
	taskRouteController := routes.NewTaskControllerRoute(taskController)

	userService := service.NewUserServiceImpl(&mongoStorage, ctx)
	userController := controller.NewUserController(userService)
	userRouteController := routes.NewRouteUserController(userController)

	authService := service.NewAuthService(&mongoStorage, ctx, &cfg)
	authController := controller.NewAuthController(authService, userService, ctx, &mongoStorage, &cfg)
	authRouteController := routes.NewAuthRouteController(authController)

	ginServer := gin.Default()
	s1TtsServer := ttsServer{
		ctx:          ctx,
		mongoStorage: &mongoStorage,
		redisStorage: &redisStorage,

		taskService:         &taskService,
		TaskController:      &taskController,
		TaskRouteController: &taskRouteController,

		userService:         &userService,
		UserController:      &userController,
		UserRouteController: &userRouteController,

		authService:         &authService,
		AuthController:      &authController,
		AuthRouteController: &authRouteController,

		cfg:    &cfg,
		server: ginServer,
	}
	if err := s1TtsServer.Run(); err != nil {
		log.Error(err)
	}

	os.Exit(0)
}
