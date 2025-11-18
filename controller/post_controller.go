package controller

import (
	"blog-gin-project/database"
	"blog-gin-project/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var input struct {
		Title   string `json:"title" form:"title" binding:"required"`
		Content string `json:"content" form:"content" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建文章"})
		return
	}

	c.Redirect(http.StatusFound, "/posts/"+strconv.Itoa(int(post.ID)))
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	userID, _ := c.Get("userID")
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限编辑此文章"})
		return
	}

	var input struct {
		Title   string `json:"title" form:"title" binding:"required"`
		Content string `json:"content" form:"content" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&post).Updates(models.Post{
		Title:   input.Title,
		Content: input.Content,
	})

	c.Redirect(http.StatusFound, "/posts/"+id)
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	userID, _ := c.Get("userID")
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
		return
	}

	database.DB.Delete(&post)
	c.Redirect(http.StatusFound, "/posts")
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	postID := c.Param("id")
	var post models.Post

	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	var input struct {
		Content string `json:"content" form:"content" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	comment := models.Comment{
		Content: input.Content,
		UserID:  userID.(uint),
		PostID:  post.ID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建评论"})
		return
	}

	c.Redirect(http.StatusFound, "/posts/"+postID)
}
