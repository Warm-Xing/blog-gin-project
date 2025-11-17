package service

import (
	"blog-gin-project/database"
	"blog-gin-project/models"
	"blog-gin-project/util"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建新的用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// Register 用户注册
func (s *UserService) Register(username, email, password string) error {
	// 检查用户名是否已存在
	var existingUser models.User
	result := database.DB.Where("username = ?", username).First(&existingUser)
	if result.Error != gorm.ErrRecordNotFound {
		return util.ErrUsernameExists
	}

	// 检查邮箱是否已存在
	result = database.DB.Where("email = ?", email).First(&existingUser)
	if result.Error != gorm.ErrRecordNotFound {
		return util.ErrEmailExists
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建新用户
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	result = database.DB.Create(&user)
	return result.Error
}

// Login 用户登录
func (s *UserService) Login(username, password string) (string, error) {
	// 查找用户
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return "", util.ErrInvalidCredentials
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", util.ErrInvalidCredentials
	}

	// 生成JWT令牌
	token, err := util.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
