package main

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/auth-with-gorm-postgres/config"
	"github.com/meshachdamilare/auth-with-gorm-postgres/controllers"
	"github.com/meshachdamilare/auth-with-gorm-postgres/routes"

	"log"
	"net/http"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController
)

func init() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("? Could not load environment variables", err)
	}
	config.ConnectDB(&conf)

	AuthController = controllers.NewAuthController(config.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(config.DB)
	UserRouteController = routes.NewRouteUserController(UserController)
	server = gin.Default()
}

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("? Could not load environment variables", err)
	}
	router := server.Group("/api")
	router.GET("/check", func(c *gin.Context) {
		message := "Testing connection"
		c.JSON(http.StatusOK, gin.H{"message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)

	log.Fatal(server.Run(":" + conf.ServerPort))

}
