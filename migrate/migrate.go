package main

import (
	"fmt"
	"github.com/meshachdamilare/auth-with-gorm-postgres/config"
	"github.com/meshachdamilare/auth-with-gorm-postgres/models"
	"log"
)

func init() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}
	config.ConnectDB(&conf)
}

func main() {
	config.DB.AutoMigrate(&models.User{})
	fmt.Println(" Migration complete.")
}
