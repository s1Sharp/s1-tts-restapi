package routes

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/s1Sharp/s1-tts-restapi/controller"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/logger"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log = logger.GetLogger()
)

type TtsServer struct {
	taskService         *service.TaskService
	TaskController      *controller.TaskController
	TaskRouteController *TaskRouteController

	userService         service.UserService
	UserController      *controller.UserController
	UserRouteController *UserRouteController

	authService         *service.AuthService
	AuthController      *controller.AuthController
	AuthRouteController *AuthRouteController

	mongoStorage *storage.MongoStorage
	redisStorage *storage.RedisStorage

	ctx    context.Context
	server *gin.Engine

	cfg *config.Config

	quitChan *chan os.Signal
}

func (s *TtsServer) GetRouter() *gin.Engine {
	return s.server
}

func (s *TtsServer) addHealthcheck() error {
	var healthcheckMessage string
	var err error
	if s.redisStorage != nil {
		healthcheckMessage, err = s.redisStorage.Client.Get("healthcheck").Result()
		if err != nil {
			log.Printf("Failed to fetch healthchek message form redis: %s", err)
		}
	} else {
		healthcheckMessage = s.cfg.HealthcheckMessage
	}

	s.server.Group("").GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": healthcheckMessage})
	})
	return err
}

func (s *TtsServer) initCors() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{s.cfg.ClientOrigin}
	corsConfig.AllowCredentials = true
	s.server.Use(cors.New(corsConfig))
}

func (s *TtsServer) Close() {
	*s.quitChan <- syscall.SIGUSR1
}

func (s *TtsServer) Run() error {

	if s.mongoStorage != nil {
		defer func() {
			// db disconnected
			if err := s.mongoStorage.Client.Disconnect(s.ctx); err != nil {
				log.Fatal(err)
			}
		}()
	}

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
	s.quitChan = &quit

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

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

func NewServer(ctx context.Context, cfg config.Config) *TtsServer {
	mongoStorage := storage.NewMongoStorage(cfg.MongoUrl, ctx)
	redisStorage := storage.NewRedisStorage(cfg.RedisUrl, cfg.HealthcheckMessage, ctx)

	taskService := service.NewTaskService(&mongoStorage, ctx)
	taskController := controller.NewTaskController(taskService)
	taskRouteController := NewTaskControllerRoute(taskController)

	userService := service.NewUserServiceImpl(&mongoStorage, ctx)
	userController := controller.NewUserController(userService)
	userRouteController := NewRouteUserController(userController)

	authService := service.NewAuthService(&mongoStorage, ctx, &cfg)
	authController := controller.NewAuthController(authService, userService, ctx, &mongoStorage, &cfg)
	authRouteController := NewAuthRouteController(authController)

	ginServer := gin.Default()
	s := TtsServer{
		ctx:          ctx,
		mongoStorage: &mongoStorage,
		redisStorage: &redisStorage,

		taskService:         &taskService,
		TaskController:      &taskController,
		TaskRouteController: &taskRouteController,

		userService:         userService,
		UserController:      &userController,
		UserRouteController: &userRouteController,

		authService:         &authService,
		AuthController:      &authController,
		AuthRouteController: &authRouteController,

		cfg:    &cfg,
		server: ginServer,
	}
	s.initServerPrivate()
	return &s
}

type TestServerInternal struct {
	mongoStorage storage.MongoStorage
	redisStorage storage.RedisStorage

	taskService         service.TaskService
	taskController      controller.TaskController
	taskRouteController TaskRouteController

	userService         service.UserService
	userController      controller.UserController
	userRouteController UserRouteController

	authService         service.AuthService
	authController      controller.AuthController
	authRouteController AuthRouteController
}

func NewTestServer(ctx context.Context, cfg config.Config, ti TestServerInternal) *TtsServer {
	var mongoStorage *storage.MongoStorage = nil
	var redisStorage *storage.RedisStorage = nil

	var taskService service.TaskService
	if ti.taskService != nil {
		taskService = ti.taskService
	} else {
		taskService = service.NewTaskService(mongoStorage, ctx)
	}
	taskController := controller.NewTaskController(taskService)
	taskRouteController := NewTaskControllerRoute(taskController)

	var userService service.UserService
	if ti.userService != nil {
		userService = ti.userService
	} else {
		userService = service.NewUserServiceImpl(mongoStorage, ctx)
	}
	userController := controller.NewUserController(userService)
	userRouteController := NewRouteUserController(userController)

	var authService service.AuthService
	if ti.userService != nil {
		authService = ti.authService
	} else {
		authService = service.NewAuthService(mongoStorage, ctx, &cfg)
	}
	authController := controller.NewAuthController(authService, userService, ctx, mongoStorage, &cfg)
	authRouteController := NewAuthRouteController(authController)

	ginServer := gin.Default()
	gin.SetMode(gin.TestMode)

	s := TtsServer{
		ctx:          ctx,
		mongoStorage: mongoStorage,
		redisStorage: redisStorage,

		taskService:         &taskService,
		TaskController:      &taskController,
		TaskRouteController: &taskRouteController,

		userService:         userService,
		UserController:      &userController,
		UserRouteController: &userRouteController,

		authService:         &authService,
		AuthController:      &authController,
		AuthRouteController: &authRouteController,

		cfg:    &cfg,
		server: ginServer,
	}
	s.initServerPrivate()
	return &s
}

func (s *TtsServer) initServerPrivate() {
	err := s.addHealthcheck()
	if err != nil {
		log.Error("error add healthcheck")
	}

	s.initCors()

	router := s.server.Group("/api/v1")
	s.TaskRouteController.TaskRoute(router)
	s.UserRouteController.UserRoute(router, s.userService, *s.cfg)
	s.AuthRouteController.AuthRoute(router, s.userService, *s.cfg)
}
