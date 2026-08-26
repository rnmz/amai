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

func parseRawArticle(post data.ArticleEntity) ArticleJSON {
	result := ArticleJSON{
		Id:      post.Id,
		Title:   post.Title,
		Poster:  post.Poster,
		Created: post.Created,
		Body:    post.Body,
		Updated: nil,
	}

	if !post.Created.Equal(post.Updated) {
		result.Updated = &post.Updated
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

	rawPost, dataErr := data.GetArticleById(db, c.Request.Context(), id)
	if dataErr != nil {
		c.Error(dataErr)
		c.Abort()
		return
	}

	post := parseRawArticle(rawPost)
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

	rawPosts, postsErr := data.GetAllArticles(db, c.Request.Context(), page)
	if postsErr != nil {
		c.Error(postsErr)
		c.Abort()
		return
	}

	var posts []ArticleJSON
	for _, e := range rawPosts {
		posts = append(posts, parseRawArticle(e))
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts, "pages": pages})
}

func ArticleCreate(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var postJson ArticleJSON

	if err := c.ShouldBindJSON(&postJson); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON body"})
		return
	}

	err := data.AddArticle(db, c.Request.Context(), data.ArticleEntity{
		Title:  postJson.Title,
		Poster: postJson.Poster,
		Body:   postJson.Body,
	})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "post created"})
}

func ArticleEdit(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var postJson ArticleJSON

	if err := c.ShouldBindJSON(&postJson); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON body"})
		return
	}
	if postJson.Id == uuid.Nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid UUID"})
		return
	}

	err := data.EditArticle(db, c.Request.Context(), data.ArticleEntity{
		Id:     postJson.Id,
		Title:  postJson.Title,
		Poster: postJson.Poster,
		Body:   postJson.Body,
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

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
