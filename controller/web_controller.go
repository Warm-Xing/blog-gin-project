package controller

import (
	"blog-gin-project/database"
	"blog-gin-project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ShowHomePage 显示首页
func ShowHomePage(c *gin.Context) {
	var posts []models.Post
	// 预加载用户信息并按创建时间倒序查询
	result := database.DB.Preload("User").Order("created_at DESC").Limit(6).Find(&posts)
	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "无法加载文章列表",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Posts": posts,
		"Year":  time.Now().Year(),
	})
}

// ListPosts 显示所有文章
func ListPosts(c *gin.Context) {
	var posts []models.Post
	result := database.DB.Preload("User").Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "无法加载文章列表",
		})
		return
	}

	c.HTML(http.StatusOK, "posts/list.html", gin.H{
		"Posts": posts,
	})
}

// ShowPost 显示单篇文章详情
func ShowPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	// 预加载用户和评论
	result := database.DB.Preload("User").Preload("Comments.User").First(&post, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "文章不存在",
		})
		return
	}

	c.HTML(http.StatusOK, "posts/show.html", gin.H{
		"Post": post,
	})
}

// ShowLoginPage 显示登录页面
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login.html", nil)
}

// ShowRegisterPage 显示注册页面
func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", nil)
}

// ShowDashboard 显示用户仪表盘
func ShowDashboard(c *gin.Context) {
	userID, _ := c.Get("userID")

	var posts []models.Post
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&posts)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Posts": posts,
	})
}

// ShowCreatePostForm 显示创建文章表单
func ShowCreatePostForm(c *gin.Context) {
	c.HTML(http.StatusOK, "posts/create.html", nil)
}

// ShowEditPostForm 显示编辑文章表单
func ShowEditPostForm(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	result := database.DB.First(&post, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "文章不存在",
		})
		return
	}

	// 检查权限
	userID, _ := c.Get("userID")
	if post.UserID != userID {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "没有权限编辑此文章",
		})
		return
	}

	c.HTML(http.StatusOK, "posts/edit.html", gin.H{
		"Post": post,
	})
}
