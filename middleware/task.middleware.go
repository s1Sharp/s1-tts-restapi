package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/s1Sharp/s1-tts-restapi/internal/models"
)

func DeserializeTasks() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tasks, ok := ctx.Get("tasks")
		if !ok {
			log.Println("failed to get middleware result")
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "server error"})
			return
		}

		taskResponse := models.DbTaskToTaskResponse(tasks.([]*models.DBTaskScheme))

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(taskResponse), "data": taskResponse})
	}
}
