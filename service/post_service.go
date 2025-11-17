package service

import (
	"blog-gin-project/database"
	"blog-gin-project/models"
	"blog-gin-project/util"
)

// PostService 文章服务
type PostService struct{}

// NewPostService 创建新的文章服务实例
func NewPostService() *PostService {
	return &PostService{}
}

// CreatePost 创建文章
func (s *PostService) CreatePost(userID uint, title, content string) (*models.Post, error) {
	post := models.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	result := database.DB.Create(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return &post, nil
}

// GetPostByID 根据ID获取文章
func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	result := database.DB.Preload("User").Preload("Comments").Preload("Comments.User").First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// GetAllPosts 获取所有文章
func (s *PostService) GetAllPosts(page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// 获取总数
	database.DB.Model(&models.Post{}).Count(&total)

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	result := database.DB.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return posts, total, nil
}

// UpdatePost 更新文章
func (s *PostService) UpdatePost(postID, userID uint, title, content string) (*models.Post, error) {
	// 获取文章
	post, err := s.GetPostByID(postID)
	if err != nil {
		return nil, util.ErrPostNotFound
	}

	// 检查权限
	if post.UserID != userID {
		return nil, util.ErrUnauthorized
	}

	// 更新文章
	post.Title = title
	post.Content = content

	result := database.DB.Save(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return post, nil
}

// DeletePost 删除文章
func (s *PostService) DeletePost(postID, userID uint) error {
	// 获取文章
	post, err := s.GetPostByID(postID)
	if err != nil {
		return util.ErrPostNotFound
	}

	// 检查权限
	if post.UserID != userID {
		return util.ErrUnauthorized
	}

	// 删除文章
	result := database.DB.Delete(&models.Post{}, postID)
	return result.Error
}
