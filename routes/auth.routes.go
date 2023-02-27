package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/auth-with-gorm-postgres/controllers"
	"github.com/meshachdamilare/auth-with-gorm-postgres/middleware"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (ac *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/register", ac.authController.SignUpUser)
	router.POST("/login", ac.authController.SignInUser)
	router.GET("/refresh", ac.authController.RefreshAccessToken)
	router.GET("/logout", middleware.AuthMiddleware, ac.authController.SignOutUser)
	router.GET("/verifyemail/:verificationCode", ac.authController.VerifyEmail)
	router.POST("/forgotpassword", ac.authController.ForgotPassword)
	router.PATCH("/resetpassword/:resetToken", ac.authController.ResetPassword)
}
