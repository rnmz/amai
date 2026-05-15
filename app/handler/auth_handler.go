package handler

import (
	"amai/blog/app/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthLogin(c *gin.Context) {
	err := auth.Login(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "login!"})
}

func AuthLogout(c *gin.Context) {
	err := auth.Logout(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "logout!"})
}
