package auth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const sessionName = "gaspoll_session"

// SessionManager keeps Echo handlers decoupled from the underlying store.
type SessionManager struct {
	store *sessions.CookieStore
}

func NewSessionManager(secret string) *SessionManager {
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return &SessionManager{store: store}
}

func (m *SessionManager) SetUser(c echo.Context, userID string) error {
	sess, _ := m.store.Get(c.Request(), sessionName)
	sess.Values["user_id"] = userID
	return sess.Save(c.Request(), c.Response())
}

func (m *SessionManager) Clear(c echo.Context) error {
	sess, _ := m.store.Get(c.Request(), sessionName)
	sess.Options.MaxAge = -1
	return sess.Save(c.Request(), c.Response())
}

func (m *SessionManager) CurrentUserID(c echo.Context) (string, error) {
	sess, _ := m.store.Get(c.Request(), sessionName)
	val, ok := sess.Values["user_id"]
	if !ok {
		return "", errors.New("no user in session")
	}
	userID, ok := val.(string)
	if !ok {
		return "", errors.New("invalid session payload")
	}
	return userID, nil
}

// Middleware ensures handler knows about authenticated user ID if present.
func (m *SessionManager) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := m.CurrentUserID(c)
		if err == nil {
			c.Set("user_id", userID)
		}
		return next(c)
	}
}
