package app

import (
	"amai/blog/app/auth"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/time/rate"
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
	router.Use(rateLimit())

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

func rateLimit() gin.HandlerFunc {
	type client struct {
		limiter *rate.Limiter
	}

	var (
		mut     sync.Mutex
		clients = make(map[string]*client)
	)

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		mut.Lock()
		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{limiter: rate.NewLimiter(5, 5)}
		}
		cl := clients[ip]
		mut.Unlock()

		if !cl.limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded. 5 req/sec only"})
			return
		}
		ctx.Next()
	}

}
