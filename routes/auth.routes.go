package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/s1Sharp/s1-tts-restapi/controller"
	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/middleware"
	"github.com/s1Sharp/s1-tts-restapi/service"
)

type AuthRouteController struct {
	authController controller.AuthController
}

func NewAuthRouteController(authController controller.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup, userService service.UserService, config config.Config) {
	router := rg.Group("/auth")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(userService, config), rc.authController.LogoutUser)
	router.GET("/verify/email/:verificationCode", rc.authController.VerifyEmail)
	router.POST("/forgotPassword", rc.authController.ForgotPassword)
}
