package config

import (
	"amai/blog/app/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routing(e *gin.Engine) {
	user := e.Group("/")
	{
		user.GET("/post/get", handler.PostGetById)
		user.GET("/post/all", handler.PostGetAll)
		user.GET("/file/get", handler.FileGet)

		user.POST("/auth/login", handler.AuthLogin)
		user.POST("/auth/logout", handler.AuthLogout)

		user.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "alive"}) })
	}

	admin := e.Group("/admin/")
	admin.Use(authMiddleware())
	{
		admin.POST("/post/create", handler.PostCreate)
		admin.PUT("/post/edit", handler.PostEdit)
		admin.DELETE("/post/delete", handler.PostDelete)

		admin.GET("/file/list", handler.GetFilesList)
		admin.POST("/file/upload", handler.FileUpload)
		admin.DELETE("/file/delete", handler.FileDelete)
		admin.GET("/whoami", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"authenticated": true}) })
	}
}
