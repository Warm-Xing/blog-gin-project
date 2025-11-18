package controller

import (
	"blog-gin-project/service"
	"blog-gin-project/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *service.CommentService
}

// NewCommentController 创建新的评论控制器实例
func NewCommentController() *CommentController {
	return &CommentController{
		commentService: service.NewCommentService(),
	}
}

// CreateCommentRequest 创建评论请求结构体
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CreateComment 创建评论
func (c *CommentController) CreateComment(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文章ID
	//postIDStr := ctx.Param("post_id")
	postIDStr := ctx.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 解析请求
	var req CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 创建评论
	comment, err := c.commentService.CreateComment(userID.(uint), uint(postID), req.Content)
	if err != nil {
		switch err {
		case util.ErrPostNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败: " + err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// GetCommentsByPostID 根据文章ID获取评论
func (c *CommentController) GetCommentsByPostID(ctx *gin.Context) {
	// 获取文章ID
	postIDStr := ctx.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取评论
	comments, err := c.commentService.GetCommentsByPostID(uint(postID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"comments": comments})
}

// DeleteComment 删除评论
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	// 删除评论
	err = c.commentService.DeleteComment(uint(commentID), userID.(uint))
	if err != nil {
		switch err {
		case util.ErrCommentNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		case util.ErrUnauthorized:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此评论"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败: " + err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}
