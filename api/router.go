package api

import (
	"blog-gin-project/controller"
	"blog-gin-project/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
//r := gin.New() 要使用这个日志中间件，需要在api/router.go文件的SetupRouter函数中，将gin.Default()替换为：  查看ai
//r.Use(gin.Recovery())
//r.Use(middleware.Logger())

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 公开路由
	public := r.Group("/api")
	{
		// 用户相关路由
		user := public.Group("/users")
		{
			user.POST("/register", controller.NewUserController().Register)
			user.POST("/login", controller.NewUserController().Login)
		}

		// 文章相关路由
		posts := public.Group("/posts")
		{
			posts.GET("", controller.NewPostController().GetAllPosts)
			posts.GET("/:id", controller.NewPostController().GetPost)

			// 评论相关路由
			posts.GET("/:post_id/comments", controller.NewCommentController().GetCommentsByPostID)
		}
	}

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// 用户相关路由
		user := authorized.Group("/users")
		{
			user.GET("/profile", controller.NewUserController().GetProfile)
		}

		// 文章相关路由
		posts := authorized.Group("/posts")
		{
			posts.POST("", controller.NewPostController().CreatePost)
			posts.PUT("/:id", controller.NewPostController().UpdatePost)
			posts.DELETE("/:id", controller.NewPostController().DeletePost)

			// 评论相关路由
			posts.POST("/:post_id/comments", controller.NewCommentController().CreateComment)
			posts.DELETE("/comments/:id", controller.NewCommentController().DeleteComment)
		}
	}

	return r
}
