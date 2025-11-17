package util

import "errors"

// 自定义错误
var (
	ErrUsernameExists     = errors.New("用户名已存在")
	ErrEmailExists        = errors.New("邮箱已存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrPostNotFound       = errors.New("文章不存在")
	ErrCommentNotFound    = errors.New("评论不存在")
	ErrUnauthorized       = errors.New("未授权操作")
)
