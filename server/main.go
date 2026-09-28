package main

import (
	"net/http"

	"github.com/bhumika019579/Vouch/server/internal/config"
	"github.com/bhumika019579/Vouch/server/internal/db"
	"github.com/bhumika019579/Vouch/server/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	database := db.Connect(cfg)
	_ = database
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	routes.SetUpRoutes(r, cfg, database)
	r.Run(":" + cfg.Port)
}
