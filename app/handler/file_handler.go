package handler

import (
	"amai/blog/app/data"
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const maxUploadSize = 32 << 20

var allowedMimeTypes = map[string]bool{
	"image/jpeg":                    true,
	"image/png":                     true,
	"image/webp":                    true,
	"image/gif":                     true,
	"text/markdown":                 true,
	"text/plain":                    true,
	"text/csv":                      true,
	"application/pdf":               true,
	"application/zip":               true,
	"application/msword":            true,
	"application/vnd.ms-powerpoint": true,
	"application/vnd.ms-excel":      true,

	"application/vnd.oasis.opendocument.text":                                   true,
	"application/vnd.oasis.opendocument.presentation":                           true,
	"application/vnd.oasis.opendocument.spreadsheet":                            true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
}

var (
	legacyOffice = map[string]bool{
		"application/msword":            true,
		"application/vnd.ms-powerpoint": true,
		"application/vnd.ms-excel":      true,
	}
	oleMagic  = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
	extRegexp = regexp.MustCompile(`^\.[a-z0-9]{1,10}$`)
)

type FileJSON struct {
	FileId  uuid.UUID `json:"id"`
	FileExt string    `json:"ext"`
}

func parseRawFiles(file data.FileEntity) FileJSON {
	return FileJSON{
		FileId:  file.FileId,
		FileExt: file.FileExt,
	}
}

type FileHandler struct {
	DB *sqlx.DB
}

func (f *FileHandler) FileGet(c *gin.Context) {
	id, ok := queryUUID(c)
	if !ok {
		return
	}

	filePath, err := data.GetFileById(f.DB, c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Header("X-Content-Type-Options", "nosniff")
	c.File(filePath)
}

func (f *FileHandler) GetFilesList(c *gin.Context) {
	db := f.DB
	ctx := c.Request.Context()

	pageQuery := c.Query("page")
	if pageQuery == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query page not set"})
		return
	}

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query page should be int"})
		return
	}
	if page <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query page must be greater than 0"})
		return
	}

	pages, err := data.GetAllPagesFile(db, ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if page > pages {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "page not found"})
		return
	}

	rawFiles, err := data.GetAllFiles(db, ctx, page)
	if err != nil {
		_ = c.Error(err)
		return
	}

	files := make([]FileJSON, 0, len(rawFiles))
	for _, e := range rawFiles {
		files = append(files, parseRawFiles(e))
	}

	c.JSON(http.StatusOK, gin.H{"files": files, "pages": pages})
}

func (f *FileHandler) FileUpload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, err := c.FormFile("file")
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"message": "file too large"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no file uploaded"})
		return
	}
	if file.Size == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "empty file uploaded"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != "" && !extRegexp.MatchString(ext) {
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"message": "invalid file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		uploadFailed(c, err)
		return
	}
	defer src.Close()

	head := make([]byte, 512)
	n, err := io.ReadFull(src, head)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		uploadFailed(c, err)
		return
	}
	head = head[:n]

	if !contentAllowed(file.Header.Get("Content-Type"), head) {
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"message": "invalid file uploaded"})
		return
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		uploadFailed(c, err)
		return
	}

	fileId, err := data.UploadFile(f.DB, src, ext)
	if err != nil {
		uploadFailed(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded.", "file_id": fileId})
}

func (f *FileHandler) FileDelete(c *gin.Context) {
	id, ok := queryUUID(c)
	if !ok {
		return
	}

	if err := data.DeleteFile(f.DB, c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}

func queryUUID(c *gin.Context) (uuid.UUID, bool) {
	raw := c.Query("id")
	if raw == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "query id not set"})
		return uuid.Nil, false
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid UUID format"})
		return uuid.Nil, false
	}
	return id, true
}

func uploadFailed(c *gin.Context, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to process file"})
}

func contentAllowed(declared string, head []byte) bool {
	declared = mediaType(declared)
	if !allowedMimeTypes[declared] {
		return false
	}

	detected := mediaType(http.DetectContentType(head))
	if detected == "application/octet-stream" {
		return legacyOffice[declared] && bytes.HasPrefix(head, oleMagic)
	}
	return allowedMimeTypes[detected]
}

func mediaType(s string) string {
	mt, _, err := mime.ParseMediaType(s)
	if err != nil {
		return ""
	}
	return mt
}
