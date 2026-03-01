package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ComputerSocietyVITC/projects-portal-backend/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

const testSecret = "test-secret"

func makeToken(t *testing.T, userID uuid.UUID, email string, roles []string, exp time.Time) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"roles":   roles,
		"exp":     exp.Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

// newRequest creates a context and returns both the context and response recorder.
func newRequest(e *echo.Echo, method, path string, headers map[string]string) (*echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func setupEcho(handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) (*echo.Echo, echo.HandlerFunc) {
	e := echo.New()
	h := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return e, h
}

func okHandler(c *echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

// ── AuthMiddleware tests ──────────────────────────────────────────────────────

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	e, h := setupEcho(okHandler, middleware.AuthMiddleware())
	c, rec := newRequest(e, http.MethodGet, "/", nil)

	_ = h(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	e, h := setupEcho(okHandler, middleware.AuthMiddleware())
	c, rec := newRequest(e, http.MethodGet, "/", map[string]string{
		"Authorization": "InvalidTokenWithoutBearer",
	})

	_ = h(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	e, h := setupEcho(okHandler, middleware.AuthMiddleware())
	c, rec := newRequest(e, http.MethodGet, "/", map[string]string{
		"Authorization": "Bearer this.is.not.a.valid.token",
	})

	_ = h(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	token := makeToken(t, uuid.New(), "user@example.com", []string{"member"}, time.Now().Add(-time.Hour))

	e, h := setupEcho(okHandler, middleware.AuthMiddleware())
	c, rec := newRequest(e, http.MethodGet, "/", map[string]string{
		"Authorization": "Bearer " + token,
	})

	_ = h(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidToken_SetsContext(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	userID := uuid.New()
	email := "user@example.com"
	roles := []string{"member"}
	token := makeToken(t, userID, email, roles, time.Now().Add(time.Hour))

	var capturedEmail interface{}
	var capturedRoles interface{}

	handler := func(c *echo.Context) error {
		capturedEmail = c.Get("email")
		capturedRoles = c.Get("roles")
		return c.String(http.StatusOK, "ok")
	}

	e, h := setupEcho(handler, middleware.AuthMiddleware())
	c, rec := newRequest(e, http.MethodGet, "/", map[string]string{
		"Authorization": "Bearer " + token,
	})

	_ = h(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if capturedEmail != email {
		t.Errorf("expected email %q, got %v", email, capturedEmail)
	}
	gotRoles, ok := capturedRoles.([]string)
	if !ok || len(gotRoles) == 0 || gotRoles[0] != roles[0] {
		t.Errorf("expected roles %v, got %v", roles, capturedRoles)
	}
}

// ── RequireRole tests ─────────────────────────────────────────────────────────

func TestRequireRole_NoRolesInContext(t *testing.T) {
	e, h := setupEcho(okHandler, middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodGet, "/", nil)
	// do NOT set roles

	_ = h(c)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestRequireRole_InsufficientRole(t *testing.T) {
	e, h := setupEcho(okHandler, middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodGet, "/", nil)
	c.Set("roles", []string{"member"})

	_ = h(c)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestRequireRole_HasRequiredRole(t *testing.T) {
	e, h := setupEcho(okHandler, middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodGet, "/", nil)
	c.Set("roles", []string{"group_head"})

	_ = h(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireRole_MultipleAllowedRoles_UserHasOne(t *testing.T) {
	e, h := setupEcho(okHandler, middleware.RequireRole("group_head", "admin"))
	c, rec := newRequest(e, http.MethodGet, "/", nil)
	c.Set("roles", []string{"admin"})

	_ = h(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireRole_UserHasMultipleRoles_OneMatches(t *testing.T) {
	e, h := setupEcho(okHandler, middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodGet, "/", nil)
	c.Set("roles", []string{"member", "group_head"})

	_ = h(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireRole_EmptyRolesList(t *testing.T) {
	e, h := setupEcho(okHandler, middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodGet, "/", nil)
	c.Set("roles", []string{})

	_ = h(c)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

// ── Auth + RequireRole integration ───────────────────────────────────────────

func TestAuthAndRequireRole_ValidTokenWithRole(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	token := makeToken(t, uuid.New(), "head@example.com", []string{"group_head"}, time.Now().Add(time.Hour))

	e, h := setupEcho(okHandler, middleware.AuthMiddleware(), middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodPost, "/invites", map[string]string{
		"Authorization": "Bearer " + token,
	})

	_ = h(c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAuthAndRequireRole_ValidTokenWithoutRole(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	token := makeToken(t, uuid.New(), "member@example.com", []string{"member"}, time.Now().Add(time.Hour))

	e, h := setupEcho(okHandler, middleware.AuthMiddleware(), middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodPost, "/invites", map[string]string{
		"Authorization": "Bearer " + token,
	})

	_ = h(c)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestAuthAndRequireRole_NoToken(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)

	e, h := setupEcho(okHandler, middleware.AuthMiddleware(), middleware.RequireRole("group_head"))
	c, rec := newRequest(e, http.MethodPost, "/invites", nil)

	_ = h(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}
