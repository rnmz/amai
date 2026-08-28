package handler

import (
	"amai/blog/app/data"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticleJSON struct {
	Id      uuid.UUID  `json:"id"`
	Title   string     `json:"title" binding:"required"`
	Poster  uuid.UUID  `json:"poster_id" binding:"required"`
	Created time.Time  `json:"created_at"`
	Updated *time.Time `json:"updated_at"`
	Body    string     `json:"body" binding:"required"`
}

func parseRawArticle(articleEntity data.ArticleEntity) ArticleJSON {
	result := ArticleJSON{
		Id:      articleEntity.Id,
		Title:   articleEntity.Title,
		Poster:  articleEntity.Poster,
		Created: articleEntity.Created,
		Body:    articleEntity.Body,
		Updated: nil,
	}

	if !articleEntity.Created.Equal(articleEntity.Updated) {
		result.Updated = &articleEntity.Updated
	}

	return result
}

func ArticleGetById(c *gin.Context) {
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

	rawArticle, dataErr := data.GetArticleById(db, c.Request.Context(), id)
	if dataErr != nil {
		c.Error(dataErr)
		c.Abort()
		return
	}

	post := parseRawArticle(rawArticle)
	c.JSON(http.StatusOK, post)
}

func ArticleGetAll(c *gin.Context) {
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

	pages, pagesErr := data.GetAllPagesArticles(db, c.Request.Context())
	if pagesErr != nil {
		c.Error(pagesErr)
		c.Abort()
		return
	}

	if page > pages {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "page not found"})
		return
	}

	articles, err := data.GetAllArticles(db, c.Request.Context(), page)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var articleJSONS []ArticleJSON
	for _, e := range articles {
		articleJSONS = append(articleJSONS, parseRawArticle(e))
	}

	c.JSON(http.StatusOK, gin.H{"articles": articleJSONS, "pages": pages})
}

func ArticleCreate(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var articleJSON ArticleJSON

	if err := c.ShouldBindJSON(&articleJSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON body"})
		return
	}

	err := data.AddArticle(db, c.Request.Context(), data.ArticleEntity{
		Title:  articleJSON.Title,
		Poster: articleJSON.Poster,
		Body:   articleJSON.Body,
	})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func ArticleEdit(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var articleJSON ArticleJSON

	if err := c.ShouldBindJSON(&articleJSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON body"})
		return
	}
	if articleJSON.Id == uuid.Nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid UUID"})
		return
	}

	err := data.EditArticle(db, c.Request.Context(), data.ArticleEntity{
		Id:     articleJSON.Id,
		Title:  articleJSON.Title,
		Poster: articleJSON.Poster,
		Body:   articleJSON.Body,
	})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post edited"})
}

func ArticleDelete(c *gin.Context) {
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

	err := data.DeleteArticle(db, c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
