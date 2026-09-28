package services

import (
	"log"

	"github.com/bhumika019579/Vouch/server/internal/models"
	"gorm.io/gorm"
)

func TagRepoLanguages(db *gorm.DB, repoID uint, languages []LanguageBreakdown) error {
	for _, lang := range languages {
		var tag models.Tag
		result := db.Where("name=?", lang.Name).First(&tag)
		if result.Error != nil {
			tag = models.Tag{Name: lang.Name}
			if err := db.Create(&tag).Error; err != nil {
				log.Printf("failed to create tag %s: %v", lang.Name, err)
				continue
			}
		}
		repoTag := models.RepoTag{
			RepoID:     repoID,
			TagID:      tag.ID,
			Percentage: lang.Percentage,
		}
		if err := db.Create(&repoTag).Error; err != nil {
			log.Printf("failed to tag repo %d with language %s: %v", repoID, lang.Name, err)
			continue
		}
	}
	return nil

}
