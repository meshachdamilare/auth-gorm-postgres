package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/auth-with-gorm-postgres/controllers"
	"github.com/meshachdamilare/auth-with-gorm-postgres/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/users")
	router.GET("/me", middleware.AuthMiddleware, uc.userController.GetMe)
}
