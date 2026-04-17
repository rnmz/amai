package handler

import (
	"amai/blog/app/data"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func FileGet(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)

	queryId := c.Query("id")
	if queryId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query id not set"})
		return
	}

	id, parseErr := uuid.Parse(queryId)
	if parseErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid UUID format"})
		return
	}

	filePath, fileErr := data.GetFileById(db, c.Request.Context(), id)
	if fileErr != nil {
		c.Error(fileErr)
		c.Abort()
		return
	}

	c.File(filePath)
}

func FileUpload(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)

	allowedFileMimeType := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/webp":      true,
		"image/gif":       true,
		"text/markdown":   true,
		"text/plain":      true,
		"application/pdf": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	}

	file, formFileErr := c.FormFile("file")

	if formFileErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no file uploaded"})
		return
	}
	if file.Size == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "empty file uploaded"})
		return
	}
	if !allowedFileMimeType[file.Header.Get("Content-Type")] {
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"message": "invalid file uploaded"})
		return
	}

	src, fileErr := file.Open()
	if fileErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to process file"})
		return
	}

	defer src.Close()

	fileId, uploadErr := data.UploadFile(db, src, filepath.Ext(file.Filename))
	if uploadErr != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "failed to process file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded. File id: " + fileId})
}

func FileDelete(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)

	queryId := c.Query("id")
	if queryId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query id not set"})
		return
	}

	id, parseErr := uuid.Parse(queryId)
	if parseErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid UUID format"})
		return
	}

	err := data.DeleteFile(db, c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}
