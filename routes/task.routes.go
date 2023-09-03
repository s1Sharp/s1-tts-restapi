package routes

import (
	"s1-tts-restapi/controller"
	"s1-tts-restapi/middleware"

	"github.com/gin-gonic/gin"
)

type TaskRouteController struct {
	taskController controller.TaskController
}

func NewTaskControllerRoute(taskController controller.TaskController) TaskRouteController {
	return TaskRouteController{taskController}
}

func (r *TaskRouteController) TaskRoute(rg *gin.RouterGroup) {
	router := rg.Group("/tasks")

	router.GET("/find", r.taskController.GetUserTasks, middleware.DeserializeTasks())
	router.GET("/:taskId", r.taskController.TaskById)
	router.POST("/", r.taskController.CreateTask)
}
