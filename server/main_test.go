package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bhumika019579/Vouch/server/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T){
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{JWTSecret: "test-secret"}
	router:=setupRouter(cfg,nil)
	w:=httptest.NewRecorder()
	req:=httptest.NewRequest(http.MethodGet,"/health",nil)
	router.ServeHTTP(w,req)
	assert.Equal(t,http.StatusOK,w.Code)
	assert.JSONEq(t,`{"status":"ok"}`,w.Body.String())
}