package main

import (
	"github.com/gin-contrib/cors"
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

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController
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

	PostController = controllers.NewPostController(config.DB)
	PostRouteController = routes.NewRoutePostController(PostController)
	server = gin.Default()
}

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", conf.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/check", func(c *gin.Context) {
		message := "Testing connection"
		c.JSON(http.StatusOK, gin.H{"message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)

	log.Fatal(server.Run(":" + conf.ServerPort))

}
