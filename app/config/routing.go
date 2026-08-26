package config

import (
	"amai/blog/app/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routing(e *gin.Engine) {
	user := e.Group("/")
	{
		user.GET("/post/get", handler.ArticleGetById)
		user.GET("/post/all", handler.ArticleGetAll)
		user.GET("/file/get", handler.FileGet)

		user.POST("/auth/login", authRateLimit(), handler.AuthLogin)
		user.POST("/auth/logout", handler.AuthLogout)

		user.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "alive"}) })
	}

	admin := e.Group("/admin/")
	admin.Use(authMiddleware())
	{
		admin.POST("/post/create", handler.ArticleCreate)
		admin.PUT("/post/edit", handler.ArticleEdit)
		admin.DELETE("/post/delete", handler.ArticleDelete)

		admin.GET("/file/list", handler.GetFilesList)
		admin.POST("/file/upload", handler.FileUpload)
		admin.DELETE("/file/delete", handler.FileDelete)
		admin.GET("/whoami", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"authenticated": true}) })
	}
}
