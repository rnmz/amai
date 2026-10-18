package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cookieName      = "sessionId"
	sessionMaxAge   = 7 * 24 * time.Hour
	refreshEvery    = 5 * time.Minute
	maxSessionIDLen = 128
)

var sessionTTL = 24 * time.Hour

type sessionStruct struct {
	User     string
	Created  time.Time
	Expires  time.Time
	LastSeen time.Time
}

type sessionKey [sha256.Size]byte

var (
	sessionsMu sync.Mutex
	sessions   = make(map[sessionKey]sessionStruct)
)

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
	id, err := c.Cookie(cookieName)
	if err != nil || id == "" || len(id) > maxSessionIDLen {
		slog.Warn("[Cookie] Cookie not found or invalid", "ip", c.ClientIP())
		return &ErrorAuth{ErrorType: ErrorCookieInvalid, ErrorMessage: "Cookie not found or invalid"}
	}

	key := keyOf(id)
	switch touchSession(key, time.Now()) {
	case sessionMissing:
		slog.Warn("[Cookie] Session doesn't exist", "ip", c.ClientIP())
		return &ErrorAuth{ErrorType: ErrorCookieNotExists, ErrorMessage: "Not found"}
	case sessionExpired:
		slog.Warn("[Cookie] Session expired", "session", fingerprint(key))
		return &ErrorAuth{ErrorType: ErrorCookieExpired, ErrorMessage: "Expired"}
	case sessionRefreshed:
		setSessionCookie(c, id, int(sessionTTL.Seconds()))
	}
	return nil
}

func Login(c *gin.Context) error {
	c.Header("Cache-Control", "no-store")

	user, pass, ok := c.Request.BasicAuth()
	if !ok {
		slog.Warn("[Auth] Basic auth header missing or malformed", "ip", c.ClientIP())
		return &ErrorAuth{ErrorType: ErrorBasicAuthNotValid, ErrorMessage: "Invalid basic auth header"}
	}

	validUser := os.Getenv("ADMIN_LOGIN")
	validPass := os.Getenv("ADMIN_PASSWORD")
	if validUser == "" || validPass == "" {
		slog.Error("[Auth] Admin credentials are not set in environment variables")
		return &ErrorAuth{ErrorType: ErrorAdminCredentialsInvalid, ErrorMessage: "Internal server configuration error"}
	}

	userValid := secureEqual(user, validUser)
	passValid := secureEqual(pass, validPass)
	if !userValid || !passValid {
		slog.Warn("[Auth] Invalid login attempt", "ip", c.ClientIP())
		return &ErrorAuth{ErrorType: ErrorBasicAuthNotValid, ErrorMessage: "Invalid username or password"}
	}

	if old, err := c.Cookie(cookieName); err == nil && old != "" {
		deleteSession(keyOf(old))
	}

	id, err := generateSessionID()
	if err != nil {
		return fmt.Errorf("generate session id: %w", err)
	}

	now := time.Now()
	key := keyOf(id)

	sessionsMu.Lock()
	sessions[key] = sessionStruct{
		User:     user,
		Created:  now,
		Expires:  now.Add(sessionTTL),
		LastSeen: now,
	}
	sessionsMu.Unlock()

	setSessionCookie(c, id, int(sessionTTL.Seconds()))
	slog.Info("[Auth] User logged in", "ip", c.ClientIP(), "session", fingerprint(key))
	return nil
}

func Logout(c *gin.Context) error {
	setSessionCookie(c, "", -1)

	id, err := c.Cookie(cookieName)
	if err != nil {
		slog.Warn("[Cookie] Logout without session cookie", "ip", c.ClientIP())
		return &ErrorAuth{ErrorType: ErrorCookieNotExists, ErrorMessage: "Cookie not found"}
	}
	if id == "" || len(id) > maxSessionIDLen {
		slog.Warn("[Cookie] Logout with invalid session cookie", "ip", c.ClientIP())
		return &ErrorAuth{ErrorType: ErrorCookieInvalid, ErrorMessage: "Cookie invalid"}
	}

	key := keyOf(id)
	if !deleteSession(key) {
		slog.Warn("[Session] Logout for non-existent or expired session", "ip", c.ClientIP())
		return &ErrorAuth{ErrorType: ErrorCookieNotExists, ErrorMessage: "Session not found"}
	}

	slog.Info("[Session] User logged out", "session", fingerprint(key))
	return nil
}

func CleanupSessions() {
	now := time.Now()
	removed := 0

	sessionsMu.Lock()
	for key, s := range sessions {
		if isExpired(s, now) {
			delete(sessions, key)
			removed++
		}
	}
	sessionsMu.Unlock()

	if removed > 0 {
		slog.Info("[Session] Cleanup completed", "removed", removed)
	}
}

type touchResult int

const (
	sessionMissing touchResult = iota
	sessionExpired
	sessionValid
	sessionRefreshed
)

func touchSession(key sessionKey, now time.Time) touchResult {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	s, ok := sessions[key]
	if !ok {
		return sessionMissing
	}
	if isExpired(s, now) {
		delete(sessions, key)
		return sessionExpired
	}
	if now.Sub(s.LastSeen) < refreshEvery {
		return sessionValid
	}

	s.LastSeen = now
	s.Expires = now.Add(sessionTTL)
	sessions[key] = s
	return sessionRefreshed
}

func deleteSession(key sessionKey) bool {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	_, ok := sessions[key]
	delete(sessions, key)
	return ok
}

func isExpired(s sessionStruct, now time.Time) bool {
	return now.After(s.Expires) || now.After(s.Created.Add(sessionMaxAge))
}

func setSessionCookie(c *gin.Context, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   os.Getenv("GIN_MODE") == "release",
		SameSite: http.SameSiteLaxMode,
	})
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func keyOf(id string) sessionKey {
	return sha256.Sum256([]byte(id))
}

func fingerprint(k sessionKey) string {
	return hex.EncodeToString(k[:4])
}

func secureEqual(a, b string) bool {
	ha, hb := sha256.Sum256([]byte(a)), sha256.Sum256([]byte(b))
	return subtle.ConstantTimeCompare(ha[:], hb[:]) == 1
}
