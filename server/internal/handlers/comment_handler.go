package handlers

import (
	"net/http"
	"strconv"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
}

func CreateComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		postID := c.Param("postId")
		var input CreateCommentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
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
		newComment := models.Comment{
			UserID:  userID,
			PostID:  uint(postIDUint),
			Content: input.Content,
		}
		if err := db.Create(&newComment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
			return
		}
		if err := db.Preload("User").First(&newComment, newComment.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "comment created but failed to load details"})
			return
		}
		c.JSON(http.StatusCreated, newComment)

	}
}
func GetComments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("postId")
		postIDuint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
			return
		}
		var post models.Post
		if err := db.First(&post, uint(postIDuint)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		var comments []models.Comment
		if err := db.Where("post_id=?", uint(postIDuint)).Preload("User").Order("created_at asc").
			Find(&comments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
			return
		}
		c.JSON(http.StatusOK, comments)
	}
}
func DeleteComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		commentID := c.Param("commentId")

		var comment models.Comment
		if err := db.First(&comment, commentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		var post models.Post
		if err := db.First(&post, comment.PostID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		if comment.UserID != userID && post.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "you can't delete this comment"})
			return
		}

		if err := db.Delete(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
	}
}
