package services

import (
	"github.com/bhumika019579/Vouch/server/internal/models"
	"gorm.io/gorm"
)

func CreateRepoFromGithub(db *gorm.DB, userID uint, repo GithubRepo) (*models.Repo, error) {
	newRepo := models.Repo{
		UserID:      userID,
		RepoURL:     repo.HTMLURL,
		RepoName:    repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		Status:      "pending",
	}

	if err := db.Create(&newRepo).Error; err != nil {
		return nil, err
	}

	return &newRepo, nil
}
