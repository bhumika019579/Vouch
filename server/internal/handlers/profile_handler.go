package handlers

import (
	"net/http"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateProfileInput struct {
	Bio    string `json:"bio"`
	Skills string `json:"skills"`
}

func GetUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var user models.User
		if err := db.Where("github_username=?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		var posts []models.Post
		if err := db.Where("user_id=?", user.ID).
			Preload("Repo").
			Preload("Repo.RepoTags").
			Preload("Repo.RepoTags.Tag").
			Preload("User").
			Order("created_at desc").
			Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"posts": posts,
		})
	}
}
func UpdateUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var input UpdateProfileInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		user.Bio = input.Bio
		user.Skills = input.Skills
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
func GetMe(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
