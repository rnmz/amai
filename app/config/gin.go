package app

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type sessionStruct struct {
	User     string
	Created  time.Time
	Expires  time.Time
	LastSeen time.Time
}

var sessions sync.Map
var sessionTTL = 24 * time.Hour

func GinApp(db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
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
		if handleCookieAuth(c) || handleAuth(c) {
			c.Next()
			return
		}
		c.Header("WWW-Authenticate", `Basic realm="Authorization Required"`)
		c.AbortWithStatus(http.StatusUnauthorized)
		time.Sleep(500 * time.Millisecond)
	}
}

func handleCookieAuth(c *gin.Context) bool {
	sessionId, err := c.Cookie("sessionId")

	if err != nil || sessionId == "" {
		return false
	}

	val, exists := sessions.Load(sessionId)
	if !exists {
		return false
	}

	session := val.(sessionStruct)

	if time.Now().After(session.Expires) {
		sessions.Delete(sessionId)
		return false
	}

	session.LastSeen = time.Now()
	session.Expires = time.Now().Add(24 * time.Hour)
	sessions.Store(sessionId, session)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return true
}

func handleAuth(c *gin.Context) bool {
	user, pass, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		return false
	}

	validUser := os.Getenv("admin_login")
	validPass := os.Getenv("admin_password")
	if validUser == "" || validPass == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return false
	}

	userValid := subtle.ConstantTimeCompare([]byte(user), []byte(validUser)) == 1
	passValid := subtle.ConstantTimeCompare([]byte(pass), []byte(validPass)) == 1
	if !userValid || !passValid {
		return false
	}

	sessionId := generateSessionID()
	sessions.Store(sessionId, sessionStruct{
		User:     user,
		Created:  time.Now(),
		Expires:  time.Now().Add(sessionTTL),
		LastSeen: time.Now(),
	})
	c.SetCookie("sessionId", sessionId, int(sessionTTL.Seconds()), "/", "", true, true)
	return true
}

func CleanupSessions() {
	sessions.Range(func(key, value any) bool {
		session := value.(sessionStruct)
		if time.Now().After(session.Expires) {
			sessions.Delete(key)
		}
		return true
	})
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
