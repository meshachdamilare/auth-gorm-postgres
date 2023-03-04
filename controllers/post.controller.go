package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/meshachdamilare/auth-with-gorm-postgres/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(DB *gorm.DB) PostController {
	return PostController{DB: DB}
}

func (pc *PostController) CreatePost(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	var payload *models.CreatePostRequest
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newPost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		User:      currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	result := pc.DB.Create(&newPost)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "sucess", "data": newPost})

}

func (pc *PostController) UpdatePost(c *gin.Context) {
	postId := c.Param("postId")
	currentUser := c.MustGet("currentUser").(models.User)

	var payload *models.UpdatePost
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var oldPost models.Post
	result := pc.DB.First(&oldPost, "id = ?", postId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}
	now := time.Now()
	updatePost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		User:      currentUser.ID,
		CreatedAt: oldPost.CreatedAt,
		UpdatedAt: now,
	}
	pc.DB.Model(&oldPost).Updates(&updatePost)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": oldPost})
}

func (pc *PostController) FindPostById(c *gin.Context) {
	postId := c.Param("postId")
	var post models.Post
	result := pc.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

func (pc *PostController) FindPosts(c *gin.Context) {
	var page = c.DefaultQuery("page", "1")
	var limit = c.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var posts []models.Post
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&posts)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}

func (pc *PostController) DeletePost(c *gin.Context) {
	postId := c.Param("postId")

	result := pc.DB.Delete(&models.Post{}, "id = ?", postId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
