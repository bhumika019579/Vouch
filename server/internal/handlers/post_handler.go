package handlers

import (
	"net/http"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatePostInput struct {
	RepoID  uint   `json:"repo_id" binding:"required"`
	Caption string `json:"caption"`
}

func CreatePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var input CreatePostInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invaild request body"})
			return
		}
		var repo models.Repo
		if err := db.First(&repo, input.RepoID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "repo not found"})
			return
		}
		if repo.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "you dont own this repo"})
			return
		}
		newPost := models.Post{
			UserID:  userID,
			RepoID:  input.RepoID,
			Caption: input.Caption,
		}
		if err := db.Create(&newPost).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
			return
		}
		if err := db.Preload("User").Preload("Repo").Preload("Repo.RepoTags").
			Preload("Repo.RepoTags.Tag").First(&newPost, newPost.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "post created but failed to fetch its details"})
			return
		}
		c.JSON(http.StatusCreated, newPost)
	}
}
func DeletePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		postID := c.Param("postId")
		var post models.Post
		if err := db.First(&post, postID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		if post.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't own this post"})
			return
		}
		if err := db.Where("post_id = ?", post.ID).Delete(&models.Comment{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post's comments"})
			return
		}
		if err := db.Where("post_id = ?", post.ID).Delete(&models.Like{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post's likes"})
			return
		}
		if err := db.Delete(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "post deleted"})

	}
}
