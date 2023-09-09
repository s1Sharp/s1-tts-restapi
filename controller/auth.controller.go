package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/models"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/service"
	"github.com/s1Sharp/s1-tts-restapi/utils"

	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService service.AuthService
	userService service.UserService
	ctx         context.Context
	storage     *storage.MongoStorage
	config      *config.Config
}

func NewAuthController(authService service.AuthService, userService service.UserService, ctx context.Context, storage *storage.MongoStorage, config *config.Config) AuthController {
	return AuthController{authService, userService, ctx, storage, config}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if ac.config.AllVerifiedByDefault {
		return
	}

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)
	log.Printf("verificationCode = %s\n", verificationCode)

	updateData := &models.UpdateInput{
		VerificationCode: verificationCode,
	}

	// Update User in Database
	_, err = ac.userService.UpdateUserById(newUser.ID.Hex(), updateData)
	if err != nil {
		log.Println(err)
	}

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       ac.config.ClientOrigin + config.VerifyEmailRoute + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, "verificationCode.html", *ac.config)
	if err != nil {
		err = ac.userService.RemoveUserById(newUser.ID.Hex())
		if err != nil {
			log.Printf("failed to remove user, need reset passord")
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}

	message := "We sent an email with a verification code to " + user.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not verified, please verify your email to login"})
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	// Generate Tokens
	accessToken, err := utils.CreateToken(ac.config.AccessTokenExpired, user.ID, ac.config.AccessPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(ac.config.RefreshTokenExpired, user.ID, ac.config.RefreshPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, ac.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, ac.config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", ac.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	sub, err := utils.ValidateToken(cookie, ac.config.RefreshPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(ac.config.AccessTokenExpired, user.ID, ac.config.AccessPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, ac.config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", ac.config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := ac.storage.UserCollection().UpdateOne(ac.ctx, query, update)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": "Could not verify email address"})
		return
	}

	fmt.Println(result)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})

}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var userCredential *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	user, err := ac.userService.FindUserByEmail(userCredential.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusOK, gin.H{"status": "fail", "message": message})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
		return
	}

	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)

	// Update User in Database
	query := bson.D{{Key: "email", Value: strings.ToLower(userCredential.Email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := ac.storage.UserCollection().UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}
	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       ac.config.ClientOrigin + "/forgotPassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	err = utils.SendEmail(user, &emailData, "resetPassword.html", *ac.config)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}
