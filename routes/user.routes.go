package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/s1Sharp/s1-tts-restapi/controller"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/middleware"
	"github.com/s1Sharp/s1-tts-restapi/service"
)

type UserRouteController struct {
	userController controller.UserController
}

func NewRouteUserController(userController controller.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService service.UserService, config config.Config) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(userService, config))
	router.GET("/me", uc.userController.GetMe)
}
