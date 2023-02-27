package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/auth-with-gorm-postgres/config"
	"github.com/meshachdamilare/auth-with-gorm-postgres/models"
	"github.com/meshachdamilare/auth-with-gorm-postgres/utils"
	"net/http"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	var acces_token string
	cookie, err := c.Cookie("access_token")
	if err != nil {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Authorization header empty.",
			})
			return
		}
		fields := strings.Fields(authHeader)
		if len(fields) != 0 && fields[0] == "Bearer" {
			acces_token = fields[1]
		}
	} else if err == nil {
		acces_token = cookie
	}

	if acces_token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": "fail", "message": "you are not logged in",
		})
		return
	}

	conf, _ := config.LoadConfig(".")

	// Recall, payload stored is userID
	payload, err := utils.ValidateToken(acces_token, conf.AccessTokenPublicKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var user models.User
	result := config.DB.First(&user, "id =?", fmt.Sprint(payload))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}
	c.Set("currentUser", user)
	c.Next()

}
