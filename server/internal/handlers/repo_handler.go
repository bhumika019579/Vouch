package handlers

import (
	"net/http"
	"strings"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/bhumika019579/Vouch/server/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserRepos(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		repos, err := services.FetchUserRepos(user.GithubAccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, repos)
	}
}
func LinkRepo(db *gorm.DB, GroqAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var chosenRepo services.GithubRepo
		if err := c.ShouldBindJSON(&chosenRepo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		parts := strings.SplitN(chosenRepo.FullName, "/", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repo full name"})
			return
		}
		owner, repoName := parts[0], parts[1]
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		var existingRepo models.Repo
		err := db.Where("user_id = ? AND repo_url = ?", userID, chosenRepo.HTMLURL).First(&existingRepo).Error
		if err == nil {
			c.JSON(http.StatusOK, existingRepo)
			return
		}
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check for existing repo"})
			return
		}
		token := user.GithubAccessToken
		newRepo, err := services.CreateRepoFromGithub(db, userID, chosenRepo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create repo: " + err.Error()})
			return
		}
		commitCount, err := services.FetchCommitCount(token, owner, repoName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch commit count: " + err.Error()})
			return
		}
		newRepo.CommitCount = commitCount
		languages, err := services.FetchRepoLanguages(token, owner, repoName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch repo languages: " + err.Error()})
			return
		}
		if err := services.TagRepoLanguages(db, newRepo.ID, languages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to tag languages: " + err.Error()})
			return
		}
		readme, err := services.FetchRepoReadme(token, owner, repoName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch readme: " + err.Error()})
			return
		}
		summary, err := services.GenerateRepoSummary(GroqAPIKey, readme)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate summary: " + err.Error()})
			return
		}
		newRepo.Summary = summary
		newRepo.Status = "linked"
		if err := db.Save(newRepo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save repo: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, newRepo)

	}
}
