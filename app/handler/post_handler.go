package handler

import (
	"amai/blog/app/data"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostJSON struct {
	Id      uuid.UUID  `json:"id"`
	Title   string     `json:"title"`
	Poster  string     `json:"poster_id"`
	Created time.Time  `json:"created_at"`
	Updated *time.Time `json:"updated_at"`
	Body    string     `json:"body"`
}

func parseRawPost(post data.Post) PostJSON {
	result := PostJSON{
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

func PostGetById(c *gin.Context) {
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

	rawPost, dataErr := data.GetPostById(db, c.Request.Context(), id)
	if dataErr != nil {
		c.Error(dataErr)
		c.Abort()
		return
	}

	post := parseRawPost(rawPost)
	c.JSON(http.StatusOK, post)
}

func PostGetAll(c *gin.Context) {
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

	pages, pagesErr := data.GetAllPages(db, c.Request.Context())
	if pagesErr != nil {
		fmt.Printf("something went wrong %s", pagesErr.Error())
		c.Error(pagesErr)
		c.Abort()
		return
	}

	if page > pages {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "page not found"})
		return
	}

	rawPosts, postsErr := data.GetAllPosts(db, c.Request.Context(), page)
	if postsErr != nil {
		c.Error(postsErr)
		c.Abort()
		return
	}

	var posts []PostJSON
	for _, e := range rawPosts {
		posts = append(posts, parseRawPost(e))
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts, "pages": pages})
}

func PostCreate(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var postJson PostJSON

	if err := c.ShouldBindJSON(&postJson); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON body"})
		return
	}

	err := data.AddPost(db, c.Request.Context(), data.Post{
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

func PostEdit(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var postJson PostJSON

	if err := c.ShouldBindJSON(&postJson); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON body"})
		return
	}
	if postJson.Id == uuid.Nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid UUID"})
		return
	}

	err := data.EditPost(db, c.Request.Context(), data.Post{
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

func PostDelete(c *gin.Context) {
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

	err := data.DeletePost(db, c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
