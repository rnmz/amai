package config

import (
	"amai/blog/app/auth"
	"errors"
	"fmt"
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

func GinApp(db *sqlx.DB) (*gin.Engine, error) {
	switch mode := os.Getenv("GIN_MODE"); mode {
	case gin.ReleaseMode, gin.DebugMode:
		gin.SetMode(mode)
	default:
		return nil, errors.New("GIN_MODE must be release or debug")
	}

	router := gin.New()
	router.MaxMultipartMemory = 16 << 20
	router.HandleMethodNotAllowed = true

	if err := router.SetTrustedProxies(trustedProxies()); err != nil {
		return nil, fmt.Errorf("set trusted proxies: %w", err)
	}

	router.Use(
		gin.CustomRecovery(recoveryHandler),
		cors(os.Getenv("TRUSTED_DOMAIN")),
		rateLimit(30, 30),
		errorHandler(),
		injectSqlx(db),
	)

	return router, nil
}

func trustedProxies() []string {
	var list []string
	for _, key := range []string{"TRUSTED_PROXY_IPV4", "TRUSTED_PROXY_IPV6"} {
		if v := os.Getenv(key); v != "" {
			list = append(list, v)
		}
	}
	return list
}

func recoveryHandler(c *gin.Context, recovered any) {
	slog.Error("[Gin] panic recovered",
		"panic", recovered,
		"ip", c.ClientIP(),
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"query", c.Request.URL.RawQuery,
		"user_agent", c.Request.UserAgent(),
		"stack", string(debug.Stack()),
	)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}

// Middleware

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		slog.Error("[Gin] error handled",
			"error", c.Errors.Last(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)

		if !c.Writer.Written() {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		}
	}
}

// todo: inject into StartServer, not into GinApp

func injectSqlx(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.CheckCookieAuth(c); err != nil {
			slog.Info("[Gin] auth required", "ip", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		c.Next()
	}
}

func authRateLimit() gin.HandlerFunc {
	return rateLimit(1, 1)
}

func rateLimit(rps, burst int) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	message := fmt.Sprintf("rate limit exceeded. %d req/sec only", rps)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		cl, ok := clients[ip]
		if !ok {
			cl = &client{limiter: rate.NewLimiter(rate.Limit(rps), burst)}
			clients[ip] = cl
		}
		cl.lastSeen = time.Now()
		mu.Unlock()

		if !cl.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": message})
			return
		}
		c.Next()
	}
}

func cors(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", allowOrigin)
		h.Set("Access-Control-Allow-Credentials", "true")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
