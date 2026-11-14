package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bhumika019579/Vouch/server/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret"

// newGuardedRouter builds a tiny app: the guard plus one dummy route.
// The dummy route returns the user_id the guard stored.
func newGuardedRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware(testSecret))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": c.MustGet("user_id")})
	})
	return r
}

// signToken builds a token from custom claims, signed with the given secret.
func signToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("could not sign token: %v", err)
	}
	return signed
}

// sendPing sends GET /ping with the given Authorization header.
func sendPing(router *gin.Engine, header string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	if header != "" {
		req.Header.Set("Authorization", header)
	}
	router.ServeHTTP(w, req)
	return w
}

func TestAuthMiddlewareRejectsBadHeaders(t *testing.T) {
	router := newGuardedRouter()

	tests := []struct {
		name   string
		header string
	}{
		{"missing header", ""},
		{"no Bearer prefix", "some-token"},
		{"wrong scheme", "Basic abc123"},
		{"Bearer with no token", "Bearer"},
		{"too many parts", "Bearer abc def"},
		{"lowercase bearer", "bearer abc123"},
		{"garbage token", "Bearer not-a-real-token"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := sendPing(router, tc.header)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) {
	router := newGuardedRouter()

	token, err := utils.GenerateJWT(42, testSecret)
	assert.NoError(t, err)

	w := sendPing(router, "Bearer "+token)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"user_id":42}`, w.Body.String())
}

func TestAuthMiddlewareRejectsBadTokens(t *testing.T) {
	router := newGuardedRouter()

	tests := []struct {
		name  string
		token string
	}{
		{"expired token", signToken(t, testSecret, jwt.MapClaims{
			"user_id": 42,
			"exp":     time.Now().Add(-time.Hour).Unix(),
		})},
		{"wrong secret", signToken(t, "other-secret", jwt.MapClaims{
			"user_id": 42,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})},
		{"missing user_id claim", signToken(t, testSecret, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})},
		{"user_id wrong type", signToken(t, testSecret, jwt.MapClaims{
			"user_id": "not-a-number",
			"exp":     time.Now().Add(time.Hour).Unix(),
		})},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := sendPing(router, "Bearer "+tc.token)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}