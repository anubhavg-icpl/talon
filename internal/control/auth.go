package control

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// sessionTTL is the dashboard/CLI session lifetime.
const sessionTTL = 30 * 24 * time.Hour

// sessionCookieName rides both the browser (Next proxy) and the WS upgrade.
const sessionCookieName = "talon_session"

// authStore is the persistence Auth needs — *DB implements it, tests fake it.
type authStore interface {
	seedUser(ctx context.Context, username, passwordHash string) error
	passwordHashFor(ctx context.Context, username string) (userID, hash string, err error)
	createSession(ctx context.Context, userID, token string, ttl time.Duration) error
	sessionUser(ctx context.Context, token string) string
	deleteSession(ctx context.Context, token string)
}

// Auth gates the API behind username/password sessions backed by Postgres.
// A nil *Auth on the Server means auth is disabled (dev escape hatch or no
// database) — every request passes through, matching the old behavior.
// cache, when non-nil, memoizes token→user lookups (one PG query per request
// otherwise) for sessionCacheTTL; logout evicts.
type Auth struct {
	db    authStore
	cache Cache
}

// NewAuth seeds the initial admin user (idempotent) and returns the Auth.
func NewAuth(ctx context.Context, db authStore, adminUser, adminPassword string) (*Auth, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if err := db.seedUser(ctx, adminUser, string(hash)); err != nil {
		return nil, err
	}
	log.Printf("control: auth enabled — admin user %q ready", adminUser)
	return &Auth{db: db}, nil
}

// SetCache enables Redis-backed session memoization.
func (a *Auth) SetCache(c Cache) { a.cache = c }

// newToken returns 256 bits of randomness, hex-encoded.
func newToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// usernameForRequest resolves the session cookie, bearer token, or access_token
// query (used by the Arsenal Shell iframe when the dashboard and core sit on
// different ports and the HttpOnly cookie is origin-scoped).
func (a *Auth) usernameForRequest(r *http.Request) string {
	token := ""
	if c, err := r.Cookie(sessionCookieName); err == nil {
		token = c.Value
	}
	if token == "" {
		if h := r.Header.Get("Authorization"); len(h) > 7 && h[:7] == "Bearer " {
			token = h[7:]
		}
	}
	if token == "" {
		// Prefer access_token; token is a common alternate used by local tooling.
		if q := r.URL.Query().Get("access_token"); q != "" {
			token = q
		} else if q := r.URL.Query().Get("token"); q != "" {
			token = q
		}
	}
	if token == "" {
		return ""
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if a.cache != nil {
		if username, ok := a.cache.Get(ctx, cacheKeySessionPrefix+token); ok {
			return username
		}
	}
	username := a.db.sessionUser(ctx, token)
	if username != "" && a.cache != nil {
		a.cache.Set(ctx, cacheKeySessionPrefix+token, username, sessionCacheTTL)
	}
	return username
}

// loginRequest is the POST /auth/login body.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// handleLogin is POST /auth/login — verifies credentials, creates a session,
// sets the cookie, and returns the token (for CLI clients).
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if s.auth == nil {
		writeError(w, http.StatusServiceUnavailable, "auth is disabled on this core")
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	userID, hash, err := s.auth.db.passwordHashFor(ctx, req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := newToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "token generation failed")
		return
	}
	if err := s.auth.db.createSession(ctx, userID, token, sessionTTL); err != nil {
		writeError(w, http.StatusInternalServerError, "session creation failed")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(sessionTTL),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, map[string]string{"token": token, "username": req.Username})
}

// handleLogout is POST /auth/logout — drops the session and clears the cookie.
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if s.auth != nil {
		if c, err := r.Cookie(sessionCookieName); err == nil {
			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()
			s.auth.db.deleteSession(ctx, c.Value)
			if s.auth.cache != nil {
				s.auth.cache.Del(ctx, cacheKeySessionPrefix+c.Value)
			}
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

// handleMe is GET /auth/me — the current session's user.
func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	if s.auth == nil {
		writeJSON(w, http.StatusOK, map[string]string{"username": "operator", "auth": "disabled"})
		return
	}
	username := s.auth.usernameForRequest(r)
	if username == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"username": username})
}

// openPath reports whether a route skips the auth middleware.
func openPath(path string) bool {
	switch path {
	case "/health", "/health/services", "/auth/login":
		return true
	}
	return false
}
