package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/auth-with-gorm-postgres/config"
	"github.com/meshachdamilare/auth-with-gorm-postgres/models"
	"github.com/meshachdamilare/auth-with-gorm-postgres/utils"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) SignUpUser(c *gin.Context) {
	var signUpPayload *models.SignUp

	if err := c.ShouldBindJSON(&signUpPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": err.Error(),
		})
		return
	}
	if signUpPayload.Password != signUpPayload.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "password do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(signUpPayload.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	now := time.Now()
	newUser := models.User{
		Name:      signUpPayload.Name,
		Email:     strings.ToLower(signUpPayload.Email),
		Password:  hashedPassword,
		Role:      "user",
		Verified:  true,
		Photo:     signUpPayload.Photo,
		Provider:  "local",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with email already exists"})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something bad happen"})
		return
	}

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Photo:     newUser.Photo,
		Role:      newUser.Role,
		Provider:  newUser.Provider,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (ac *AuthController) SignInUser(c *gin.Context) {
	var signInPayload *models.SignIn
	var user models.User

	if err := c.ShouldBindJSON(&signInPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail", "message": err.Error(),
		})
		return
	}

	result := ac.DB.First(&user, "email = ?", strings.ToLower(signInPayload.Email))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, signInPayload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	conf, _ := config.LoadConfig(".")

	payload := user.ID // Here we -> userID is the payload stored inside jwt with cookie

	accessToken, err := utils.CreateToken(conf.AccessTokenExpiresIn, payload, conf.AccessTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(conf.RefreshTokenExpiresIn, payload, conf.RefreshTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	fmt.Println(accessToken)
	fmt.Println(refreshToken)
	c.SetCookie("access_token", accessToken, conf.AccessTokenMaxAge, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, conf.RefreshTokenMaxAge, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", conf.AccessTokenMaxAge*60, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})

}

func (ac *AuthController) RefreshAccessToken(c *gin.Context) {
	message := "could not refresh the token"
	cookie, err := c.Cookie("refresh_token")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}
	conf, _ := config.LoadConfig(".")

	// payload stored in the cookie is userID

	payload, err := utils.ValidateToken(cookie, conf.RefreshTokenPublicKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(payload))
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}
	accessToken, err := utils.CreateToken(conf.AccessTokenExpiresIn, user.ID, conf.AccessTokenPrivateKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.SetCookie("access_token", accessToken, conf.AccessTokenMaxAge*60, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", conf.AccessTokenMaxAge*60, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) SignOutUser(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_in", "", -1, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
