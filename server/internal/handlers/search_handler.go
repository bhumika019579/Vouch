package handlers

import (
	"net/http"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SearchRepos(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
			return
		}
		searchTeam := "%" + query + "%"
		var posts []models.Post
		if err := db.Preload("User").Preload("Repo").Preload("Repo.RepoTags").Preload("Repo.RepoTags.Tag").
			Joins("JOIN repos ON repos.id=posts.repo_id").
			Joins("JOIN users ON users.id=posts.user_id").
			Joins("JOIN repo_tags ON repo_tags.repo_id = repos.id").
			Joins("JOIN tags ON tags.id = repo_tags.tag_id").
			Where("repos.repo_name ILIKE ? OR users.github_username ILIKE ? OR tags.name ILIKE ?", searchTeam, searchTeam, searchTeam).
			Distinct("posts.*").
			Order("posts.created_at DESC").
			Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": posts})

	}
}
