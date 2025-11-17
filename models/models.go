package models

import (
	"time"

	"gorm.io/gorm"
)

// User 模型表示系统用户
type User struct {
	gorm.Model
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:100;not null" json:"-"` // json:"-" 表示序列化时忽略此字段
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Posts    []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// Post 模型表示博客文章
type Post struct {
	gorm.Model
	Title    string    `gorm:"size:200;not null" json:"title"`
	Content  string    `gorm:"type:text;not null" json:"content"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// Comment 模型表示文章评论
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	PostID  uint   `gorm:"not null" json:"post_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post    Post   `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// BeforeCreate 在创建记录前设置创建时间
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 在更新记录前设置更新时间
func (p *Post) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
