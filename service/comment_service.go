package service

import (
	"blog-gin-project/database"
	"blog-gin-project/models"
	"blog-gin-project/util"
)

// CommentService 评论服务
type CommentService struct{}

// NewCommentService 创建新的评论服务实例
func NewCommentService() *CommentService {
	return &CommentService{}
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(userID, postID uint, content string) (*models.Comment, error) {
	// 检查文章是否存在
	var post models.Post
	result := database.DB.First(&post, postID)
	if result.Error != nil {
		return nil, util.ErrPostNotFound
	}

	comment := models.Comment{
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}

	result = database.DB.Create(&comment)
	if result.Error != nil {
		return nil, result.Error
	}

	// 预加载用户信息
	database.DB.Preload("User").First(&comment, comment.ID)

	return &comment, nil
}

// GetCommentsByPostID 根据文章ID获取评论
func (s *CommentService) GetCommentsByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	result := database.DB.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(commentID, userID uint) error {
	var comment models.Comment
	result := database.DB.First(&comment, commentID)
	if result.Error != nil {
		return util.ErrCommentNotFound
	}

	// 检查权限（只有评论作者可以删除自己的评论）
	if comment.UserID != userID {
		return util.ErrUnauthorized
	}

	result = database.DB.Delete(&comment)
	return result.Error
}
