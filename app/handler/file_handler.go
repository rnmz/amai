package handler

import (
	"amai/blog/app/data"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FileJSON struct {
	FileId  uuid.UUID `json:"id"`
	FileExt string    `json:"ext"`
}

func parseRawFiles(file data.FileEntity) FileJSON {
	result := FileJSON{
		FileId:  file.FileId,
		FileExt: file.FileExt,
	}

	return result
}

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

func GetFilesList(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)

	pageQuery := c.Query("page")
	if pageQuery == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query page not set"})
		return
	}

	page, convErr := strconv.Atoi(pageQuery)
	if convErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query page should be int"})
		return
	}
	if page <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query page must be greater than 0"})
		return
	}

	pages, pagesErr := data.GetAllPagesFile(db, c.Request.Context())
	if pagesErr != nil {
		c.Error(pagesErr)
		c.Abort()
		return
	}

	if page > pages {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "page not found"})
		return
	}

	rawFiles, filesErr := data.GetAllFiles(db, c.Request.Context(), page)
	if filesErr != nil {
		c.Error(filesErr)
		c.Abort()
		return
	}

	var files []FileJSON
	for _, e := range rawFiles {
		files = append(files, parseRawFiles(e))
	}

	c.JSON(http.StatusOK, gin.H{"files": files, "pages": pages})
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

	buf := make([]byte, 512)
	n, readErr := src.Read(buf)
	if readErr != nil && readErr != io.EOF {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to process file"})
		return
	}
	if _, seekErr := src.Seek(0, io.SeekStart); seekErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to process file"})
		return
	}

	detectedType := http.DetectContentType(buf[:n])
	if !allowedFileMimeType[detectedType] {
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"message": "invalid file uploaded"})
		return
	}

	fileId, uploadErr := data.UploadFile(db, src, filepath.Ext(file.Filename))
	if uploadErr != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "failed to process file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded.", "file_id": fileId})
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
