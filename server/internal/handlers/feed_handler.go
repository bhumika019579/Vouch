package handlers

import (
	"net/http"
	"strconv"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetFeed(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil || limit < 1 {
			limit = 20
		}
		offset := (page - 1) * limit
		var posts []models.Post
		if err := db.Preload("User").Preload("Repo").Preload("Repo.RepoTags").Preload("Repo.RepoTags.Tag").
			Order("created_at desc").Limit(limit).Offset(offset).
			Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch feed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"page":  page,
			"limit": limit,
			"posts": posts,
		})

	}
}
