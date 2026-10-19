package config

import (
	"amai/blog/app/handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Routing(db *sqlx.DB, e *gin.Engine) {
	fHandler := handler.FileHandler{DB: db}
	aHandler := handler.ArticleHandler{DB: db}
	user := e.Group("/")
	{
		user.GET("/article/get", aHandler.ArticleGetById)
		user.GET("/article/all", aHandler.ArticleGetAll)
		user.GET("/file/get", fHandler.FileGet)

		user.POST("/auth/login", authRateLimit(), handler.AuthLogin)
		user.POST("/auth/logout", handler.AuthLogout)

		user.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "alive"}) })
	}

	admin := e.Group("/admin/")
	admin.Use(authMiddleware())
	{
		admin.POST("/article/create", aHandler.ArticleCreate)
		admin.PUT("/article/edit", aHandler.ArticleEdit)
		admin.DELETE("/article/delete", aHandler.ArticleDelete)

		admin.GET("/file/list", fHandler.GetFilesList)
		admin.POST("/file/upload", fHandler.FileUpload)
		admin.DELETE("/file/delete", fHandler.FileDelete)
		admin.GET("/whoami", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"authenticated": true}) })
	}
}
