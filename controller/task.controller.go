package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/s1Sharp/s1-tts-restapi/internal/models"
	"github.com/s1Sharp/s1-tts-restapi/service"
	"net/http"
	"strconv"
	"strings"
)

type TaskController interface {
	CreateTask(ctx *gin.Context)
	TaskById(ctx *gin.Context)
	GetUserTasks(ctx *gin.Context)
}

type TaskControllerImpl struct {
	postService service.TaskService
}

func NewTaskController(taskService service.TaskService) TaskController {
	return &TaskControllerImpl{taskService}
}

func (pc *TaskControllerImpl) CreateTask(ctx *gin.Context) {
	var task *models.CreateTaskScheme

	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newTask, err := pc.postService.CreateTask(task)

	if err != nil {
		if strings.Contains(err.Error(), "task already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"task_id": newTask.ID}})
}

func (pc *TaskControllerImpl) TaskById(ctx *gin.Context) {
	taskId := ctx.Param("id")

	task, err := pc.postService.TaskById(taskId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": task})
}

func (pc *TaskControllerImpl) GetUserTasks(ctx *gin.Context) {
	var user, ok = ctx.GetQuery("user")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "user required"})
		return
	}
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	tasks, err := pc.postService.GetUserTasks(user, intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.Set("tasks", tasks)
	ctx.Next()
}
