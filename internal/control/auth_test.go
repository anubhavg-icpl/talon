package control

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeAuthStore is an in-memory authStore for tests.
type fakeAuthStore struct {
	mu       sync.Mutex
	users    map[string]string // username -> bcrypt hash
	sessions map[string]string // token -> username
}

func newFakeAuthStore() *fakeAuthStore {
	return &fakeAuthStore{users: map[string]string{}, sessions: map[string]string{}}
}

func (f *fakeAuthStore) seedUser(_ context.Context, username, hash string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.users[username]; !ok {
		f.users[username] = hash
	}
	return nil
}

func (f *fakeAuthStore) passwordHashFor(_ context.Context, username string) (string, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	hash, ok := f.users[username]
	if !ok {
		return "", "", errors.New("no such user")
	}
	return "uid-" + username, hash, nil
}

func (f *fakeAuthStore) createSession(_ context.Context, userID, token string, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sessions[token] = strings.TrimPrefix(userID, "uid-")
	return nil
}

func (f *fakeAuthStore) sessionUser(_ context.Context, token string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.sessions[token]
}

func (f *fakeAuthStore) deleteSession(_ context.Context, token string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.sessions, token)
}

func newAuthedServer(t *testing.T) *Server {
	t.Helper()
	auth, err := NewAuth(context.Background(), newFakeAuthStore(), "admin", "s3cret")
	if err != nil {
		t.Fatalf("NewAuth: %v", err)
	}
	return NewServer(nil, NewStore(), WithAuth(auth))
}

func loginBody(user, pass string) *strings.Reader {
	return strings.NewReader(`{"username":"` + user + `","password":"` + pass + `"}`)
}

func TestLoginFlowAndMiddleware(t *testing.T) {
	t.Parallel()
	srv := newAuthedServer(t)
	h := srv.Handler()

	// Protected route without credentials → 401.
	req := httptest.NewRequest(http.MethodGet, "/runs", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("no-auth /runs status=%d want 401", rec.Code)
	}

	// Open routes stay open.
	for _, path := range []string{"/health", "/health/services"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusUnauthorized {
			t.Fatalf("open route %s returned 401", path)
		}
	}

	// Bad password → 401.
	req = httptest.NewRequest(http.MethodPost, "/auth/login", loginBody("admin", "wrong"))
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("bad-login status=%d want 401", rec.Code)
	}

	// Good login → 200 + token + cookie.
	req = httptest.NewRequest(http.MethodPost, "/auth/login", loginBody("admin", "s3cret"))
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login status=%d want 200 (%s)", rec.Code, rec.Body.String())
	}
	var cookie *http.Cookie
	for _, c := range rec.Result().Cookies() {
		if c.Name == sessionCookieName {
			cookie = c
		}
	}
	if cookie == nil || cookie.Value == "" {
		t.Fatal("login did not set session cookie")
	}

	// Cookie auth works on a protected route.
	req = httptest.NewRequest(http.MethodGet, "/runs", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("cookie /runs status=%d want 200", rec.Code)
	}

	// Bearer auth works too (CLI path).
	req = httptest.NewRequest(http.MethodGet, "/runs", nil)
	req.Header.Set("Authorization", "Bearer "+cookie.Value)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("bearer /runs status=%d want 200", rec.Code)
	}

	// /auth/me resolves the user.
	req = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "admin") {
		t.Fatalf("me status=%d body=%s", rec.Code, rec.Body.String())
	}

	// Logout kills the session.
	req = httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("logout status=%d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/runs", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("post-logout /runs status=%d want 401", rec.Code)
	}
}

func TestAuthDisabledPassesThrough(t *testing.T) {
	t.Parallel()
	srv := NewServer(nil, NewStore()) // no WithAuth
	req := httptest.NewRequest(http.MethodGet, "/runs", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("auth-disabled /runs status=%d want 200", rec.Code)
	}
}

func TestSeedUserIdempotent(t *testing.T) {
	t.Parallel()
	store := newFakeAuthStore()
	if _, err := NewAuth(context.Background(), store, "admin", "one"); err != nil {
		t.Fatal(err)
	}
	// Second seed with a different password must not overwrite the first.
	if _, err := NewAuth(context.Background(), store, "admin", "two"); err != nil {
		t.Fatal(err)
	}
	if len(store.users) != 1 {
		t.Fatalf("users=%d want 1", len(store.users))
	}
	srv := NewServer(nil, NewStore(), WithAuth(&Auth{db: store}))
	req := httptest.NewRequest(http.MethodPost, "/auth/login", loginBody("admin", "one"))
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login with original password failed: %d", rec.Code)
	}
}
