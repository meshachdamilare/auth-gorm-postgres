package main

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/auth-with-gorm-postgres/config"
	"log"
	"net/http"
)

var server *gin.Engine

func init() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("? Could not load environment variables", err)
	}
	config.ConnectDB(&conf)
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

	log.Fatal(server.Run(":" + conf.ServerPort))

}
