package app

import (
	"amai/blog/app/auth"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GinApp(db *sqlx.DB) *gin.Engine {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.MaxMultipartMemory = 16 << 20 // 16 MiB
	router.HandleMethodNotAllowed = true

	router.Use(gin.CustomRecovery(ginCustomRecovery))
	router.Use(errorHandler())
	router.Use(injectSqlx(db))

	trustedProxyIpV4 := os.Getenv("TRUSTED_PROXY_IPV4")
	trustedProxyIpV6 := os.Getenv("TRUSTED_PROXY_IPV6")
	err := router.SetTrustedProxies([]string{trustedProxyIpV4, trustedProxyIpV6})

	if err != nil {
		fmt.Printf("SetTrustedProxies error. Message %s", err.Error())
	}

	return router
}

func ginCustomRecovery(c *gin.Context, recovered any) {
	args := []any{
		"Panic: ", recovered,
		"Ip", c.ClientIP(),
		"Method", c.Request.Method,
		"Path", c.Request.URL.Path,
		"Query", c.Request.URL.RawQuery,
	}
	if gin.Mode() == gin.DebugMode {
		args = append(args,
			"UserAgent", c.Request.UserAgent(),
			"Stack", string(debug.Stack()),
		)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}

// Middleware
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.Error(errors.New("panic recovered"))

				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
					return
				}

				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
		}
	}
}

func injectSqlx(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.CheckCookieAuth(c)
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="Authorization Required"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			time.Sleep(500 * time.Millisecond)
			return
		}
		c.Next()
	}
}
