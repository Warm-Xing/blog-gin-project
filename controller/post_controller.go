package controller

import (
	"blog-gin-project/service"
	"blog-gin-project/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostController 文章控制器
type PostController struct {
	postService *service.PostService
}

// NewPostController 创建新的文章控制器实例
func NewPostController() *PostController {
	return &PostController{
		postService: service.NewPostService(),
	}
}

// CreatePostRequest 创建文章请求结构体
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200"`
	Content string `json:"content" binding:"required,min=10"`
}

// UpdatePostRequest 更新文章请求结构体
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200"`
	Content string `json:"content" binding:"required,min=10"`
}

// CreatePost 创建文章
func (c *PostController) CreatePost(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析请求
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 创建文章
	post, err := c.postService.CreatePost(userID.(uint), req.Title, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    post,
	})
}

// GetPost 获取文章详情
func (c *PostController) GetPost(ctx *gin.Context) {
	// 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取文章
	post, err := c.postService.GetPostByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"post": post})
}

// GetAllPosts 获取所有文章
func (c *PostController) GetAllPosts(ctx *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取文章列表
	posts, total, err := c.postService.GetAllPosts(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败: " + err.Error()})
		return
	}

	// 计算总页数
	totalPages := (int(total) + pageSize - 1) / pageSize

	ctx.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

// UpdatePost 更新文章
func (c *PostController) UpdatePost(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 解析请求
	var req UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 更新文章
	post, err := c.postService.UpdatePost(uint(id), userID.(uint), req.Title, req.Content)
	if err != nil {
		switch err {
		case util.ErrPostNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		case util.ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限更新此文章"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败: " + err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post":    post,
	})
}

// DeletePost 删除文章
func (c *PostController) DeletePost(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 删除文章
	err = c.postService.DeletePost(uint(id), userID.(uint))
	if err != nil {
		switch err {
		case util.ErrPostNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		case util.ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败: " + err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
