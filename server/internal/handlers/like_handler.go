package handlers

import (
	"net/http"
	"strconv"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ToggleLike(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		postID := c.Param("postId")

		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
			return
		}
		var like models.Like
		result := db.Where("user_id = ? AND post_id = ?", userID, uint(postIDUint)).First(&like)
		if result.Error != nil {
			newLike := models.Like{
				UserID: userID,
				PostID: uint(postIDUint),
			}
			if err := db.Create(&newLike).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like post"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"liked": "true"})
			return
		}
		if err := db.Delete(&like).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"liked": "false"})
	}
}

func GetAllLikes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("postId")
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
			return
		}

		var post models.Post
		if err := db.First(&post, uint(postIDUint)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		var likes []models.Like
		result := db.Where("post_id=?", uint(postIDUint)).Preload("User").Order("created_at desc").
			Find(&likes)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch likes"})
			return
		}

		c.JSON(http.StatusOK, likes)

	}

}
