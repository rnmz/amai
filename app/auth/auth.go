package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
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

type AuthError struct {
	ErrorType    string
	ErrorMessage string
}

const (
	AuthErrorBasicAuthNotValid       = "NOT VALID"
	AuthErrorAdminCredentialsInvalid = "ADMIN CREDENTIALS IS INVALID"

	AuthErrorCookieInvalid   = "COOKIE INVALID"
	AuthErrorCookieNotExists = "COOKIE NOT EXISTS"
	AuthErrorCookieExpired   = "COOKIE EXPIRED"
)

func (e *AuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType, e.ErrorMessage)
}

func CheckCookieAuth(c *gin.Context) error {
	sessionId, err := c.Cookie("sessionId")

	if err != nil {
		return &AuthError{ErrorType: AuthErrorCookieInvalid, ErrorMessage: "Cookie not found or invalid"}
	}

	if sessionId == "" {
		return &AuthError{ErrorType: AuthErrorCookieInvalid, ErrorMessage: "Empty sessionId"}
	}

	val, exists := sessions.Load(sessionId)
	if !exists {
		return &AuthError{ErrorType: AuthErrorCookieNotExists, ErrorMessage: "No found"}
	}

	session := val.(sessionStruct)

	if time.Now().After(session.Expires) {
		sessions.Delete(sessionId)
		return &AuthError{ErrorType: AuthErrorCookieExpired, ErrorMessage: "Message"}
	}

	currTime := time.Now()
	session.LastSeen = currTime
	session.Expires = currTime.Add(sessionTTL)
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
	return nil
}

func Login(c *gin.Context) error {
	user, pass, ok := c.Request.BasicAuth()
	if !ok {
		return &AuthError{ErrorType: AuthErrorBasicAuthNotValid, ErrorMessage: "Something wrong with basic auth"}
	}

	validUser := os.Getenv("admin_login")
	validPass := os.Getenv("admin_password")
	if validUser == "" || validPass == "" {
		return &AuthError{ErrorType: AuthErrorAdminCredentialsInvalid, ErrorMessage: "Admin login or password not set"}
	}

	userValid := subtle.ConstantTimeCompare([]byte(user), []byte(validUser)) == 1
	passValid := subtle.ConstantTimeCompare([]byte(pass), []byte(validPass)) == 1
	if !userValid || !passValid {
		return &AuthError{ErrorType: AuthErrorBasicAuthNotValid, ErrorMessage: ""}
	}

	sessionId := generateSessionID()
	currTime := time.Now()
	sessions.Store(sessionId, sessionStruct{
		User:     user,
		Created:  currTime,
		Expires:  currTime.Add(sessionTTL),
		LastSeen: currTime,
	})
	c.SetCookie("sessionId", sessionId, int(sessionTTL.Seconds()), "/", "", true, true)
	return nil
}

func Logout(c *gin.Context) error {
	sessionId, err := c.Cookie("sessionId")

	if err != nil {
		return &AuthError{ErrorType: AuthErrorCookieNotExists}
	}

	if sessionId == "" {
		return &AuthError{ErrorType: AuthErrorCookieInvalid}
	}

	_, exists := sessions.Load(sessionId)

	if !exists {
		return &AuthError{ErrorType: AuthErrorCookieNotExists}
	}

	sessions.Delete(sessionId)
	return nil
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
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
