package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/s1Sharp/s1-tts-restapi/internal/models"
	"github.com/s1Sharp/s1-tts-restapi/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}
