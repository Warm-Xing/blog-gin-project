package routes

import (
	"blog-gin-project/controller"
	"blog-gin-project/middleware"
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// LoadTemplates 加载HTML模板并注册自定义函数
func LoadTemplates(r *gin.Engine) {
	// 定义自定义模板函数
	funcMap := template.FuncMap{
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
	}

	// 解析所有模板文件并关联函数映射
	templates, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		panic("无法解析模板文件: " + err.Error())
	}

	// 配置Gin模板引擎
	r.SetHTMLTemplate(templates)
}

// SetupWebRoutes 设置Web路由
func SetupWebRoutes(r *gin.Engine) {
	// 配置静态文件
	r.Static("/static", "./static")

	// 首页路由
	r.GET("/", controller.ShowHomePage)

	// 文章相关路由
	r.GET("/posts", controller.ListPosts)
	r.GET("/posts/:id", controller.ShowPost)

	// 用户认证路由
	r.GET("/login", controller.ShowLoginPage)
	r.GET("/register", controller.ShowRegisterPage)

	// 需要认证的路由组
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/dashboard", controller.ShowDashboard)
		auth.GET("/posts/create", controller.ShowCreatePostForm)
		//auth.POST("/posts", postController.CreatePost)
		//auth.GET("/posts/:id/edit", controller.ShowEditPostForm)
		//auth.POST("/posts/:id/update", controller.UpdatePost)
		//auth.POST("/posts/:id/delete", controller.DeletePost)
		//auth.POST("/posts/:id/comments", controller.CreateComment)
	}
}
