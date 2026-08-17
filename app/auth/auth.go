package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type sessionStruct struct {
	User     string
	Created  time.Time
	Expires  time.Time
	LastSeen time.Time
}

var sessions sync.Map
var sessionTTL = 24 * time.Hour

type ErrorAuth struct {
	ErrorType    string
	ErrorMessage string
}

const (
	ErrorBasicAuthNotValid       = "NOT VALID"
	ErrorAdminCredentialsInvalid = "ADMIN CREDENTIALS IS INVALID"

	ErrorCookieInvalid   = "COOKIE INVALID"
	ErrorCookieNotExists = "COOKIE NOT EXISTS"
	ErrorCookieExpired   = "COOKIE EXPIRED"
)

func (e *ErrorAuth) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType, e.ErrorMessage)
}

func CheckCookieAuth(c *gin.Context) error {
	sessionId, err := c.Cookie("sessionId")

	if err != nil {
		slog.Warn("[Cookie] Cookie not found or invalid")
		return &ErrorAuth{ErrorType: ErrorCookieInvalid, ErrorMessage: "Cookie not found or invalid"}
	}

	if sessionId == "" {
		slog.Warn("[Cookie] Empty session id")
		return &ErrorAuth{ErrorType: ErrorCookieInvalid, ErrorMessage: "Empty sessionId"}
	}

	val, exists := sessions.Load(sessionId)
	if !exists {
		slog.Warn("[Cookie] Session doesn't exists")
		return &ErrorAuth{ErrorType: ErrorCookieNotExists, ErrorMessage: "No found"}
	}

	session := val.(sessionStruct)

	if time.Now().After(session.Expires) {
		sessions.Delete(sessionId)
		slog.Warn("[Cookie] Session expired", "sessionId", sessionId)
		return &ErrorAuth{ErrorType: ErrorCookieExpired, ErrorMessage: "Expired"}
	}

	currTime := time.Now()
	session.LastSeen = currTime
	session.Expires = currTime.Add(sessionTTL)
	sessions.Store(sessionId, session)

	var isSecure bool
	if os.Getenv("GIN_MODE") == "release" {
		isSecure = true
	} else {
		isSecure = false
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
	slog.Info("[Cookie] Session cookie updated", "sessionId", sessionId)
	return nil
}

func Login(c *gin.Context) error {
	user, pass, ok := c.Request.BasicAuth()
	if !ok {
		slog.Warn("[Auth] Basic auth header missing or malformed")
		return &ErrorAuth{ErrorType: ErrorBasicAuthNotValid, ErrorMessage: "Invalid basic auth header"}
	}

	validUser := os.Getenv("ADMIN_LOGIN")
	validPass := os.Getenv("ADMIN_PASSWORD")
	if validUser == "" || validPass == "" {
		slog.Error("[Auth] Admin credentials are not set in environment variables")
		return &ErrorAuth{ErrorType: ErrorAdminCredentialsInvalid, ErrorMessage: "Internal server configuration error"}
	}

	userValid := subtle.ConstantTimeCompare([]byte(user), []byte(validUser)) == 1
	passValid := subtle.ConstantTimeCompare([]byte(pass), []byte(validPass)) == 1
	if !userValid || !passValid {
		slog.Warn("[Auth] Invalid login credentials attempt", "username", user)
		return &ErrorAuth{ErrorType: ErrorBasicAuthNotValid, ErrorMessage: "Invalid username or password"}
	}

	sessionId := generateSessionID()
	currTime := time.Now()
	sessions.Store(sessionId, sessionStruct{
		User:     user,
		Created:  currTime,
		Expires:  currTime.Add(sessionTTL),
		LastSeen: currTime,
	})

	var isSecure bool
	if os.Getenv("GIN_MODE") == "release" {
		isSecure = true
	} else {
		isSecure = false
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})

	slog.Info("[Auth] User logged in successfully", "username", user, "sessionId", sessionId)
	return nil
}

func Logout(c *gin.Context) error {
	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		slog.Warn("[Cookie] Logout attempted without sessionId cookie")
		return &ErrorAuth{ErrorType: ErrorCookieNotExists}
	}

	if sessionId == "" {
		slog.Warn("[Cookie] Logout attempted with empty sessionId cookie")
		return &ErrorAuth{ErrorType: ErrorCookieInvalid}
	}

	_, exists := sessions.Load(sessionId)
	if !exists {
		slog.Warn("[Session] Logout attempted for non-existent or expired session", "sessionId", sessionId)
		return &ErrorAuth{ErrorType: ErrorCookieNotExists}
	}

	sessions.Delete(sessionId)

	var isSecure bool
	if os.Getenv("GIN_MODE") == "release" {
		isSecure = true
	} else {
		isSecure = false
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
	slog.Info("[Session] User logged out successfully", "sessionId", sessionId)
	return nil
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func CleanupSessions() {
	slog.Debug("[Session] Running background session cleanup")
	deletedCount := 0

	sessions.Range(func(key, value any) bool {
		session := value.(sessionStruct)
		if time.Now().After(session.Expires) {
			sessions.Delete(key)
			deletedCount++
			slog.Debug("[Session] Expired session removed during cleanup", "sessionId", key)
		}
		return true
	})

	if deletedCount > 0 {
		slog.Info("[Session] Background cleanup completed", "removedSessionsCount", deletedCount)
	}
}
