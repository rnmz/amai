package app

import (
	"amai/blog/app/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routing(e *gin.Engine) {
	user := e.Group("/")
	{
		user.GET("/post/get", func(ctx *gin.Context) { handler.PostGetById(ctx) })
		user.GET("/post/all", func(ctx *gin.Context) { handler.PostGetAll(ctx) })
		user.GET("/file/get", func(ctx *gin.Context) { handler.FileGet(ctx) })

		user.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "alive"}) })
	}

	admin := e.Group("/admin/")
	admin.Use(authMiddleware())
	{
		admin.POST("/post/create", func(ctx *gin.Context) { handler.PostCreate(ctx) })
		admin.PUT("/post/edit", func(ctx *gin.Context) { handler.PostEdit(ctx) })
		admin.DELETE("/post/delete", func(ctx *gin.Context) { handler.PostDelete(ctx) })

		admin.POST("/file/upload", func(ctx *gin.Context) { handler.FileUpload(ctx) })
		admin.DELETE("/file/delete", func(ctx *gin.Context) { handler.FileDelete(ctx) })
	}
}
