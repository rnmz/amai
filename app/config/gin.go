package config

import (
	"amai/blog/app/auth"
	"log/slog"
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
	} else if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		panic("GIN MODE ERROR: GIN_MODE MUST BE release OR debug")
	}

	router := gin.New()

	router.MaxMultipartMemory = 16 << 20 // 16 MiB
	router.HandleMethodNotAllowed = true

	router.Use(gin.CustomRecovery(ginCustomRecovery))
	router.Use(errorHandler())
	router.Use(injectSqlx(db))
	router.Use(rateLimit())
	router.Use(cors())

	trustedProxies := make([]string, 0, 2)
	if v := os.Getenv("TRUSTED_PROXY_IPV4"); v != "" {
		trustedProxies = append(trustedProxies, v)
	}
	if v := os.Getenv("TRUSTED_PROXY_IPV6"); v != "" {
		trustedProxies = append(trustedProxies, v)
	}

	if len(trustedProxies) > 0 {
		if err := router.SetTrustedProxies(trustedProxies); err != nil {
			slog.Error("SetTrustedProxies error", "error", err)
		}
	} else {
		if err := router.SetTrustedProxies(nil); err != nil {
			slog.Error("SetTrustedProxies error", "error", err)
		}
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
	slog.Error("[Gin] Panic recovered", "args", args)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}

// Middleware
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			if !c.Writer.Written() {
				lastErr := c.Errors.Last()
				slog.Error("[Gin] Error handled", "error", lastErr)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			slog.Info("[Gin] Auth required")
			return
		}
		c.Next()
	}
}

func rateLimit() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mut     sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mut.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mut.Unlock()
		}
	}()

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		mut.Lock()
		cl, exists := clients[ip]
		if !exists {
			cl = &client{limiter: rate.NewLimiter(30, 30)}
			clients[ip] = cl
		}
		cl.lastSeen = time.Now()
		mut.Unlock()

		if !cl.limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded. 5 req/sec only"})
			return
		}
		ctx.Next()
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("TRUSTED_DOMAIN"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
