package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bhumika019579/Vouch/server/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)
var protectedRoutes = []struct {
	name   string
	method string
	path   string
}{
	{"get repos", http.MethodGet, "/api/repos"},
	{"link repo", http.MethodPost, "/api/repos"},
	{"create post", http.MethodPost, "/api/posts"},
	{"delete post", http.MethodDelete, "/api/posts/1"},
	{"create comment", http.MethodPost, "/api/posts/1/comments"},
	{"delete comment", http.MethodDelete, "/api/comments/1"},
	{"toggle like", http.MethodPost, "/api/posts/1/like"},
	{"update profile", http.MethodPatch, "/api/me"},
	{"get me", http.MethodGet, "/api/me"},
}
func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{JWTSecret: "test-secret"}
	return setupRouter(cfg, nil)
}

func TestHealthEndpoint(t *testing.T){
	router:=newTestRouter()
	w:=httptest.NewRecorder()
	req:=httptest.NewRequest(http.MethodGet,"/health",nil)
	router.ServeHTTP(w,req)
	assert.Equal(t,http.StatusOK,w.Code)
	assert.JSONEq(t,`{"status":"ok"}`,w.Body.String())
}
func TestProtectedRoutesWithoutToken(t *testing.T){
	router:=newTestRouter()
	for _,tc:=range protectedRoutes{
		t.Run(tc.name,func(t *testing.T){
			w:=httptest.NewRecorder()
			req:=httptest.NewRequest(tc.method,tc.path,strings.NewReader(`{}`))
			router.ServeHTTP(w,req)
			assert.Equal(t,http.StatusUnauthorized,w.Code)
		})
	}
}
func TestProtectedRoutesWithInvalidToken(t *testing.T) {
	router := newTestRouter()

	for _, tc := range protectedRoutes {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{}`))
			req.Header.Set("Authorization", "Bearer not-a-real-token")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestUnknownRoute(t *testing.T) {
	router := newTestRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}